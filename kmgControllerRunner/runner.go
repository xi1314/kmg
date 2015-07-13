package kmgControllerRunner

import (
	//"fmt"
	"github.com/bronze1man/kmg/kmgErr"
	"github.com/bronze1man/kmg/kmgLog"
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	"github.com/bronze1man/kmg/kmgReflect"
	"github.com/bronze1man/kmg/kmgTime"
	"net/http"
	"reflect"
	"strings"
	"time"
)

var EnterPointApiName = ""
var controllerObjMap = map[string]func(ctx *kmgHttp.Context){} //key 带点号的完整的类名.
var controllerFuncType = reflect.TypeOf((func(ctx *kmgHttp.Context))(nil))

func clearController() {
	controllerObjMap = map[string]func(ctx *kmgHttp.Context){}
}

//注册controller,请先注册后使用,该函数不能并发调用,不能在http开始处理之后再注册(data race问题)
// 不允许重复注册.
func RegisterController(obj interface{}) {
	v := reflect.ValueOf(obj)
	t := v.Type()

	objName := kmgReflect.GetTypeFullName(t)
	objName = strings.Replace(objName, "/", ".", -1)
	for i := 0; i < t.NumMethod(); i++ {
		if t.Method(i).PkgPath != "" {
			// 不是public不需要注册进来.
			continue
		}
		mv := v.Method(i)
		mvt := mv.Type()
		if mvt.AssignableTo(controllerFuncType) {
			name := objName + "." + t.Method(i).Name
			//_, ok := controllerObjMap[name]
			//if ok {
			//	panic(fmt.Errorf("[RegisterController] Repeat register controller name[%s]", name))
			//}
			controllerObjMap[name] = mv.Interface().(func(ctx *kmgHttp.Context))
		}
	}
}

func GetControllerNameList() []string {
	out := []string{}
	for key := range controllerObjMap {
		out = append(out, key)
	}
	return out
}

var HttpHandler = http.HandlerFunc(HttpHandlerFunc)

//httpHandler
func HttpHandlerFunc(w http.ResponseWriter, req *http.Request) {
	ctx := kmgHttp.NewContextFromHttp(w, req) //TODO 这里忽略了错误，此处应该如何处理错误
	HttpProcessorList[0](ctx, HttpProcessorList[1:])
	ctx.WriteToResponseWriter(w, req)
}

func ContextHandle(ctx *kmgHttp.Context) {
	HttpProcessorList[0](ctx, HttpProcessorList[1:])
}

type HttpProcessor func(ctx *kmgHttp.Context, processorList []HttpProcessor)

var HttpProcessorList = []HttpProcessor{
	PanicHandler,
	Dispatcher,
}

func PanicHandler(ctx *kmgHttp.Context, processorList []HttpProcessor) {
	err := kmgErr.PanicToErrorAndLog(func() {
		processorList[0](ctx, processorList[1:])
	})
	if err != nil {
		ctx.Error(err)
		return
	}
	return
}

func Dispatcher(ctx *kmgHttp.Context, processorList []HttpProcessor) {
	apiName := ctx.InStr("n")
	if apiName == "" && EnterPointApiName != "" {
		if ctx.GetRequestUrl() == "/favicon.ico" {
			// 避免网站图标请求,占用大量资源.
			ctx.NotFound("api not found")
			return
		}
		apiName = EnterPointApiName
	}
	apiFunc, ok := controllerObjMap[apiName]
	if !ok {
		ctx.NotFound("api not found")
		return
	}

	apiFunc(ctx)
	return
}

// 默认不用这个,容易搞的测试里面到处都是log.
// TODO 静态文件的log问题
// TODO 尝试搞出更好用的log系统.
func RequestLogger(ctx *kmgHttp.Context, processorList []HttpProcessor) {
	startTime := time.Now()
	processorList[0](ctx, processorList[1:])
	time := time.Since(startTime)
	log := ctx.Log()
	log.ProcessTime = kmgTime.DurationFormat(time)
	kmgLog.Log("Request", log)
}
