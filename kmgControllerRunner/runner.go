package kmgControllerRunner

import (
	"github.com/bronze1man/kmg/kmgErr"
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
	ctx := kmgHttp.NewContextFromHttpRequest(req) //TODO 这里忽略了错误，此处应该如何处理错误
	apiName := ctx.InStr("n")
	if apiName == "" && EnterPointApiName != "" {
		apiName = EnterPointApiName
	}
	apiFunc, ok := controllerObjMap[apiName]
	if !ok {
		ctx.NotFound("api not found")
		ctx.WriteToResponseWriter(w, req)
		return
	}
	err := kmgErr.PanicToError(func() {
		apiFunc(ctx)
	})
	if err != nil {
		ctx.Response = err.Error()
		ctx.ResponseCode = 500
		ctx.WriteToResponseWriter(w, req)
		return
	}
	ctx.WriteToResponseWriter(w, req)
}
