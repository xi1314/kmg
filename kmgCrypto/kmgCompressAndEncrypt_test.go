package kmgCrypto

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/bronze1man/kmg/kmgRand"
	"github.com/bronze1man/kmg/kmgTest"
)

func TestCompressAndEncryptBytes(ot *testing.T) {
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
		ob := CompressAndEncryptBytesEncode(key, origin)
		ret, err := CompressAndEncryptBytesDecode(key, ob)
		kmgTest.EqualMsg(err, nil, origin)
		kmgTest.EqualMsg(ret, origin)
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
		_, err := Decrypt([]byte("1"), origin)
		kmgTest.Ok(err != nil)
	}
	origin := []byte("1234567890123456712345678901234567")
	//多次加密得到不同的结果 (随机化iv)
	kmgTest.Ok(!bytes.Equal(CompressAndEncryptBytesEncode(key, origin), CompressAndEncryptBytesEncode(key, origin)))
	//修改任何一个字节都会报错 (hash)

	ob := CompressAndEncryptBytesEncode(key, origin)
	for i := 0; i < len(ob); i++ {
		newOb := make([]byte, len(ob))
		newOb[i] = -newOb[i]
		_, err := CompressAndEncryptBytesDecode(key, newOb)
		kmgTest.Ok(err != nil)
	}

	fmt.Printf("%#v", kmgRand.MustCryptoRandBytes(32))
}
