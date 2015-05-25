package kmgBase64

import "encoding/base64"

//@deprecated
func MustStdBase64DecodeString(s string) (out []byte) {
	return MustStdBase64DecodeStringToByte(s)
}

//@deprecated
func StdBase64Decode(s []byte) (out []byte, err error) {
	return StdBase64DecodeByteToByte(s)
}

// @deprecated
// use Base64EncodeStringToString instead
func MustBase64EncodeStringToString(input string) string {
	return base64.URLEncoding.EncodeToString([]byte(input))
}
