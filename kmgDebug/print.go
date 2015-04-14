package kmgDebug

import (
	"encoding/json"
	"fmt"
)

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
