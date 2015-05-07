package kmgCrypto

import (
	"bytes"
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
)

func TestKmgCipherV2(ot *testing.T) {
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
		ob := EncryptV2([]byte("1"), origin)
		ret, err := DecryptV2([]byte("1"), ob)
		kmgTest.Equal(err, nil)
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
		_, err := Decrypt([]byte("1"), origin)
		kmgTest.Ok(err != nil)
	}
	origin := []byte("1234567890123456712345678901234567")
	//多次加密得到不同的结果 (随机化iv)
	kmgTest.Ok(!bytes.Equal(EncryptV2([]byte("1"), origin), EncryptV2([]byte("1"), origin)))
	//修改任何一个字节都会报错 (hash)

	ob := EncryptV2([]byte("1"), origin)
	for i := 0; i < len(ob); i++ {
		newOb := make([]byte, len(ob))
		newOb[i] = -newOb[i]
		_, err := DecryptV2([]byte("1"), newOb)
		kmgTest.Ok(err != nil)
	}
}
