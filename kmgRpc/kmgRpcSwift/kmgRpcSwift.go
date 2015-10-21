package kmgRpcSwift

import (
	"github.com/bronze1man/kmg/kmgCache"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgGoSource/kmgFormat"
	"path/filepath"
	"strings"
	"github.com/bronze1man/kmg/kmgFileToXcode"
)

type GenerateRequest struct {
	ObjectPkgPath   string
	ObjectName      string
	ObjectIsPointer bool
	OutFilePath     string //输出的文件路径,仅用于写入文件  如 /root/xxx/src/RpcDemo.swift
	OutClassName    string // swift的类的名字 如 RpcDemo
	OutProjectName  string
	ApiNameFilterCb func(name string) bool
	NeedSource      bool
}

//生成代码
// 此处只生成swift代码,不生成golang代码.
// 限制: 输出只能有一个参数,
func MustGenerateCode(req *GenerateRequest) {
	config := reflectToTplConfig(req)
	outBs := tplGenerateCode(config)
	outB := kmgFormat.RemoteEmptyLine([]byte(outBs))
	kmgFile.MustWriteFileWithMkdir(req.OutFilePath, outB)
	if req.NeedSource{
		//生成xxx-Bridging-Header.h NSData+Compression.h NSData+Compression.m
		BridgingHeaderContent := `
	//  请将该文件放到根目录的项目名文件下
	#import "NSData+Compression.h"
	`
		path := strings.Split(req.OutFilePath,"/")
		parPath := strings.Join(path[:(len(path)-1)],"/") + "/"
		BridgingHeaderPath := parPath + req.OutProjectName+"-Bridging-Header.h"
		NSDataCompressionMethodPath := parPath + "NSData+Compression.m"
		NSDataCompressionHeadPath := parPath + "NSData+Compression.h"
		InfoListPath := parPath + "Info.plist"
		podfilePath := parPath + "Podfile"
		kmgFile.MustWriteFileWithMkdir(BridgingHeaderPath, []byte(BridgingHeaderContent))
		kmgFile.MustWriteFileWithMkdir(NSDataCompressionMethodPath, []byte(NSDataCompressionMethod()))
		kmgFile.MustWriteFileWithMkdir(NSDataCompressionHeadPath, []byte(NSDataCompressionHead()))
		kmgFile.MustWriteFileWithMkdir(InfoListPath, []byte(InfoList()))
		kmgFile.MustWriteFileWithMkdir(podfilePath, []byte(Podfile(req.OutProjectName)))
	}
}

type tplConfig struct {
	ClassName      string        //类名称 如 RpcDemo
	InnerClassList []*InnerClass //里面包含的类的类型定义的名称 包括rpc辅助类,如 xxxRequest 和 golang里面用户定义的struct.
	innerClassMap  map[string]*InnerClass
	ApiList        []Api //api列表 包括所有大写开头额Api名称
}

func (config *tplConfig) addInnerClass(class *InnerClass) {
	_, ok := config.innerClassMap[class.Name]
	if ok {
		panic("InnerClass name repeat [" + class.Name + "]")
	}
	config.innerClassMap[class.Name] = class
	config.InnerClassList = append(config.InnerClassList, class)
}

type Api struct {
	Name             string         //在这个系统里面的名字
	InArgsList       []NameTypePair //输入变量列表
	OutTypeString    string         // 有可能是void
	OutTypeFieldName string         // 输出的那个变量的在response里面的名字,如果没有表示直接返回response
}

func (api *Api) getClientFuncInParameter() string {
	outputList := []string{}
	for _, arg := range api.InArgsList {
		outputList = append(outputList, arg.Name+": "+arg.TypeStr)
	}
	return strings.Join(outputList, ",")
}

type NameTypePair struct {
	Name    string
	TypeStr string
}

type InnerClass struct {
	Name      string //此处只有一个层次的名称,如果原先有package会被直接灭掉.
	FieldList []NameTypePair
	IsPublic  bool
}

// 使用缓存 生成代码
func MustGenerateCodeWithCache(req *GenerateRequest) {
	pkgFilePath := kmgConfig.DefaultEnv().PathInProject(filepath.Join("src", req.ObjectPkgPath))
	kmgCache.MustMd5FileChangeCache("kmgRpc_"+req.OutFilePath, []string{req.OutFilePath, pkgFilePath}, func() {
		MustGenerateCode(req)
	})
}
