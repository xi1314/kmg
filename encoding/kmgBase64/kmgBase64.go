package kmgBase64

import "encoding/base64"

func MustStdBase64DecodeStringToByte(s string) (out []byte) {
	out, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return
}

func StdBase64DecodeByteToByte(s []byte) (out []byte, err error) {
	out = make([]byte, base64.StdEncoding.DecodedLen(len(s)))
	nw, err := base64.StdEncoding.Decode(out, s)
	if err != nil {
		return nil, err
	}
	return out[:nw], nil
}

func Base64EncodeStringToString(input string) string {
	return base64.URLEncoding.EncodeToString([]byte(input))
}

func Base64EncodeByteToString(input []byte) string {
	return base64.URLEncoding.EncodeToString(input)
}

func MustBase64DecodeStringToString(input string) string {
	output, err := base64.URLEncoding.DecodeString(input)
	if err != nil {
		panic(err)
	}
	return string(output)
}

func Base64DecodeStringToByte(input string) (b []byte, err error) {
	return base64.URLEncoding.DecodeString(input)
}
