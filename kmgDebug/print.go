package kmgDebug

import (
	"encoding/json"
	"fmt"
)

//提供一个漂亮的调试显示接口
func Println(objList ...interface{}) {
	outList := make([]interface{}, len(objList))
	for i := range objList {
		b, err := json.MarshalIndent(objList[i], "", " ")
		if err != nil {
			outList[i] = "[Println]error:" + err.Error()
			continue
		}
		outList[i] = string(b)
	}
	fmt.Println(outList...)
}
