package kmgHttp

import (
	"net/http"
	"strconv"
)

//该对象上的方法不应该被并发调用.
type HttpContext struct {
	Req *http.Request
	W   http.ResponseWriter
}

//根据key返回输入参数,包括post和url的query的数据,如果没有,或者不是整数返回0 返回类型为int
func (c *HttpContext) InNum(key string) int {
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
func (c *HttpContext) InStr(key string) string {
	return c.Req.FormValue(key)
}

//TODO 如何处理错误
func (c *HttpContext) MustPost() {

}

func (c *HttpContext) IsGet() bool {
	return c.Req.Method == "GET"
}
func (c *HttpContext) IsPost() bool {
	return c.Req.Method == "POST"
}

//TODO 如何处理错误
func (c *HttpContext) MustInNum(key string) int {
	return 0
}

//TODO 如何处理错误
func (c *HttpContext) MustInStr(key string) string {
	return ""
}

func (c *HttpContext) Redirect(url string) {
	return http.Redirect(c.W, c.Req, url, 302)
}

func (c *HttpContext) WriteString(s string) (int, error) {
	return c.W.Write([]byte(s))
}

/*
目前用的比较少
func (c *HttpContext)InHas(key string)bool{
    return false
}

//这个返回类型可能有问题
func (c *HttpContext)InArray(key string)[]string{
    return nil
}
*/
