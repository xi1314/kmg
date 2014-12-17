package kmgCrypto

import (
	"crypto/hmac"
	"crypto/sha256"
)

func CheckSha256MAC(message, messageMAC, key []byte) bool {
	expectedMAC := GetSha256MAC(message, key)
	return hmac.Equal(messageMAC, expectedMAC)
}

func GetSha256MAC(message, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	return mac.Sum(nil)
}
