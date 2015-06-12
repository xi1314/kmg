package kmgControllerRunner

import (
	"flag"
	"fmt"
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	"net/http"
)

// 以命令行方式运行应用.这个函数会阻塞执行.
// http方式运行 ./MainApp -http=:80
// https 方式 加http跳转运行 (必须绑一个域名,如果有更复杂的需求,此处无法完成)
//    ./MainApp -https=:443 -http=:80 -type=Both -certFile=xxx.crt -keyFile=xxx.key domain=xxx.com
func StartServerCommand() {
	req := ServerRequest{}
	flag.StringVar(&req.HttpAddr, "http", ":8080", "listen addr") //默认值应该不需要root权限.
	flag.StringVar(&req.HttpsAddr, "https", ":443", "")
	flag.StringVar(&req.HttpsCertFilePath, "certFile", "", "")
	flag.StringVar(&req.HttpsKeyFilePath, "keyFile", "", "")
	flag.StringVar(&req.Domain, "domain", "", "")
	t := ""
	flag.StringVar(&t, "type", "Http", "Http or Both")
	flag.Parse()
	req.Type = ServerType(t)
	StartServer(req)
}

type ServerRequest struct {
	HttpAddr          string
	HttpsAddr         string
	HttpsCertFilePath string
	HttpsKeyFilePath  string
	Domain            string
	Type              ServerType
}

type ServerType string

const (
	ServerTypeHttp ServerType = "Http" //只有http
	ServerTypeBoth ServerType = "Both" //http跳转到https
)

func StartServer(sReq ServerRequest) {
	var targetScheme string
	switch sReq.Type {
	case ServerTypeHttp:
		targetScheme = "http"
	case ServerTypeBoth:
		targetScheme = "https"
	default:
		panic("not expect ServerType " + sReq.Type)
	}
	if sReq.Domain != "" {
		oldHttpHandler := HttpHandler
		HttpHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if sReq.Domain != req.Host {
				kmgHttp.Redirect301ToNewHost(w, req, targetScheme, sReq.Domain)
				return
			}
			oldHttpHandler.ServeHTTP(w, req)
		})
	}
	http.Handle("/", HttpHandler)
	switch sReq.Type {
	case ServerTypeHttp:
		fmt.Println("start http at", sReq.HttpAddr)
		err := http.ListenAndServe(sReq.HttpAddr, nil)
		if err != nil {
			panic(err)
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
		err := http.ListenAndServeTLS(sReq.HttpsAddr, sReq.HttpsCertFilePath, sReq.HttpsKeyFilePath,
			nil)
		if err != nil {
			panic(err)
		}
	default:
		panic("impossible execute path")
	}
}
