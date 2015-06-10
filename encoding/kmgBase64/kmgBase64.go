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

// 应用范围比较广的base64转义方案
//   1.输出的字符串里面可能包含下列特殊字符 -_= 这三种
//   2.输出的字符串区分大小写,
//   3.不要放在文件名的地方,mac os 和windows的文件名不区分大小写. 请使用kmgBase32
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
