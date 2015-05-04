package kmgRand

import (
	"crypto/rand"
	"encoding/hex"
)

func CryptoRandBytes(length int) []byte {
	buf := make([]byte, length)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	return buf
}

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

const alphaNumMap = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func MustCryptoRandToAlphaNum(length int) string {
	return MustCryptoRandFromByteList(length, alphaNumMap)
}

const realableAlphaNumMap = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"

func MustCryptoRandToReadableAlphaNum(length int) string {
	return MustCryptoRandFromByteList(length, realableAlphaNumMap)
}

const numMap = "0123456789"

func MustCryptoRandToNum(length int) string {
	return MustCryptoRandFromByteList(length, numMap)
}

func MustCryptoRandFromByteList(length int, list string) string {
	var bytes = make([]byte, 2*length)
	var outBytes = make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	mapLen := len(list)
	for i := 0; i < length; i++ {
		outBytes[i] = list[(int(bytes[2*i])*256+int(bytes[2*i+1]))%(mapLen)]
	}
	return string(outBytes)
}
