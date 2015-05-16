package kmgBase64

//@deprecated
func MustStdBase64DecodeString(s string) (out []byte) {
	return MustStdBase64DecodeStringToByte(s)
}

//@deprecated
func StdBase64Decode(s []byte) (out []byte, err error) {
	return StdBase64DecodeByteToByte(s)
}
