package kmgHttp

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

func NewRequestFromByte(r []byte) (req *http.Request, err error) {
	return http.ReadRequest(bufio.NewReader(bytes.NewReader(r)))
}

//sometimes it is hard to remember how to get response from bytes ...
func NewResponseFromBytes(r []byte) (resp *http.Response, err error) {
	return http.ReadResponse(bufio.NewReader(bytes.NewBuffer(r)), &http.Request{})
}

//sometimes it is hard to remember how to dump response to bytes
func DumpResponseToBytes(resp *http.Response) (b []byte, err error) {
	return httputil.DumpResponse(resp, true)
}

func ResponseReadAllBody(resp *http.Response) (b []byte, err error) {
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func RequestReadAllBody(req *http.Request) (b []byte, err error) {
	defer req.Body.Close()
	return ioutil.ReadAll(req.Body)
}

func MustResponseReadAllBody(resp *http.Response) (b []byte) {
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return b
}

func UrlGetContent(url string) (b []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	return ResponseReadAllBody(resp)
}

func HeaderToString(header http.Header) (s string) {
	buf := &bytes.Buffer{}
	header.Write(buf)
	return string(buf.Bytes())
}

//把request转换成[]byte,并且使body可以被再次读取
func MustRequestToStringCanRead(req *http.Request) (s string) {
	oldBody := req.Body
	defer oldBody.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(body))
	buf := &bytes.Buffer{}
	req.Write(buf)
	req.Body = ioutil.NopCloser(bytes.NewReader(body))
	return string(buf.Bytes())
}

func MustRequestFromString(reqString string) (req *http.Request) {
	req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader([]byte(reqString))))
	if err != nil {
		panic(err)
	}
	return req
}

// 进行http代理,并且代理到writer上面去
// 调用时,请修改req的参数,避免自己调用自己
// 如果出现错误,(对方服务器连不上之类的,不会修改w,会返回一个error
// 不跟踪redirect(跟踪redirect会导致redirect的请求的内容被返回)
func HttpProxyToWriter(w http.ResponseWriter, req *http.Request) (err error) {
	resp, err := HttpRoundTrip(req)
	if err != nil {
		return err
	}
	HttpResponseToWrite(resp, w)
	return
}

func HttpRoundTrip(req *http.Request) (resp *http.Response, err error) {
	req.RequestURI = ""
	if req.Proto == "" {
		req.Proto = "HTTP/1.1"
	}
	if req.URL.Scheme == "" {
		req.URL.Scheme = "http"
	}
	resp, err = http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	return
}

func HttpResponseToWrite(resp *http.Response, w http.ResponseWriter) {
	defer resp.Body.Close()
	for k, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(k, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

//会把request里面的东西全部都读出来(body)
func HttpRequestClone(in *http.Request) (out *http.Request, err error) {
	buf := &bytes.Buffer{}
	err = in.Write(buf)
	if err != nil {
		return
	}
	out, err = http.ReadRequest(bufio.NewReader(buf))
	if err != nil {
		return
	}
	return
}

func MustHttpRequestClone(in *http.Request) *http.Request {
	out, err := HttpRequestClone(in)
	if err != nil {
		panic(err)
	}
	return out
}

func MustAddFileToHttpPathToDefaultServer(httpPath string, localFilePath string) {
	err := AddFileToHttpPathToServeMux(http.DefaultServeMux, httpPath, localFilePath)
	if err != nil {
		panic(err)
	}
}

func AddFileToHttpPathToServeMux(mux *http.ServeMux, httpPath string, localFilePath string) error {
	fi, err := os.Stat(localFilePath)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(httpPath, "/") {
		httpPath = "/" + httpPath
	}
	if fi.IsDir() {
		if !strings.HasSuffix(httpPath, "/") {
			httpPath += "/"
		}
		mux.Handle(httpPath, http.StripPrefix(httpPath, http.FileServer(http.Dir(localFilePath))))
	} else {
		mux.HandleFunc(httpPath, func(w http.ResponseWriter, req *http.Request) {
			http.ServeFile(w, req, localFilePath)
		})
	}
	return nil
}

func AddUriProxyToDefaultServer(uri, targetUrl string) {
	http.DefaultServeMux.HandleFunc(uri, func(w http.ResponseWriter, req *http.Request) {
		proxyReq := MustHttpRequestClone(req)
		target, err := url.Parse(targetUrl)
		if err != nil {
			panic(err)
		}
		proxyReq.Host = target.Host
		proxyReq.URL.Host = target.Host
		proxyReq.URL.Scheme = target.Scheme
		err = HttpProxyToWriter(w, proxyReq)
		if err != nil {
			panic(err)
		}
	})
}
