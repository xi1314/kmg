package kmgCrypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"github.com/bronze1man/kmg/kmgCompress"
	"github.com/bronze1man/kmg/kmgRand"
)

var magicCode4 = []byte{0xa7, 0x97, 0x6d, 0x15}

/*
先压缩,后加密,
对称加密,正确包含psk的所有功能,
让Aes加密和magicCode进行传输错误验证.
压缩网络包时,效果不明显,劣化也不明显,大部分包都无法压缩,一般仅有3%左右的包可以压缩.
AES-CTR magicCode完整性验证 zlib压缩(1字节最坏情况控制) 最坏情况多21个字节
*/
func CompressAndEncryptBytesEncodeV2(key *[32]byte, data []byte) (output []byte) {
	//先压缩
	data = compressV2(data)
	//后加密
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
func CompressAndEncryptBytesDecodeV2(key *[32]byte, data []byte) (output []byte, err error) {
	//先解密
	if len(data) < 16+4+1 {
		return nil, errors.New("[kmgCipher.CompressAndEncryptBytesDecode] input data too small")
	}
	aseKey := key[:]
	Iv := data[:16] //此处把cap改低一点,使后面的append不会覆盖已有数据
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
	output, err = uncompressV2(output)
	return output, err
}

func compressV2(inData []byte) (outData []byte) {
	outData = kmgCompress.ZlibMustCompress(inData)
	if len(outData) >= len(inData) {
		return append([]byte{0}, inData...)
	}
	return append([]byte{1}, outData...)
}
func uncompressV2(inData []byte) (outData []byte, err error) {
	if len(inData) == 0 {
		return nil, errors.New("[uncopressV2] len(inData)==0")
	}
	if inData[0] == 0 {
		return inData[1:], nil
	}
	return kmgCompress.ZlibUnCompress(inData[1:])
}
