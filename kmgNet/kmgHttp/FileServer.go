package kmgHttp

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/bronze1man/kmg/kmgFile"
)

// 向http默认服务器加入一个本地文件或目录
func MustAddFileToHttpPathToDefaultServer(httpPath string, localFilePath string) {
	MustAddFileToHttpPathToServeMux(http.DefaultServeMux, httpPath, localFilePath)
}

func MustAddFileToHttpPathToServeMux(mux *http.ServeMux, httpPath string, localFilePath string) {
	localFilePath, err := kmgFile.Realpath(localFilePath)
	//_, err := os.Stat(localFilePath)
	if err != nil {
		panic(err)
	}
	if !strings.HasPrefix(httpPath, "/") {
		httpPath = "/" + httpPath
	}
	if !strings.HasSuffix(httpPath, "/") {
		httpPath = httpPath + "/"
	}
	mux.HandleFunc(httpPath, CompressHandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		urlPath := req.URL.Path
		relPath := strings.TrimPrefix(urlPath, httpPath)
		filePath := filepath.Join(localFilePath, relPath)
		fi, err := os.Stat(filePath)
		if err != nil {
			http.NotFound(w, req)
			return
		}
		if fi.IsDir() {
			http.NotFound(w, req)
			return
		}
		http.ServeFile(w, req, filePath)
	}))
	/*
		if fi.IsDir() {
			if !strings.HasSuffix(httpPath, "/") {
				httpPath += "/"
			}
			mux.Handle(httpPath, http.StripPrefix(httpPath, http.FileServer(http.Dir(localFilePath))))
		} else {
			mux.HandleFunc(httpPath, func(w http.ResponseWriter, req *http.Request) {
				http.ServeFile(w, req, localFilePath)
			})
		}
	*/
	return
}
