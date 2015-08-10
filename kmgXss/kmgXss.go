package kmgXss

import (
	"encoding/json"
	"fmt"
	"html/template"
)

//在模板上面经常需要直接传入int等数据类型.
func Urlv(obj interface{}) string {
	s := fmt.Sprint(obj)
	// 系统的QueryEscape版本 存在bug,主要是' '变'+'这个坑,
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' || c == '-' || c == '.' || c == '_' {
			out = append(out, byte(c))
		} else {
			out = append(out, '%', "0123456789ABCDEF"[c>>4], "0123456789ABCDEF"[c&15])
		}
	}
	return string(out)
}

func H(s string) string {
	return template.HTMLEscapeString(s)
}

func Jsonv(obj interface{}) string {
	out, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return string(out)
}
