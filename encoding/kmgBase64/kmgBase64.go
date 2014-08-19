package kmgBase64

import "encoding/base64"

func MustStdBase64DecodeString(s string) (out []byte) {
	out, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return
}

func StdBase64Decode(s []byte) (out []byte, err error) {
	out = make([]byte, base64.StdEncoding.DecodedLen(len(s)))
	nw, err := base64.StdEncoding.Decode(out, s)
	if err != nil {
		return nil, err
	}
	return out[:nw], nil
}
