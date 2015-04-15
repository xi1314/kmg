package kmgControllerRunner

import (
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	"net/http"
	"reflect"
	"strings"
)

var EnterPointApiName = ""
var controllerObjMap = map[string]func(ctx *kmgHttp.Context){} //key 带点号的完整的类名.
var controllerFuncType = reflect.TypeOf((func(ctx *kmgHttp.Context))(nil))

func clearController() {
	controllerObjMap = map[string]func(ctx *kmgHttp.Context){}
}

//注册controller
func RegisterController(obj interface{}) {
	v := reflect.ValueOf(obj)
	t := v.Type()
	objName := t.PkgPath() + "." + t.Name()
	objName = strings.Replace(objName, "/", ".", -1)
	for i := 0; i < t.NumMethod(); i++ {
		mv := v.Method(i)
		mvt := mv.Type()
		if mvt.AssignableTo(controllerFuncType) {
			name := objName + "." + t.Method(i).Name
			controllerObjMap[name] = mv.Interface().(func(ctx *kmgHttp.Context))
		}
	}
}

var HttpHandler = http.HandlerFunc(HttpHandlerFunc)

//httpHandler
func HttpHandlerFunc(w http.ResponseWriter, req *http.Request) {
	ctx := &kmgHttp.Context{W: w, Req: req}
	apiName := ctx.InStr("n")
	if apiName == "" && EnterPointApiName != "" {
		apiName = EnterPointApiName
	}
	apiFunc, ok := controllerObjMap[apiName]
	if !ok {
		ctx.WriteString("api not found")
		return
	}
	apiFunc(ctx)
}
