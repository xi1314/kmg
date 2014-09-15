package kmgRand

import (
	"crypto/rand"
	"encoding/hex"
)

//读出给定长度的加密的已经Hex过的字符串(结果字符串就是那么长)
func MustCryptoRandToHex(length int) string {
	readLen := length/2 + length%2
	buf := make([]byte, length+length%2)
	_, err := rand.Read(buf[:readLen])
	if err != nil {
		panic(err)
	}
	hex.Encode(buf, buf[:readLen])
	return string(buf[:length])
}
