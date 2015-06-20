package kmgBase64

import (
	"encoding/base64"
	"strings"
)

// urlbase64 并且去掉了=符号
//   1.输出的字符串里面可能包含下列特殊字符 -_ 这三种
//   2.输出的字符串区分大小写,
//   3.不要放在文件名的地方,mac os 和windows的文件名不区分大小写. 请使用kmgBase32
func EncodeByteToStringV2(input []byte) string {
	out := base64.URLEncoding.EncodeToString([]byte(input))
	return strings.Replace(out, "=", "", -1)
}

func DecodeStringToByteV2(input string) (b []byte, err error) {
	if len(input)%4 != 0 {
		input += strings.Repeat("=", 4-len(input)%4)
	}
	return base64.URLEncoding.DecodeString(input)
}
