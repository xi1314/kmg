package kmgHttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bronze1man/kmg/encoding/kmgBase64"
	"github.com/bronze1man/kmg/encoding/kmgJson"
	"github.com/bronze1man/kmg/kmgCrypto"
	"github.com/bronze1man/kmg/kmgErr"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strconv"
)

//该对象上的方法不应该被并发调用.
type Context struct {
	Method           string
	RequestUrl       string
	Request          map[string]string
	RequestFile      map[string]*multipart.FileHeader
	Response         string
	ResponseFileName string
	ResponseFile     *bytes.Buffer
	RedirectUrl      string
	ResponseCode     int
	Req              *http.Request
	httpHeader       map[string]string
	sessionMap       map[string]string
}

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

var SessionKey = kmgBase64.MustStdBase64DecodeStringToByte("JOr7fL1TBkU/VqatYYc0D2wERVNUoECzM78HYWaJhIE=")

func NewContextFromHttp(w http.ResponseWriter, req *http.Request) *Context {
	context := &Context{
		Method:      req.Method,
		Request:     map[string]string{},
		RequestFile: map[string]*multipart.FileHeader{},
		RequestUrl:  req.URL.String(),
		//Session:      kmgSession.GetSession(w, req),
		ResponseCode: 200,
		Req:          req,
	}
	err := req.ParseForm()
	if err != nil {
		panic(err)
	}
	for key, value := range req.Form {
		context.Request[key] = value[0] //TODO 这里没有处理同一个 key 多个 value 的情况
	}
	contentType := req.Header.Get("Content-Type")
	if contentType == "" {
		return context
	}
	contentType, _, err = mime.ParseMediaType(contentType)
	if err != nil {
		panic(err)
	}
	if contentType != "multipart/form-data" {
		return context
	}
	err = req.ParseMultipartForm(defaultMaxMemory)
	if err != nil {
		panic(err)
	}
	for key, value := range req.MultipartForm.File {
		context.RequestFile[key] = value[0]
	}
	for key, value := range req.MultipartForm.Value {
		context.Request[key] = value[0]
	}
	return context
}

