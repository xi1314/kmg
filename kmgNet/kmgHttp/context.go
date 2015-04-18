package kmgHttp

import (
	"encoding/json"
	"net/http"
	"strconv"
)

//该对象上的方法不应该被并发调用.
type Context struct {
	Method       string
	Request      map[string]string
	Response     string
	RedirectUrl  string
	ResponseCode int
}

func NewContextFromHttpRequest(req *http.Request) *Context {
	return &Context{
		Method: req.Method,
		Request: func() map[string]string {
			m := map[string]string{}
			err := req.ParseForm()
			if err != nil {
				return m
			}
			for key, value := range req.Form {
				m[key] = value[0] //TODO 这里没有处理同一个 Key 多个 Value 的情况
			}
			return m
		}(),
		ResponseCode: 200,
	}
}

//根据key返回输入参数,包括post和url的query的数据,如果没有,或者不是整数返回0 返回类型为int
func (c *Context) InNum(key string) int {
	value, ok := c.Request[key]
	if !ok {
		return 0
	}
	num, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return num
}

//根据key返回输入参数,包括post和url的query的数据,如果没有返回"" 类型为string
func (c *Context) InStr(key string) string {
	value, ok := c.Request[key]
	if !ok {
		return ""
	}
	return value
}

//TODO 如何处理错误
func (c *Context) MustPost() {
}

func (c *Context) IsGet() bool {
	return c.Method == "GET"
}
func (c *Context) IsPost() bool {
	return c.Method == "POST"
}

//TODO 如何处理错误
func (c *Context) MustInNum(key string) int {
	return 0
}

//TODO 如何处理错误
func (c *Context) MustInStr(key string) string {
	return ""
}

func (c *Context) Redirect(url string) {
	c.RedirectUrl = url
	c.ResponseCode = 302
}

func (c *Context) NotFound(msg string) {
	c.Response = msg
	c.ResponseCode = 404
}

func (c *Context) Error(msg string) {

}

func (c *Context) WriteString(s string) {
	c.Response += s
}

func (c *Context) WriteJson(obj interface{}) {
	json, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	c.Response += string(json)
}

func (ctx *Context) WriteToResponseWriter(w http.ResponseWriter, req *http.Request) {
	if ctx.RedirectUrl != "" {
		http.Redirect(w, req, ctx.RedirectUrl, ctx.ResponseCode)
		return
	}
	w.WriteHeader(ctx.ResponseCode)
	w.Write([]byte(ctx.Response))
}

/*
目前用的比较少
func (c *Context)InHas(key string)bool{
    return false
}

//这个返回类型可能有问题
func (c *Context)InArray(key string)[]string{
    return nil
}
*/
