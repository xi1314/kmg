package kmgHttp

import (
	"net/http"
	"net/url"
	"strings"
)

// url明文反向代理.
// uri表示,以此uri作为前缀可以进入这个地方,此处代理不会修改用户的uri
// 例如: 注册 AddUriProxyToDefaultServer("/data/upload/","http://b.com/")
// 用户访问 http://xxx.com/data/upload/745dfc0e7d24b9232a0226116e115967.jpg 此处会去访问 http://b.com/data/upload/745dfc0e7d24b9232a0226116e115967.jpg
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

// url明文反向代理.
// uri 是否以/开头和结束无所谓, targetUrl 是否以/结束无所谓.
// uri表示,以此path作为前缀可以进入这个地方,此处代理可能会修改请求的path
// 例如: 注册 MustAddUriProxyRefToUriToDefaultServer("/data/upload/","http://b.com/")
// 用户访问 http://xxx.com/data/upload/745dfc0e7d24b9232a0226116e115967.jpg 此处会去访问 http://b.com/745dfc0e7d24b9232a0226116e115967.jpg
func MustAddUriProxyRefToUriToDefaultServer(uri, targetUrl string) {
	uri = "/" + strings.Trim(uri, "/")
	targetUrlObj, err := url.Parse(targetUrl)
	if err != nil {
		panic(err)
	}
	handler := func(w http.ResponseWriter, req *http.Request) {
		proxyReq := MustHttpRequestClone(req)
		proxyReq.Host = targetUrlObj.Host
		proxyReq.URL.Host = targetUrlObj.Host
		proxyReq.URL.Scheme = targetUrlObj.Scheme
		// uri指向单个文件特殊情况 ("/favicon.ico","http://b.com/favicon.ico"
		refPath := strings.TrimPrefix(strings.TrimPrefix(req.URL.Path, uri), "/")
		if refPath == "" {
			proxyReq.URL.Path = strings.TrimSuffix(targetUrlObj.Path, "/")
		} else {
			proxyReq.URL.Path = strings.TrimSuffix(targetUrlObj.Path, "/") + "/" + refPath // 发起请求的时候使用的是 proxyReq.URL.Path 不使用 proxyReq.RequestURI
		}
		proxyReq.URL, err = url.Parse(proxyReq.URL.String())
		if err != nil {
			panic(err)
		}
		proxyReq.RequestURI = proxyReq.URL.RequestURI()
		err = HttpProxyToWriter(w, proxyReq)
		if err != nil {
			panic(err)
		}
	}
	// 目录
	http.DefaultServeMux.HandleFunc(uri+"/", handler)
	// 单文件,由于无法区分是单文件还是目录,只能注册2遍.
	http.DefaultServeMux.HandleFunc(uri, handler)
}
