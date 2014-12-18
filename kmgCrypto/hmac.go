package kmgCrypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"hash/adler32"
	"encoding/binary"
	"bytes"
)

//长度32字节
func CheckSha256MAC(message, messageMAC, key []byte) bool {
	expectedMAC := GetSha256MAC(message, key)
	return hmac.Equal(messageMAC, expectedMAC)
}

//速度比较慢 118MB/s
func GetSha256MAC(message, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	return mac.Sum(nil)
}

//长度4字节 速度非常快 1GB/s
func CheckAdler32(message, messageMAC []byte) bool {
	get:=GetAdler32(message)
	return bytes.Equal(get,messageMAC)
}

func GetAdler32(message []byte) []byte{
	outi:= adler32.Checksum(message)
	out:=make([]byte,4)
	binary.LittleEndian.PutUint32(out,outi)
	return out
}
