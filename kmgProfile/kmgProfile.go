package kmgProfile

import (
	"expvar"
	"fmt"
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	"net/http"
	"net/http/pprof"
	"path/filepath"
	"runtime/debug"
)

// 可以使用PrefixPath提高安全性
// 如果你不关心安全性 可以使用 / 以便使用golang的默认值. (kmgHttp.ClearHttpDefaultServer() 避免多次注册)
func RegisterProfile(prefixPath string) {
	http.Handle(filepath.Join(prefixPath, "/debug/pprof/profile"), http.HandlerFunc(pprof.Profile))
	http.Handle(filepath.Join(prefixPath, "/debug/pprof/symbol"), http.HandlerFunc(pprof.Symbol))
	http.Handle(filepath.Join(prefixPath, "/debug/pprof/heap"), http.HandlerFunc(heap))
	http.Handle(filepath.Join(prefixPath, "/debug/pprof/block"), pprof.Handler("block"))
	http.Handle(filepath.Join(prefixPath, "/debug/pprof/goroutine"), pprof.Handler("goroutine"))
	http.Handle(filepath.Join(prefixPath, "/debug/pprof/threadcreate"), pprof.Handler("threadcreate"))

	http.Handle(filepath.Join(prefixPath, "/debug/pprof/"), http.HandlerFunc(pprof.Index))
	http.Handle(filepath.Join(prefixPath, "/debug/vars"), http.HandlerFunc(ExpvarHandler))
	http.Handle(filepath.Join(prefixPath, "/debug/gc"), http.HandlerFunc(GcHandler))
}

// 暂时使用默认http的handler
func StartProfileOnAddr(prefixPath string, profileAddr string) {
	kmgHttp.ClearHttpDefaultServer()
	RegisterProfile("/48qcA6SYYyGGXg/")
	go func() {
		err := http.ListenAndServe(profileAddr, nil)
		if err != nil {
			panic(err)
		}
	}()
}

// Replicated from expvar.go as not public.
// TODO 自行实现一个,灭掉自带的几个变量
func ExpvarHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, "{\n")
	first := true
	expvar.Do(func(kv expvar.KeyValue) {
		if !first {
			fmt.Fprintf(w, ",\n")
		}
		first = false
		fmt.Fprintf(w, "%q: %s", kv.Key, kv.Value)
	})
	fmt.Fprintf(w, "\n}\n")
}

func GcHandler(w http.ResponseWriter, r *http.Request) {
	debug.FreeOSMemory()
	w.Header().Set("Content-Type", "application/text; charset=utf-8")
	w.Write([]byte("SUCCESS"))
}

func heap(w http.ResponseWriter, r *http.Request) {
	debug.FreeOSMemory()
	pprof.Handler("heap").ServeHTTP(w, r)
}
