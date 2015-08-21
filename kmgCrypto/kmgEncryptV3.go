package kmgCrypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"github.com/bronze1man/kmg/kmgRand"
)

/*
只加密,不压缩
对称加密,正确包含psk的所有功能,
让Aes加密和magicCode进行传输错误验证.
AES-CTR magicCode完整性验证 最坏情况多20个字节
*/
func EncryptV3(key *[32]byte, data []byte) (output []byte) {
	Iv := kmgRand.MustCryptoRandBytes(16) //此处只会报操作系统不支持的错误.
	block, err := aes.NewCipher((*key)[:])
	if err != nil {
		panic(err) //此处只会由于key长度错误而报错,而此处key长度不会错误.
	}
	afterCbcSize := len(data)
	output = make([]byte, afterCbcSize+16+4)
	copy(output[:16], Iv)
	copy(output[16:len(data)+16], data)
	copy(output[len(data)+16:len(data)+16+4], magicCode4)
	ctr := cipher.NewCTR(block, Iv)
	ctr.XORKeyStream(output[16:], output[16:])
	return output
}

/*
对称解密,
	不会修改输入的数据
*/
func DecryptV3(key *[32]byte, data []byte) (output []byte, err error) {
	//先解密
	if len(data) < 20 {
		return nil, errors.New("[kmgCipher.CompressAndEncryptBytesDecode] input data too small")
	}
	aseKey := key[:]
	Iv := data[:16]
	block, err := aes.NewCipher(aseKey)
	if err != nil {
		return nil, err
	}
	output = make([]byte, len(data)-16)
	ctr := cipher.NewCTR(block, Iv)
	ctr.XORKeyStream(output, data[16:])
	beforeCbcSize := len(data) - 16 - 4
	if !bytes.Equal(magicCode4, output[beforeCbcSize:beforeCbcSize+4]) {
		return nil, errors.New("[kmgCipher.CompressAndEncryptBytesDecode] magicCode not match")
	}
	output = output[:beforeCbcSize]
	//后解压缩
	return output, nil
}
