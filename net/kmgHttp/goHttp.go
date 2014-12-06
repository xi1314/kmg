package kmgHttp

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
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

func NewHttpsCertNotCheckClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

//一个不验证证书,并且使用http代理的http客户端
// httpProxy look like http://127.0.0.1:9876
func MustNewTestClientWithHttpProxy(httpProxy string) *http.Client {
	var Proxy func(*http.Request) (*url.URL, error)
	if httpProxy != "" {
		u, err := url.Parse(httpProxy)
		if err != nil {
			panic(err)
		}
		Proxy = http.ProxyURL(u)
	}
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			Proxy: Proxy,
		},
	}
}
