package kmgControllerRunner

import (
	//"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/bronze1man/kmg/kmgErr"
	"github.com/bronze1man/kmg/kmgLog"
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	"github.com/bronze1man/kmg/kmgReflect"
	"github.com/bronze1man/kmg/kmgTime"
)

// EnterPointApiName 最前面不要加 /?n=
var EnterPointApiName = ""
var controllerObjMap = map[string]func(ctx *kmgHttp.Context){} //key 带点号的完整的类名.
var controllerFuncType = reflect.TypeOf((func(ctx *kmgHttp.Context))(nil))

func clearController() {
	controllerObjMap = map[string]func(ctx *kmgHttp.Context){}
}

//注册controller,请先注册后使用,该函数不能并发调用,不能在http开始处理之后再注册(data race问题)
// 允许重复注册.
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

// 此处可以随意定义Api的名称,用来解决调用者的名称向前兼容问题.
// 普通情况下可以使用 RegisterController,减少信息重复.
func RegisterControllerFunc(name string, f func(ctx *kmgHttp.Context)) {
	controllerObjMap[name] = f
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

/**
Example:
在 kmgControllerRunner.StartServerCommand() 前调用
needAuthFilter 返回 true 表示需要验证身份
下面这个例子是：只有在访问前缀是 "/?n=Admin.User.Info" 时才验证身份
kmgControllerRunner.AddHTTPBasicAuthenticationDispatcher("UserName", "Password", func(ctx *kmgHttp.Context) bool {
	return strings.HasPrefix(ctx.GetRequestUrl(), "/?n=Admin.User.Info")
})
*/
func AddHTTPBasicAuth(username, password string) {
	AddHTTPBasicAuthWithFilter(username, password, nil)
}

func AddHTTPBasicAuthWithFilter(username, password string, needAuthFilter func(ctx *kmgHttp.Context) bool) {
	f := func(ctx *kmgHttp.Context, processorList []HttpProcessor) {
		authName, authPass, ok := ctx.GetRequest().BasicAuth()
		needAuth := true
		if needAuthFilter != nil {
			needAuth = needAuthFilter(ctx)
		}
		if needAuth && (authName != username || authPass != password || !ok) {
			ctx.SetResponseHeader("WWW-Authenticate", `Basic realm="kmgControllerRunner"`)
			ctx.SetResponseCode(401)
			return
		}
		processorList[0](ctx, processorList[1:])
	}
	HttpProcessorList = append([]HttpProcessor{f}, HttpProcessorList...)
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
