package kmgCrypto

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

// @deprecated
var GenUUIDErrors = errors.New("gen uuid fail")

// @deprecated
func GenUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := rand.Read(uuid)
	if n != len(uuid) || err != nil {
		return "", GenUUIDErrors
	}
	return hex.EncodeToString(uuid), nil
}

// @deprecated
func MustGenUUID() string {
	val, err := GenUUID()
	if err != nil {
		panic(err)
	}
	return val
}
