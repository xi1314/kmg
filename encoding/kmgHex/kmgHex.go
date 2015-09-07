package kmgHex

import (
	"encoding/hex"
	"strings"
)

// 生成的字符串是大写
func UpperEncodeBytesToString(b []byte) string {
	return strings.ToUpper(hex.EncodeToString(b))
}

func EncodeBytesToString(b []byte) string {
	return hex.EncodeToString(b)
}

// 生成的字符串是小写
func EncodeStringToString(s string) string {
	return hex.EncodeToString([]byte(s))
}

// 可以解码大写,也可以解码小写
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
func MustDecodeStringToByteArray(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}
