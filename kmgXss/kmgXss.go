package kmgXss

import (
	"encoding/json"
	"html/template"
)

func Urlv(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' || c == '-' || c == '.' {
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
