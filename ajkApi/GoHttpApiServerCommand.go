package ajkApi

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"net/http"

	//"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgConsole"
	"github.com/bronze1man/kmg/kmgCrypto"
	"github.com/bronze1man/kmg/sessionStore"
	//"github.com/bronze1man/kmg/sessionStore/memcacheProvider"
)

var AdditionHttpHandler []HttpHandlerConfig

type HttpHandlerConfig struct {
	Path    string
	Handler http.Handler
}

//start a golang http api server
type GoHttpApiServerCommand struct {
	http          string
	https         string
	randPort      bool
	isHttps       bool
	tcpListenAddr string
}

func RunGoHttpApiServerCmd() {
	command := &GoHttpApiServerCommand{}
	flag.StringVar(&command.http, "http", ":18080", "listen http port of the server")
	flag.StringVar(&command.https, "https", "", "listen https port of the server")
	flag.BoolVar(&command.randPort, "randPort", false, "if can not listen on default port ,will listen on random port")
	flag.Parse()
	if command.https != "" {
		command.isHttps = true
		command.tcpListenAddr = command.https
	} else {
		command.tcpListenAddr = command.http
	}
	jsonHttpHandler := &JsonHttpHandler{
		ApiManager: DefaultApiManager,
		SessionStoreManager: &sessionStore.Manager{
			Provider: sessionStore.NewMemoryProvider(),
			//Provider: memcacheProvider.New(kmgConfig.DefaultParameter().MemcacheHostList...),
		},
	}
	http.Handle("/api", &HttpApiFilterManager{
		Filters: []HttpApiFilter{
			jsonHttpHandler.Filter,
		},
	})
	http.Handle("/api.deflate", &HttpApiFilterManager{
		Filters: []HttpApiFilter{
			HttpApiDeflateCompressFilter,
			jsonHttpHandler.Filter,
		},
	})
	for _, handlerConfig := range AdditionHttpHandler {
		http.Handle(handlerConfig.Path, handlerConfig.Handler)
	}
	l, err := command.listen()
	kmgConsole.ExitOnErr(err)
	fmt.Printf("Listen on %s\n", l.Addr().String())
	if command.isHttps {
		tlsConfig, err := kmgCrypto.CreateTlsConfig()
		if err != nil {
			kmgConsole.ExitOnErr(fmt.Errorf("fail at kmgTls.CreateTlsConfig,error:%s", err.Error()))
		}
		l = tls.NewListener(l, tlsConfig)
	}
	err = http.Serve(l, nil)
	kmgConsole.ExitOnErr(err)
}

//first try addr,if err happened try random address.
func (command *GoHttpApiServerCommand) listen() (l net.Listener, err error) {
	l, err = net.Listen("tcp", command.tcpListenAddr)
	if err == nil {
		return
	}
	if command.randPort {
		l, err = net.Listen("tcp", ":0")
		return
	}
	return
}
