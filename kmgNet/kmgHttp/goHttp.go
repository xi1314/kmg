package kmgHttp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"crypto/tls"
	"github.com/bronze1man/kmg/kmgNet"
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

func MustUrlGetContent(url string) (b []byte) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	b, err = ResponseReadAllBody(resp)
	if err != nil {
		panic(err)
	}
	return b
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
	resp, err = http.DefaultTransport.RoundTrip(req) //使用这个避免跟进 redirect
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

//清空默认的Http服务器的路径表
func ClearHttpDefaultServer() {
	http.DefaultServeMux = http.NewServeMux()
}

func Redirect301ToNewHost(w http.ResponseWriter, req *http.Request, scheme string, host string) {
	u := req.URL
	u.Host = host
	if u.Scheme == "" {
		u.Scheme = scheme
	}
	http.Redirect(w, req, u.String(), 301)
}

func MustUrlGetContentProcess(url string) (b []byte) {
	fmt.Print("\nConnnecting\r")
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	printProgress(0, resp.ContentLength, 0, 0)
	defer resp.Body.Close()
	buf := bytes.Buffer{}
	bufBytes := make([]byte, 32*1024)
	lastTime := time.Now()
	lastBytes := 0
	for {
		n, err := resp.Body.Read(bufBytes)
		if n > 0 {
			buf.Write(bufBytes[:n])
			now := time.Now()
			if now.After(lastTime.Add(1 * time.Second)) {
				thisBytes := buf.Len()
				printProgress(int64(thisBytes), resp.ContentLength, now.Sub(lastTime), thisBytes-lastBytes)
				lastTime = now
				lastBytes = thisBytes
			}
		}
		if err == io.EOF {
			fmt.Println()
			return buf.Bytes()
		}
		if err != nil {
			panic(err)
		}
	}
}

func printProgress(get int64, total int64, dur time.Duration, lastBytes int) {
	percent := 0.0
	if total <= 0 {
		percent = 0.0
	} else if total < get {
		percent = 1.0
	} else {
		percent = float64(get) / float64(total)
	}
	showNum := int(percent * 40)
	notShowNum := 40 - showNum
	fmt.Printf("%s%s %.2f%% %s/%s %s     \r",
		strings.Repeat("#", showNum), strings.Repeat(" ", notShowNum), percent*100,
		kmgNet.SizeString(get), kmgNet.SizeString(total), kmgNet.SpeedString(lastBytes, dur))
}

// 异步开启一个http服务器,这个服务器可以使用返回的closer关闭
func MustGoHttpAsyncListenAndServeWithCloser(addr string, handler http.Handler) (closer func() error) {
	srv := &http.Server{Addr: addr, Handler: handler}
	if addr == "" {
		addr = ":80"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	go func() {
		err := srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
		if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			panic(err)
		}
	}()
	return ln.Close
}

func MustGoHttpsAsyncListenAndServeWithCloser(addr string, tlsConfig *tls.Config, handler http.Handler) (closer func() error) {
	srv := &http.Server{Addr: addr, Handler: handler}
	/*
		cert, err := tls.X509KeyPair([]byte(certS), []byte(keyS))
		if err != nil {
			return err
		}
		srv.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			NextProtos:   []string{"http/1.1"},
		}
	*/
	srv.TLSConfig = tlsConfig
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	tlsListener := tls.NewListener(tcpKeepAliveListener{ln.(*net.TCPListener)}, srv.TLSConfig)
	go func() {
		err := srv.Serve(tlsListener)
		if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			panic(err)
		}
	}()
	return tlsListener.Close
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func GoListenAndServeTLSWithCertContent(addr string, certS string, keyS string, handler http.Handler) error {
	srv := &http.Server{Addr: addr, Handler: handler}
	cert, err := tls.X509KeyPair([]byte(certS), []byte(keyS))
	if err != nil {
		return err
	}
	srv.TLSConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"http/1.1"},
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	tlsListener := tls.NewListener(tcpKeepAliveListener{ln.(*net.TCPListener)}, srv.TLSConfig)
	return srv.Serve(tlsListener)
}

// 去掉cert和key,复制粘帖带来的各种格式错误
func FormatHttpsCertOrKey(inS string) string {
	inS = strings.TrimSpace(inS)
	return inS
}
