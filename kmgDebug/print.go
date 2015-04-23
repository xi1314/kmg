package kmgDebug

import (
	"encoding/json"
	"fmt"
)

//提供一个漂亮的调试显示接口
func Println(objList ...interface{}) {
	outList := make([]interface{}, len(objList)+1)
	outList[0] = "[kmgDebug.Println]"
	for i := range objList {
		b, err := json.MarshalIndent(objList[i], "", " ")
		if err != nil {
			outList[i+1] = "[Println]error:" + err.Error()
			continue
		}
		outList[i+1] = string(b)
	}
	fmt.Println(outList...)
}
