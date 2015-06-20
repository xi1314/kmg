package kmgHex

import "encoding/hex"

func EncodeStringToString(s string) string {
	return hex.EncodeToString([]byte(s))
}

func DecodeStringToString(s string) (string, error) {
	b, err := hex.DecodeString(s)
	return string(b), err
}

func MustDecodeStringToString(s string) string {
	b, err := DecodeStringToString(s)
	if err != nil {
		panic(err)
	}
	return b
}
