package kmgBase32

import "encoding/base32"

//使用标准编码的base32协议
// 字母都是大写
// 有=这个特殊字符
// 包含I和1,0和o之类的东西,人类识别有些困难
func StdBase32EncodeStringToString(input string) string {
	return base32.StdEncoding.EncodeToString([]byte(input))
}
