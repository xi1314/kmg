package kmgHttpReverseProxy

import (
	"flag"
	"fmt"
	"github.com/bronze1man/InlProxy/AesCtrConnWrapper"
	"github.com/bronze1man/InlProxy/AesTunnel"
	"github.com/bronze1man/kmg/encoding/kmgYaml"
	"github.com/bronze1man/kmg/kmgConsole"
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	"net"
	"net/http"
	"net/http/httputil"
)

func AddCommandList() {
	kmgConsole.AddCommandWithName("HttpReverseProxy", func() {
		var confPath string
		flag.StringVar(&confPath, "conf", "", "yaml conf path")
		flag.Parse()
		req := &ServerRequest{}
		err := kmgYaml.ReadFile(confPath, req)
		if err != nil {
			fmt.Printf("read config file %s %s\n", confPath, err.Error())
			return
		}
		_, err = RunServer(*req)
		kmgConsole.ExitOnErr(err)
	})
}

type ServerRequest struct {
	HttpAddr  string
	HttpsAddr string
	HttpsCert string // 此处是证书内容
	HttpsKey  string // 此处是密钥内容
	Domain    string //配置域名后,未知域名会跳转到该域名上
	Type      ServerType

	NextAddr string
	AesKey   string //不为空时,使用AesKey作为 Aes加密算法的密钥
}

type ServerType string

const (
	ServerTypeHttp ServerType = "Http" //只有http
	ServerTypeBoth ServerType = "Both" //http跳转到https https提供服务
)

func RunServer(sReq ServerRequest) (closer func() error, err error) {
	dialer := func(network, address string) (net.Conn, error) {
		return net.Dial("tcp", sReq.NextAddr)
	}
	if sReq.AesKey != "" {
		dialer = AesTunnel.NewAesStartDialer(AesCtrConnWrapper.NewKey([]byte(sReq.AesKey)), dialer)
	}
	proxyer := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.RequestURI = ""
			if req.Proto == "" {
				req.Proto = "HTTP/1.1"
			}
			if req.URL.Scheme == "" {
				req.URL.Scheme = "http"
			}
			req.URL.Host = req.Host
		},
		Transport: &http.Transport{
			Dial: dialer,
		},
	}
	handler := proxyer.ServeHTTP
	var targetScheme string
	switch sReq.Type {
	case ServerTypeHttp:
		targetScheme = "http"
	case ServerTypeBoth:
		targetScheme = "https"
	default:
		panic("not expect ServerType [" + string(sReq.Type) + "]")
	}
	if sReq.Domain != "" {
		oldHttpHandler := handler
		handler = func(w http.ResponseWriter, req *http.Request) {
			if sReq.Domain != req.Host {
				kmgHttp.Redirect301ToNewHost(w, req, targetScheme, sReq.Domain)
				return
			}
			oldHttpHandler(w, req)
		}
	}
	switch sReq.Type {
	case ServerTypeHttp:
		fmt.Println("start http at", sReq.HttpAddr)
		err := http.ListenAndServe(sReq.HttpAddr, http.HandlerFunc(handler))
		if err != nil {
			return nil, err
		}
	case ServerTypeBoth:
		fmt.Printf("start http[%s] https[%s]\n", sReq.HttpAddr, sReq.HttpsAddr)
		go func() {
			err := http.ListenAndServe(sReq.HttpAddr, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				kmgHttp.Redirect301ToNewHost(w, req, "https", sReq.Domain)
				return
			}))
			panic(err)
		}()
		err := kmgHttp.GoListenAndServeTLSWithCertContent(sReq.HttpsAddr, sReq.HttpsCert, sReq.HttpsKey, http.HandlerFunc(handler))
		if err != nil {
			return nil, err
		}
	default:
		panic("impossible execute path")
	}
	return nil, nil
}
