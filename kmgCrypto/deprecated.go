package kmgCrypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
)

// @deprecated
var GenUUIDErrors = errors.New("gen uuid fail")

// @deprecated
func GenUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := rand.Read(uuid)
	if n != len(uuid) || err != nil {
		return "", GenUUIDErrors
	}
	return hex.EncodeToString(uuid), nil
}

// @deprecated
func MustGenUUID() string {
	val, err := GenUUID()
	if err != nil {
		panic(err)
	}
	return val
}

/*
对称加密,  安全性没有这么简单
 key为任意长度,使用简单
 data为任意长度,使用简单
 使用aes,cbc,32位密码
 输入密码hash使用sha384
 数据padding使用PKCS5Padding
 不会修改输入的数据
 @deprecated
*/
func Encrypt(key []byte, data []byte) (output []byte, err error) {
	keyHash := sha512.Sum384(key)
	aseKey := keyHash[:32]
	cbcIv := keyHash[32:]
	block, err := aes.NewCipher(aseKey)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	paddingSize := blockSize - len(data)%blockSize
	afterCbcSize := paddingSize + len(data)
	output = make([]byte, afterCbcSize)
	copy(output, data)
	copy(output[len(data):], bytes.Repeat([]byte{byte(paddingSize)}, paddingSize))
	blockmode := cipher.NewCBCEncrypter(block, cbcIv)
	blockmode.CryptBlocks(output, output)
	return output, nil
}

/*
对称解密,
	不会修改输入的数据
	@deprecated
*/
func Decrypt(key []byte, data []byte) (output []byte, err error) {
	if len(data) == 0 {
		return nil, errors.New("[kmgCipher.Decrypt] input data zero length")
	}
	keyHash := sha512.Sum384(key)
	aseKey := keyHash[:32]
	cbcIv := keyHash[32:]
	block, err := aes.NewCipher(aseKey)
	if err != nil {
		return nil, err
	}
	if len(data)%block.BlockSize() != 0 {
		return nil, errors.New("[kmgCipher.Decrypt] input not full blocks")
	}
	output = make([]byte, len(data))
	blockmode := cipher.NewCBCDecrypter(block, cbcIv)
	blockmode.CryptBlocks(output, data)
	paddingSize := int(output[len(output)-1])
	if paddingSize > block.BlockSize() {
		return nil, errors.New("[kmgCipher.Decrypt] paddingSize out of range")
	}
	beforeCbcSize := len(data) - paddingSize
	return output[:beforeCbcSize], nil
}

// @deprecated
func EncryptString(key string, data []byte) (output string, err error) {
	outputByte, err := Encrypt([]byte(key), []byte(data))
	if err != nil {
		return
	}
	return base64.URLEncoding.EncodeToString(outputByte), nil
}

// @deprecated
func DecryptString(key string, data string) (output []byte, err error) {
	dataByte, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		return
	}
	outputByte, err := Decrypt([]byte(key), dataByte)
	if err != nil {
		return
	}
	return outputByte, nil
}
