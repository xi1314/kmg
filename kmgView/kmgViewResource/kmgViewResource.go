package kmgViewResource

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgCrypto"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/third/kmgQiniu"
	"path/filepath"
	"sort"
	"strings"
)

type ResourceUploadRequest struct {
	ResourceList  []string //传入一堆资源目录的列表,然后传到七牛上.
	Qiniu         *kmgQiniu.Context
	QiniuPrefix   string
	OutGoFilePath string
	FuncPrefix    string
}

func ResourceUpload(req *ResourceUploadRequest) {
	req.QiniuPrefix = strings.Trim(req.QiniuPrefix, "/")
	fileList := kmgFile.MustGetAllFileFromPathList(req.ResourceList)
	//1 css
	cssContent := []byte{}
	jsPathList := []string{}
	for _, file := range fileList {
		if kmgFile.HasExt(file, ".css") {
			cssContent = append(cssContent, kmgFile.MustReadFile(file)...)
			cssContent = append(cssContent, '\n')
		}
		if kmgFile.HasExt(file, ".js") {
			jsPathList = append(jsPathList, file)
		}
	}
	sort.Strings(jsPathList)
	jsContent := []byte{}
	for _, file := range jsPathList {
		jsContent = append(jsContent, kmgFile.MustReadFile(file)...)
		jsContent = append(jsContent, '\n')
	}
	cssFileName := kmgCrypto.Md5Hex(cssContent) + ".css"
	jsFileName := kmgCrypto.Md5Hex(jsContent) + ".js"
	writeFilePath := []kmgFile.PathAndContentPair{
		{
			Path:    cssFileName,
			Content: cssContent,
		},
		{
			Path:    jsFileName,
			Content: jsContent,
		},
	}

	for _, pair := range writeFilePath {
		req.Qiniu.MustUploadFromBytes(req.QiniuPrefix+"/"+pair.Path, pair.Content)
	}
	packageName := filepath.Base(filepath.Dir(req.OutGoFilePath))
	if req.FuncPrefix == "" {
		req.FuncPrefix = "getManaged"
	}
	outGoContent := []byte(`package ` + packageName + `
func ` + req.FuncPrefix + `JsUrl()string{
	return ` + fmt.Sprintf("%#v", req.Qiniu.GetSchemeAndDomain()+"/"+req.QiniuPrefix+"/"+jsFileName) + `
}
func ` + req.FuncPrefix + `CssUrl()string{
	return ` + fmt.Sprintf("%#v", req.Qiniu.GetSchemeAndDomain()+"/"+req.QiniuPrefix+"/"+cssFileName) + `
}
`)
	kmgFile.MustWriteFile(req.OutGoFilePath, outGoContent)
}
