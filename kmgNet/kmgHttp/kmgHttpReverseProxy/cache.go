package kmgHttpReverseProxy

import (
	"github.com/bronze1man/kmg/kmgCache"
	"github.com/bronze1man/kmg/kmgCompress"
	"github.com/bronze1man/kmg/kmgCrypto"
	"github.com/bronze1man/kmg/kmgLog"
	"github.com/bronze1man/kmg/kmgNet"
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	"github.com/bronze1man/kmg/kmgStrings"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func NewProxyHttpHandlerFromDialer(dialer kmgNet.Dialer) func(rw http.ResponseWriter, req *http.Request) {
	/*
		return &httputil.ReverseProxy{
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
	*/
	transport := &http.Transport{
		Dial: dialer,
	}
	httpClient := &http.Client{
		Transport: transport,
	}
	ttlCacher := kmgCache.NewMemoryTtlCacheV2()
	return func(w http.ResponseWriter, req *http.Request) {
		// 如果是get请求,进行1秒钟的缓存.
		// TODO 只有图片,js,css有缓存,其他的都没有缓存.
		if req.Method == "GET" {
			req.RequestURI = ""
			if req.Proto == "" {
				req.Proto = "HTTP/1.1"
			}
			if req.URL.Scheme == "" {
				req.URL.Scheme = "http"
			}
			req.URL.Host = req.Host
			ustring := req.URL.String()
			if kmgStrings.IsInSlice([]string{".css", ".js", ".gif", ".woff2", ".png", ".jpeg", ".jpg", ".ico"}, filepath.Ext(ustring)) {
				getCache(ttlCacher, httpClient, ustring, w, req)
				return
			}
		}
		req1, err := kmgHttp.HttpRequestClone(req)
		if err != nil {
			http.Error(w, "Wrong format", 400)
			kmgLog.Log("InfoServerError", err.Error(), kmgHttp.NewLogStruct(req))
			return
		}
		err = httpProxyToWrite(transport, w, req1)
		if err != nil {
			http.Error(w, "gateway error", 503)
			kmgLog.Log("InfoServerError", err.Error(), kmgHttp.NewLogStruct(req))
			return
		}
		return
	}
}

type getCacheResponse struct {
	Content        []byte
	ContentType    string
	DeflateContent []byte
	ETag           string
}

func httpProxyToWrite(transport http.RoundTripper, w http.ResponseWriter, req *http.Request) (err error) {
	req.RequestURI = ""
	if req.Proto == "" {
		req.Proto = "HTTP/1.1"
	}
	if req.URL.Scheme == "" {
		req.URL.Scheme = "http"
	}
	req.URL.Host = req.Host
	resp, err := transport.RoundTrip(req) //使用这个避免跟进 redirect
	if err != nil {
		return err
	}
	kmgHttp.HttpResponseToWrite(resp, w)
	return nil
}

func getCache(ttlCacher *kmgCache.MemoryTtlCacheV2, httpClient *http.Client, ustring string, w http.ResponseWriter, req *http.Request) {
	entryi, err := ttlCacher.Do(ustring, func() (value interface{}, ttl time.Duration, err error) {
		resp, err := httpClient.Get(ustring)
		if err != nil {
			return
		}
		b, err := kmgHttp.ResponseReadAllBody(resp)
		if err != nil {
			return
		}
		entry := getCacheResponse{
			ETag:           kmgCrypto.Md5Hex(b),
			Content:        b,
			ContentType:    resp.Header.Get("Content-Type"),
			DeflateContent: kmgCompress.FlateMustCompress(b),
		}
		return entry, time.Hour, nil
	})
	if err != nil {
		http.Error(w, "gateway error", 503)
		kmgLog.Log("InfoServerError", err.Error(), kmgHttp.NewLogStruct(req))
		return
	}
	entry := entryi.(getCacheResponse)
	if strings.HasPrefix(entry.ContentType, "image/") {
		w.Header().Set("Cache-Control", "max-age=864000")
	}
	if strings.HasPrefix(entry.ContentType, "application/javascript") {
		w.Header().Set("Cache-Control", "max-age=40000")
	}
	if strings.HasPrefix(entry.ContentType, "text/css") {
		w.Header().Set("Cache-Control", "max-age=40000")
	}
	etag := entry.ETag
	if req.Header.Get("If-None-Match") == etag {
		w.WriteHeader(304)
		return
	}
	w.Header().Set("Content-Type", entry.ContentType)
	w.Header().Set("ETag", etag)
	if strings.Contains(req.Header.Get("Accept-Encoding"), "deflate") {
		w.Header().Set("Content-Encoding", "deflate")
		b := entry.DeflateContent
		w.Header().Set("Content-Length", strconv.Itoa(len(b)))
		w.Write(b)
	} else {
		w.Header().Set("Content-Length", strconv.Itoa(len(entry.Content)))
		w.Write(entry.Content)
	}
	return
}