//返回一个新的测试上下文,这个上下文的所有参数都是空的
func NewTestContext() *Context {
	return &Context{
		RequestUrl:   "/testContext",
		Request:      map[string]string{},
		RequestFile:  map[string]*multipart.FileHeader{},
		ResponseCode: 200,
		sessionMap:   map[string]string{},
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
	return c.Request[key]
}

func (c *Context) InStrDefault(key string, def string) string {
	out := c.Request[key]
	if out == "" {
		return def
	}
	return out
}

func (c *Context) MustPost() {
	if !c.IsPost() {
		panic(errors.New("Need post"))
	}
}

func (c *Context) IsGet() bool {
	return c.Method == "GET"
}
func (c *Context) IsPost() bool {
	return c.Method == "POST"
}

func (c *Context) MustInNum(key string) int {
	s := c.InNum(key)
	if s == 0 {
		panic(fmt.Errorf("Need %s parameter", key))
	}
	return s
}

func (c *Context) InHas(key string) bool {
	return c.Request[key] != ""
}

func (c *Context) MustInStr(key string) string {
	s := c.InStr(key)
	if s == "" {
		panic(fmt.Errorf("Need %s parameter", key))
	}
	return s
}

func (c *Context) MustInJson(key string, obj interface{}) {
	s := c.MustInStr(key)
	err := json.Unmarshal([]byte(s), obj)
	if err != nil {
		panic(err)
	}
	return
}

func (c *Context) InFile(key string) *multipart.FileHeader {
	if file, ok := c.RequestFile[key]; ok {
		return file
	}
	return &multipart.FileHeader{}
}

func (c *Context) SetInStr(key string, value string) *Context {
	c.Request[key] = value
	return c
}

func (c *Context) SetRequest(data map[string]string) *Context {
	c.Request = data
	return c
}

func (c *Context) SetPost() *Context {
	c.Method = "POST"
	return c
}

func (c *Context) sessionInit() {
	if c.sessionMap != nil {
		return
	}
	cookie, err := c.Req.Cookie("kmgSession")
	if err != nil {
		kmgErr.LogErrorWithStack(err)
		c.sessionMap = map[string]string{}
		//没有Cooke
		return
	}
	output, err := kmgBase64.Base64DecodeStringToByte(cookie.Value)
	if err != nil {
		kmgErr.LogErrorWithStack(err)
		c.sessionMap = map[string]string{}
		return
	}
	output, err = kmgCrypto.DecryptV2(SessionKey, output)
	if err != nil {
		kmgErr.LogErrorWithStack(err)
		c.sessionMap = map[string]string{}
		return
	}
	err = json.Unmarshal(output, &c.sessionMap)
	if err != nil {
		kmgErr.LogErrorWithStack(err)
		c.sessionMap = map[string]string{}
		return
	}
}

//向Session里面设置一个字符串
func (c *Context) SessionSetStr(key string, value string) *Context {
	c.sessionInit()
	c.sessionMap[key] = value
	return c
}

//从Session里面获取一个字符串
func (c *Context) SessionGetStr(key string) string {
	c.sessionInit()
	return c.sessionMap[key]
}

func (c *Context) SessionSetJson(key string, value interface{}) *Context {
	json, err := json.Marshal(value)
	if err != nil {
		panic(err) //不能Marshal一定是代码的问题
	}
	c.SessionSetStr(key, string(json))
	return c
}

func (c *Context) SessionGetJson(key string, obj interface{}) (err error) {
	out := c.SessionGetStr(key)
	if out == "" {
		return errors.New("Session Empty")
	}
	err = json.Unmarshal([]byte(out), obj)
	return err
}

//清除Session里面的内容.
//更换Session的Id.
func (c *Context) SessionClear() *Context {
	c.sessionInit()
	c.sessionMap = map[string]string{}
	return c
}

//仅把Session传递过去的上下文,其他东西都恢复默认值
func (c *Context) NewTestContextWithSession() *Context {
	nc := NewTestContext()
	nc.sessionMap = c.sessionMap
	return nc
}

func (c *Context) Redirect(url string) {
	c.RedirectUrl = url
	c.ResponseCode = 302
}

func (c *Context) NotFound(msg string) {
	c.Response = msg
	c.ResponseCode = 404
}

func (c *Context) Error(err error) {
	c.Response += err.Error()
}

func (c *Context) WriteString(s string) {
	c.Response += s
}

func (c *Context) WriteAttachmentFile(file *bytes.Buffer, fileName string) {
	c.ResponseFile = file
	c.ResponseFileName = fileName
}

func (c *Context) WriteJson(obj interface{}) {
	json, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	c.Response += string(json)
}

func (c *Context) WriteToResponseWriter(w http.ResponseWriter, req *http.Request) {
	for key, value := range c.httpHeader {
		w.Header().Set(key, value)
	}
	if c.sessionMap != nil {
		http.SetCookie(w, &http.Cookie{
			Name:  "kmgSession",
			Value: kmgBase64.Base64EncodeByteToString(kmgCrypto.EncryptV2(SessionKey, kmgJson.MustMarshal(c.sessionMap))),
		})
	}
	if c.RedirectUrl != "" {
		http.Redirect(w, req, c.RedirectUrl, c.ResponseCode)
		return
	}
	if c.Response != "" {
		w.WriteHeader(c.ResponseCode)
		w.Write([]byte(c.Response))
		return
	}
	if c.ResponseFile == nil {
		return
	}
	w.Header().Set("Content-Disposition", "attachment;filename="+c.ResponseFileName)
	w.WriteHeader(c.ResponseCode)
	_, err := io.Copy(w, c.ResponseFile)
	if err != nil {
		panic(err)
	}
}

func (c *Context) SetHeader(key, value string) {
	if c.httpHeader == nil {
		c.httpHeader = map[string]string{}
	}
	c.httpHeader[key] = value
}

func (c *Context) CurrentUrl() string {
	return c.RequestUrl
}

func (c *Context) GetResponseCode() int {
	return c.ResponseCode
}

func (c *Context) GetResponseContent() []byte {
	return []byte(c.Response)
}

/*
//这个返回类型可能有问题
func (c *Context)InArray(key string)[]string{
    return nil
}
*/
