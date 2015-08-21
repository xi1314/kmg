package kmgCrypto

import (
	"bytes"
	"testing"

	"github.com/bronze1man/kmg/kmgRand"
	"github.com/bronze1man/kmg/kmgTest"
)

func TestCompressAndEncryptBytes(ot *testing.T) {
	EncryptTester(CompressAndEncryptBytesEncode, CompressAndEncryptBytesDecode, 58)
}

func EncryptTester(encrypt func(key *[32]byte, data []byte) (output []byte),
	decrypt func(key *[32]byte, data []byte) (output []byte, err error),
	maxOverhead int) {
	key := &[32]byte{0xd8, 0x51, 0xea, 0x81, 0xb9, 0xe, 0xf, 0x2f, 0x8c, 0x85, 0x5f, 0xb6, 0x14, 0xb2}
	//加密数据,可以正确解密测试
	for _, origin := range [][]byte{
		[]byte(""),
		[]byte("1"),
		[]byte("12"),
		[]byte("123"),
		[]byte("1234"),
		[]byte("12345"),
		[]byte("123456"),
		[]byte("1234567"),
		[]byte("12345678"),
		[]byte("123456789"),
		[]byte("1234567890"),
		[]byte("123456789012345"),
		[]byte("1234567890123456"),
		[]byte("12345678901234567"),
		bytes.Repeat([]byte("1234567890"), 100),
	} {
		ob := encrypt(key, origin)
		ret, err := decrypt(key, ob)
		kmgTest.Equal(err, nil, origin)
		kmgTest.Equal(ret, origin)
	}
	//任意数据传入解密不会挂掉,并且会报错
	for _, origin := range [][]byte{
		[]byte(""),
		[]byte("1"),
		[]byte("12"),
		[]byte("123"),
		[]byte("1234"),
		[]byte("12345"),
		[]byte("123456"),
		[]byte("1234567"),
		[]byte("12345678"),
		[]byte("123456789"),
		[]byte("1234567890"),
		[]byte("123456789012345"),
		[]byte("1234567890123456"),
		[]byte("12345678901234567"),
	} {
		_, err := decrypt(key, origin)
		kmgTest.Ok(err != nil)
	}
	origin := []byte("1234567890123456712345678901234567")
	//多次加密得到不同的结果 (随机化iv)
	kmgTest.Ok(!bytes.Equal(encrypt(key, origin), encrypt(key, origin)))
	//修改任何一个字节都会报错 (hash)

	ob := encrypt(key, origin)
	for i := 0; i < len(ob); i++ {
		newOb := make([]byte, len(ob))
		newOb[i] = -newOb[i]
		_, err := decrypt(key, newOb)
		kmgTest.Ok(err != nil)
	}

	for _, i := range []int{1, 10, 100, 1000, 10000} {
		ob := encrypt(key, kmgRand.MustCryptoRandBytes(i))
		kmgTest.Ok(len(ob)-i <= maxOverhead, i)
	}
}
