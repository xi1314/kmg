package kmgHttp

import (
	"flag"
	"fmt"
	"github.com/bronze1man/kmg/kmgConsole"
	"net/http"
	"os"
)

func AddCommandList() {
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "FileHttpServer",
		Runner: runFileHttpServer,
	})
}

func runFileHttpServer() {
	listenAddr := ""
	path := ""
	flag.StringVar(&listenAddr, "l", ":80", "listen address")
	flag.StringVar(&path, "path", "", "root path of the file server")
	flag.Parse()
	var err error
	if path == "" {
		path, err = os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "os.Getwd() fail %s", err)
			return
		}
	}
	http.Handle("/", http.FileServer(http.Dir(path)))
	fmt.Println("start server at", listenAddr)
	err = http.ListenAndServe(listenAddr, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "os.Getwd() fail %s", err)
		return
	}
	return
}
