package kmgHttp

import (
	"net/http"
	"strconv"
)

//该对象上的方法不应该被并发调用.
type Context struct {
	Req *http.Request
	W   http.ResponseWriter
}

//根据key返回输入参数,包括post和url的query的数据,如果没有,或者不是整数返回0 返回类型为int
func (c *Context) InNum(key string) int {
	v := c.Req.FormValue(key)
	if v == "" {
		return 0
	}
	out, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return out
}

//根据key返回输入参数,包括post和url的query的数据,如果没有返回"" 类型为string
func (c *Context) InStr(key string) string {
	return c.Req.FormValue(key)
}

//TODO 如何处理错误
func (c *Context) MustPost() {

}

func (c *Context) IsGet() bool {
	return c.Req.Method == "GET"
}
func (c *Context) IsPost() bool {
	return c.Req.Method == "POST"
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
	http.Redirect(c.W, c.Req, url, 302)
	return
}

func (c *Context) WriteString(s string) (int, error) {
	return c.W.Write([]byte(s))
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
