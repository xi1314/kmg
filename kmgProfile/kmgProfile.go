package kmgProfile

import (
	"expvar"
	"fmt"
	"net/http"
	"net/http/pprof"
	"path/filepath"
	"runtime/debug"
)

// 可以使用PrefixPath提高安全性
// 如果你不关心安全性 可以使用 / 以便使用golang的默认值.
func RegisterProfile(prefixPath string) {
	http.Handle(filepath.Join(prefixPath, "/debug/pprof/profile"), http.HandlerFunc(pprof.Profile))
	http.Handle(filepath.Join(prefixPath, "/debug/pprof/symbol"), http.HandlerFunc(pprof.Symbol))
	http.Handle(filepath.Join(prefixPath, "/debug/pprof/heap"), pprof.Handler("heap"))
	http.Handle(filepath.Join(prefixPath, "/debug/pprof/block"), pprof.Handler("block"))

	http.Handle(filepath.Join(prefixPath, "/debug/pprof/"), http.HandlerFunc(pprof.Index))
	http.Handle(filepath.Join(prefixPath, "/debug/vars"), http.HandlerFunc(ExpvarHandler))
	http.Handle(filepath.Join(prefixPath, "/debug/gc"), http.HandlerFunc(GcHandler))
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
