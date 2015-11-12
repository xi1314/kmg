package kmgCrypto

import (
	"testing"

	"github.com/bronze1man/kmg/kmgTest"
)

func Test1(ot *testing.T) {
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
	} {
		ob, err := Encrypt([]byte("1"), origin)
		kmgTest.Equal(err, nil)
		ret, err := Decrypt([]byte("1"), ob)
		kmgTest.Equal(err, nil)
		kmgTest.Equal(ret, origin)

		sob, err := EncryptString("1", origin)
		kmgTest.Equal(err, nil)
		sret, err := DecryptString("1", sob)
		kmgTest.Equal(err, nil)
		kmgTest.Equal(sret, origin)
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
}
