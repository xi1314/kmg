package kmgDebug

import (
	"encoding/json"
	"fmt"
	"github.com/bronze1man/kmg/kmgReflect"
	"os"
	"reflect"
)

//提供一个漂亮的调试显示接口
// TODO 根据不同类型使用不同的显示方法
func Println(objList ...interface{}) {
	s := Sprintln(objList...)
	os.Stdout.WriteString(s)
	return
}

func Sprintln(objList ...interface{}) string {
	outList := make([]interface{}, len(objList)+1)
	outList[0] = "[kmgDebug.Println]"
	for i := range objList {
		if kmgReflect.IsNil(reflect.ValueOf(objList[i])) {
			outList[i+1] = "nil"
			continue
		}
		switch obj := objList[i].(type) {
		case []byte:
			outList[i+1] = fmt.Sprintf("%#v", obj)
		default:
			b, err := json.MarshalIndent(objList[i], "", " ")
			if err != nil {
				outList[i+1] = "[Println]error:" + err.Error()
				continue
			}
			outList[i+1] = string(b)
		}

	}
	return fmt.Sprintln(outList...)
}
