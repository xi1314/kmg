package kmgBase64

import "encoding/base64"

func MustStdBase64DecodeString(s string) (out []byte) {
	out, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return
}
