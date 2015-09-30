package kmgGoLiteral

import (
	"fmt"
	//"bytes"
	//"reflect"
)

// 未完成,请不要调用.
// TODO 完成它
func MustMarshalToString(obj interface{}) (s string) {
	return fmt.Sprintf("%#v", obj) //不可以使用 fmt.Sprintf("%#v",obj) 会导出私有变量.
	//buf:=&bytes.Buffer{}

}

/*
func mustMarshalReflectWithBuffer(v reflect.Value,buf *bytes.Buffer){
	switch v.Kind() {
	case reflect.String:

	case reflect.Slice:
	case reflect.Struct:
	case reflect.Ptr:
	default:
		panic(fmt.Errorf("[mustMarshalReflectWithBuffer] can not handle reflect kind %s",v.Kind()))
	}
}
*/
