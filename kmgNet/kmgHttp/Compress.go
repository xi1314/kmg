package kmgHttp

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type ResponseWriterWraper struct {
	io.Writer
	http.ResponseWriter
}

// TODO 这个地方缺少猜类型功能,导致不能直接wrapper到任意Handler上去.(会导致 html Type 错误.)
// TODO 有时候压缩会有负效果,此处应该可以自动判断出来.
func (w ResponseWriterWraper) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func CompressHandlerFunc(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		acceptEncoding := r.Header.Get("Accept-Encoding")
		switch {
		case strings.Contains(acceptEncoding, "deflate"):
			w.Header().Set("Content-Encoding", "deflate")
			gzw, err := flate.NewWriter(w, -1)
			if err != nil {
				panic(err)
			}
			defer gzw.Close()
			gzr := ResponseWriterWraper{Writer: gzw, ResponseWriter: w}
			f(gzr, r)
			return
		case strings.Contains(acceptEncoding, "gzip"):
			w.Header().Set("Content-Encoding", "gzip")
			gzw := gzip.NewWriter(w)
			defer gzw.Close()
			gzr := ResponseWriterWraper{Writer: gzw, ResponseWriter: w}
			f(gzr, r)
			return
		default:
			f(w, r)
		}
	}
}

func CompressHandler(fn http.Handler) http.Handler {
	return http.HandlerFunc(CompressHandlerFunc(fn.ServeHTTP))
}

// a flate(DEFLATE) compress wrap around http request and response,
// !!not handle any http header!!
func HttpHandleCompressFlateWrap(fn http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		oldBody := r.Body
		defer oldBody.Close()
		r.Body = flate.NewReader(oldBody)
		//w.Header().Set("Content-Encoding", "deflate")
		gzw, err := flate.NewWriter(w, -1)
		if err != nil {
			panic(err)
		}
		defer gzw.Close()
		gzr := ResponseWriterWraper{Writer: gzw, ResponseWriter: w}
		fn.ServeHTTP(gzr, r)
	})
}

// a flate(DEFLATE) compress wrap around http request and response,
// !!not handle any http header!!
func HttpHandleCompressGzipWrap(fn http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		oldBody := r.Body
		defer oldBody.Close()
		var err error
		r.Body, err = gzip.NewReader(oldBody)
		if err != nil {
			panic(err)
		}
		//w.Header().Set("Content-Encoding", "gzip")
		gzw := gzip.NewWriter(w)
		defer gzw.Close()
		gzr := ResponseWriterWraper{Writer: gzw, ResponseWriter: w}
		fn.ServeHTTP(gzr, r)
	})
}
