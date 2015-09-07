package kmgHttp

import (
	"flag"
	"fmt"
	"github.com/bronze1man/kmg/kmgConsole"
	"github.com/bronze1man/kmg/kmgCrypto"
	"github.com/bronze1man/kmg/kmgFile"
	"net/http"
	"os"
	"path/filepath"
)

func AddCommandList() {
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "FileHttpServer",
		Runner: runFileHttpServer,
	})
	kmgConsole.AddCommandWithName("FileHttpGet", func() {
		requestUrl := ""
		key := ""
		flag.StringVar(&requestUrl, "url", "", "")
		flag.StringVar(&key, "key", "", "crypto key use to decrypt respond")
		flag.Parse()
		if requestUrl == "" {
			panic("FileHttpGet -url http://xxx")
		}
		b := MustUrlGetContent(requestUrl)
		if key == "" {
			fmt.Println(string(b))
		} else {
			b, err := kmgCrypto.CompressAndEncryptBytesDecodeV2(kmgCrypto.Get32PskFromString(key), b)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(b))
		}
	})
}

func runFileHttpServer() {
	listenAddr := ""
	path := ""
	key := ""
	flag.StringVar(&listenAddr, "l", ":80", "listen address")
	flag.StringVar(&path, "path", "", "root path of the file server")
	flag.StringVar(&key, "key", "", "crypto key use to encrypt all request of this server")
	flag.Parse()
	var err error
	if path == "" {
		path, err = os.Getwd()
		if err != nil {
			fmt.Printf("os.Getwd() fail %s", err)
			return
		}
	}
	if key == "" {
		http.Handle("/", http.FileServer(http.Dir(path)))
	}
	if key != "" {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			realPath := filepath.Join(path, r.URL.Path)
			if !kmgFile.MustFileExist(realPath) {
				w.Write([]byte("File Not Exist"))
				return
			}
			w.Write(
				kmgCrypto.CompressAndEncryptBytesEncode(
					kmgCrypto.Get32PskFromString(key),
					kmgFile.MustReadFile(realPath),
				),
			)
		})
	}
	fmt.Println("start server at", listenAddr)
	err = http.ListenAndServe(listenAddr, nil)
	if err != nil {
		fmt.Printf("http.ListenAndServe() fail %s", err)
		return
	}
	return
}
