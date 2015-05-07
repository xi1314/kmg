package kmgCrypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"errors"
	"github.com/bronze1man/kmg/kmgRand"
)

/*
对称加密,正确包含psk的所有功能,不考虑性能
 key为任意长度
 data为任意长度
 会在原文的最后面加上原文的sha512的内容,并且加密
 使用随机化iv
 使用aes,cbc,32位密码
 输入密码hash使用sha384
 数据padding使用PKCS5Padding
 不会修改输入的数据
*/
func EncryptV2(key []byte, data []byte) (output []byte) {
	keyHash := sha512.Sum512(key)
	aseKey := keyHash[:32]
	cbcIv := kmgRand.MustCryptoRandBytes(16) //除了操作系统不支持,此处不会报错.
	block, err := aes.NewCipher(aseKey)
	if err != nil {
		panic(err) //此处只会由于key长度错误而报错,此处key长度不会错误
	}
	blockSize := block.BlockSize()
	paddingSize := blockSize - len(data)%blockSize
	afterCbcSize := paddingSize + len(data)
	output = make([]byte, afterCbcSize+64+16)
	copy(output[:16], cbcIv)
	copy(output[16:], data)
	hashData := sha512.Sum512(data)
	copy(output[len(data)+16:len(data)+64+16], hashData[:])
	copy(output[len(data)+64+16:], bytes.Repeat([]byte{byte(paddingSize)}, paddingSize))
	blockmode := cipher.NewCBCEncrypter(block, cbcIv)
	blockmode.CryptBlocks(output[16:], output[16:])
	return output
}

/*
对称解密,
	不会修改输入的数据
*/
func DecryptV2(key []byte, data []byte) (output []byte, err error) {
	if len(data) < 16+64 {
		return nil, errors.New("[kmgCipher.Decrypt] input data too small")
	}
	keyHash := sha512.Sum512(key)
	aseKey := keyHash[:32]
	cbcIv := data[:16]
	block, err := aes.NewCipher(aseKey)
	if err != nil {
		return nil, err
	}
	if len(data)%block.BlockSize() != 0 {
		return nil, errors.New("[kmgCipher.Decrypt] input not full blocks")
	}
	output = make([]byte, len(data)-16)
	blockmode := cipher.NewCBCDecrypter(block, cbcIv)
	blockmode.CryptBlocks(output, data[16:])
	paddingSize := int(output[len(output)-1])
	if paddingSize > block.BlockSize() {
		return nil, errors.New("[kmgCipher.Decrypt] paddingSize out of range")
	}
	beforeCbcSize := len(data) - paddingSize - 64 - 16
	hashData := sha512.Sum512(output[:beforeCbcSize])
	if !bytes.Equal(hashData[:], output[beforeCbcSize:beforeCbcSize+64]) {
		return nil, errors.New("[kmgCipher.Decrypt] hash not match")
	}
	return output[:beforeCbcSize], nil
}
