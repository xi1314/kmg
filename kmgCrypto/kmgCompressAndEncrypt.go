package kmgCrypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"github.com/bronze1man/kmg/kmgCompress"
	"github.com/bronze1man/kmg/kmgRand"
)

/*
先压缩,后加密,
对称加密,正确包含psk的所有功能,
不管是否可以压缩.保证最多会比明文数据增加58个字节.
*/
func CompressAndEncryptBytesEncode(key *[32]byte, data []byte) (output []byte) {
	//先压缩 压缩不了只会多16个字节而已,不需要进行优化
	data = kmgCompress.FlateMustCompress(data)
	//后加密
	cbcIv := kmgRand.MustCryptoRandBytes(16) //此处只会报操作系统不支持的错误.
	block, err := aes.NewCipher((*key)[:])
	if err != nil {
		panic(err) //此处只会由于key长度错误而报错,而此处key长度不会错误.
	}
	blockSize := block.BlockSize()
	paddingSize := blockSize - len(data)%blockSize
	afterCbcSize := paddingSize + len(data)
	output = make([]byte, afterCbcSize+16+16)
	copy(output[:16], cbcIv)
	copy(output[16:], data)
	hashData := sha512.Sum512(data)
	copy(output[len(data)+16:len(data)+16+16], hashData[:16]) //只有前16位,其他的废弃
	copy(output[len(data)+16+16:], bytes.Repeat([]byte{byte(paddingSize)}, paddingSize))
	blockmode := cipher.NewCBCEncrypter(block, cbcIv)
	blockmode.CryptBlocks(output[16:], output[16:])
	return output
}

/*
对称解密,
	不会修改输入的数据
*/
func CompressAndEncryptBytesDecode(key *[32]byte, data []byte) (output []byte, err error) {
	//先解密
	if len(data) < 16+16 {
		return nil, errors.New("[kmgCipher.CompressAndEncryptBytesDecode] input data too small")
	}
	aseKey := key[:]
	cbcIv := data[:16]
	block, err := aes.NewCipher(aseKey)
	if err != nil {
		return nil, err
	}
	if len(data)%block.BlockSize() != 0 {
		return nil, errors.New("[kmgCipher.CompressAndEncryptBytesDecode] input not full blocks")
	}
	output = make([]byte, len(data)-16)
	blockmode := cipher.NewCBCDecrypter(block, cbcIv)
	blockmode.CryptBlocks(output, data[16:])
	paddingSize := int(output[len(output)-1])
	if paddingSize > block.BlockSize() {
		return nil, errors.New("[kmgCipher.CompressAndEncryptBytesDecode] paddingSize out of range")
	}
	beforeCbcSize := len(data) - paddingSize - 16 - 16
	hashData := sha512.Sum512(output[:beforeCbcSize])
	if !bytes.Equal(hashData[:16], output[beforeCbcSize:beforeCbcSize+16]) {
		return nil, errors.New("[kmgCipher.CompressAndEncryptBytesDecode] hash not match")
	}
	output = output[:beforeCbcSize]
	//后解压缩
	output, err = kmgCompress.FlateUnCompress(output)
	return output, err
}

func CompressAndEncryptBase64Encode(key *[32]byte, data []byte) (output string) {
	outB := CompressAndEncryptBytesEncode(key, data)
	return base64.URLEncoding.EncodeToString(outB)
}

func CompressAndEncryptBase64Decode(key *[32]byte, data string) (output []byte, err error) {
	dataB, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return CompressAndEncryptBytesDecode(key, dataB)
}
