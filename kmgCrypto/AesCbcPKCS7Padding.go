package kmgCrypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

var ErrDecryptedDataTooSmall = fmt.Errorf("encrypted data too small")

//key must be 32 bytes, if it is not 32 byte,it will panic
// 警告: 这个加密没有hash认证
// 0-16: iv
// 16-32: data
func AesCbcPKCS7PaddingEncrypt(in, key []byte) (out []byte) {
	out = make([]byte, len(in)+aes.BlockSize*2)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	in = PKCS7Padding(in)
	iv := make([]byte, aes.BlockSize)
	_, err = rand.Read(iv)
	if err != nil {
		panic(err)
	}
	copy(out[:aes.BlockSize], iv)
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(out[aes.BlockSize:], in)

	return out[:len(in)+aes.BlockSize]
}

//key must be 32 bytes, if it is not 32 byte,it will panic.
// only repeat error because of input data
func AesCbcPKCS7PaddingDecrypt(in, key []byte) (out []byte, err error) {
	if len(in) < aes.BlockSize*2 {
		return nil, ErrDecryptedDataTooSmall
	}
	out = make([]byte, len(in))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	iv := in[:aes.BlockSize]
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(out, in[aes.BlockSize:])
	return UnPKCS7Padding(out[:len(in)-aes.BlockSize]), nil
}

func PKCS7Padding(data []byte) []byte {
	blockSize := 16
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)

}

func UnPKCS7Padding(data []byte) []byte {
	length := len(data)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}
