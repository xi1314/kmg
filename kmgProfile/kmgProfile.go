package kmgProfile

import (
	"expvar"
	"fmt"
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgSys"
	"github.com/bronze1man/kmg/kmgView"
	"github.com/bronze1man/kmg/kmgView/kmgBootstrap"
	"net/http"
	"net/http/pprof"
	"os"
	"runtime/debug"
	runtimePprof "runtime/pprof"
	"strings"
	"sync"
	"time"
)

// 可以使用PrefixPath提高安全性
// 如果你不关心安全性 可以使用 / 以便使用golang的默认值. (kmgHttp.ClearHttpDefaultServer() 避免多次注册)
func RegisterProfile(prefixPath string) {
	registerProfile(prefixPath,http.DefaultServeMux)
}

func StartProfileOnAddr(prefixPath string, profileAddr string) {

	mux := http.NewServeMux()
	registerProfile(prefixPath, mux)
	go func() {
		err := http.ListenAndServe(profileAddr, mux)
		if err != nil {
			panic(err)
		}
	}()
}

var initOnce sync.Once

func registerProfile(prefixPath string, mux *http.ServeMux) {
	prefixPath = strings.Trim(prefixPath, "/")
	if prefixPath != "" {
		prefixPath = "/" + prefixPath
	}

	mux.HandleFunc(prefixPath+"/pprof/profile", pprof.Profile)
	mux.HandleFunc(prefixPath+"/pprof/symbol", pprof.Symbol)
	mux.HandleFunc(prefixPath+"/pprof/heap", heap)
	mux.Handle(prefixPath+"/pprof/block", pprof.Handler("block"))
	mux.Handle(prefixPath+"/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle(prefixPath+"/pprof/threadcreate", pprof.Handler("threadcreate"))
	mux.HandleFunc(prefixPath+"/pprof/", pprof.Index)
	mux.HandleFunc(prefixPath+"/vars", ExpvarHandler)
	mux.HandleFunc(prefixPath+"/gc", GcHandler)
	mux.HandleFunc(prefixPath+"/", Index)

	initOnce.Do(func() {
		gStartTime = time.Now()
		expvar.Publish("startTime", expvar.Func(startTime))
		expvar.Publish("uptime", expvar.Func(uptime))
	})
}

// Replicated from expvar.go as not public.
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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("SUCCESS"))
}

func heap(w http.ResponseWriter, r *http.Request) {
	debug.FreeOSMemory()
	pprof.Handler("heap").ServeHTTP(w, r)
}

func Index(w http.ResponseWriter, r *http.Request) {
	content := kmgBootstrap.Table{}
	for _, url := range []string{
		"gc",
		"vars",
		//"pprof/profile?debug=1",
		"pprof/heap?debug=1",
		"pprof/threadcreate?debug=1",
		"pprof/goroutine?debug=1",
		"pprof/goroutine?debug=2",
	} {
		content.DataList = append(content.DataList, []kmgView.HtmlRenderer{
			kmgBootstrap.A{Href: url, Title: url},
		})
	}
	w.Write([]byte(kmgBootstrap.NewWrap("debug page", content).HtmlRender()))
}

var gStartTime time.Time

func startTime() interface{} {
	return gStartTime.String()
}
func uptime() interface{} {
	return time.Since(gStartTime).String()
}

func CpuProfile(funcer func()) {
	selfPath, err := kmgSys.GetCurrentExecutePath()
	if err != nil {
		panic(err)
	}
	tmpPath := kmgFile.NewTmpFilePath()
	f, err := os.Create(tmpPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	runtimePprof.StartCPUProfile(f)
	defer runtimePprof.StopCPUProfile()
	funcer()
	runtimePprof.StopCPUProfile()
	f.Close()
	kmgCmd.CmdSlice([]string{"go", "tool", "pprof", "-top", "-cum", selfPath, tmpPath}).MustRun()
	kmgCmd.CmdSlice([]string{"go", "tool", "pprof", "-top", "-cum", "-lines", selfPath, tmpPath}).MustRun()
}
