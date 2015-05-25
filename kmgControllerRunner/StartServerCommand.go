package kmgControllerRunner

import (
	"flag"
	"fmt"
	"github.com/bronze1man/kmg/kmgConsole"
	"net/http"
)

func StartServerCommand() {
	var addr string
	flag.StringVar(&addr, "l", ":8080", "listen addr")
	flag.Parse()

	http.Handle("/", HttpHandler)
	fmt.Println("start at", addr)
	err := http.ListenAndServe(addr, nil)
	kmgConsole.ExitOnErr(err)
}
