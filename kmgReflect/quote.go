package kmgReflect

import "encoding/json"

//反序列化 golang的双引号 "encoding/json" -> encoding/json
func UnquoteGolangDoubleQuote(in string) (out string, err error) {
	//先使用json的吧.
	err = json.Unmarshal([]byte(in), &out)
	return
}
