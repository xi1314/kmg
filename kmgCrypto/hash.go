package kmgCrypto

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

func Sha256Hex(data []byte) string {
	out := sha256.Sum256(data)
	return hex.EncodeToString(out[:])
}

func Sha512Hex(data []byte) string {
	out := sha512.Sum512(data)
	return hex.EncodeToString(out[:])
}
func Sha512HexFromString(data string) string {
	out := sha512.Sum512([]byte(data))
	return hex.EncodeToString(out[:])
}

func Md5Hex(data []byte) string {
	hash := md5.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}
