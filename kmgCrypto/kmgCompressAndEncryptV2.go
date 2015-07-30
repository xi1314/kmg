package kmgCrypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"errors"
	"github.com/bronze1man/kmg/kmgCompress"
	"github.com/bronze1man/kmg/kmgRand"
)

/*
先压缩,后加密,
对称加密,正确包含psk的所有功能,
特别优化体积
实际iv长度只有4字节,实际hash长度只有4字节
不管是否可以压缩.保证最多会比明文数据增加9个字节.
压缩网络包时,效果不明显,劣化也不明显,大部分包都无法压缩,一般仅有3%左右的包可以压缩.
*/
func CompressAndEncryptBytesEncodeV2(key *[32]byte, data []byte) (output []byte) {
	//先压缩
	data = compressV2(data)
	//后加密
	Iv := kmgRand.MustCryptoRandBytes(4) //此处只会报操作系统不支持的错误.
	block, err := aes.NewCipher((*key)[:])
	if err != nil {
		panic(err) //此处只会由于key长度错误而报错,而此处key长度不会错误.
	}
	afterCbcSize := len(data)
	output = make([]byte, afterCbcSize+4+4)
	copy(output[:4], Iv)
	copy(output[4:], data)
	hashData := sha512.Sum512(data)
	copy(output[len(data)+4:len(data)+4+4], hashData[:4]) //只有前4位,其他的废弃
	ctrIv := append(Iv, (*key)[:12]...)
	ctr := cipher.NewCTR(block, ctrIv)
	ctr.XORKeyStream(output[4:], output[4:])
	return output
}

/*
对称解密,
	不会修改输入的数据
*/
func CompressAndEncryptBytesDecodeV2(key *[32]byte, data []byte) (output []byte, err error) {
	//先解密
	if len(data) < 4+4+1 {
		return nil, errors.New("[kmgCipher.CompressAndEncryptBytesDecode] input data too small")
	}
	aseKey := key[:]
	Iv := data[:4:4] //此处把cap改低一点,使后面的append不会覆盖已有数据
	block, err := aes.NewCipher(aseKey)
	if err != nil {
		return nil, err
	}
	output = make([]byte, len(data)-4)
	ctrIv := append(Iv, (*key)[:12]...)
	ctr := cipher.NewCTR(block, ctrIv)
	ctr.XORKeyStream(output, data[4:])
	beforeCbcSize := len(data) - 4 - 4
	hashData := sha512.Sum512(output[:beforeCbcSize])
	if !bytes.Equal(hashData[:4], output[beforeCbcSize:beforeCbcSize+4]) {
		return nil, errors.New("[kmgCipher.CompressAndEncryptBytesDecode] hash not match")
	}
	output = output[:beforeCbcSize]
	//后解压缩
	output, err = uncompressV2(output)
	return output, err
}

func compressV2(inData []byte) (outData []byte) {
	outData = kmgCompress.FlateMustCompress(inData)
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
	return kmgCompress.FlateUnCompress(inData[1:])
}

/*
func compressV3(inData []byte) (outData []byte) {
	outData,err := zappy.Encode(nil,inData)
	if err!=nil{
		panic(err)
	}
	if len(outData) >= len(inData) {
		return append([]byte{0}, inData...)
	}
	return append([]byte{1}, outData...)
}
func uncompressV3(inData []byte) (outData []byte, err error) {
	if len(inData) == 0 {
		return nil, errors.New("[uncopressV2] len(inData)==0")
	}
	if inData[0] == 0 {
		return inData[1:], nil
	}
	return zappy.Decode(nil,inData[1:])
}

func compressV4(inData []byte) (outData []byte) {
	return inData
}
func uncompressV4(inData []byte) (outData []byte, err error) {
	return outData,nil
}
*/
