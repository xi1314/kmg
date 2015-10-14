package kmgRpc

import (
	"strings"

	"github.com/bronze1man/kmg/kmgCache"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgGoSource/kmgFormat"
	"path/filepath"
)

type GenerateRequest struct {
	//Object               interface{} //需要分析的对象
	ObjectPkgPath        string //TODO 对此处进行封装,解决描述对象问题.
	ObjectName           string
	ObjectIsPointer      bool
	OutFilePath          string //生成的文件路径
	OutPackageImportPath string //生成的package的importPath
	ApiNameFilterCb      func(name string) bool
}

//生成代码
// 会把Object上面所有的公开函数都拿去生成一遍
func MustGenerateCode(req *GenerateRequest) {
	config := reflectToTplConfig(req)
	outBs := tplGenerateCode(config)
	outB := []byte(outBs)
	outB1, err := kmgFormat.Source(outB)
	if err == nil {
		outB = outB1
	}
	kmgFile.MustWriteFileWithMkdir(req.OutFilePath, outB)
	return
}

// 使用缓存 生成代码
func MustGenerateCodeWithCache(req *GenerateRequest) {
	pkgFilePath := kmgConfig.DefaultEnv().PathInProject(filepath.Join("src", req.ObjectPkgPath))
	kmgCache.MustMd5FileChangeCache("kmgRpc_"+req.OutFilePath, []string{req.OutFilePath, pkgFilePath}, func() {
		MustGenerateCode(req)
	})
}

type tplConfig struct {
	OutPackageName string          //生成的package的名字 testPackage
	ObjectName     string          //对象名字	如 Demo
	ObjectTypeStr  string          //对象的类型表示	如 *Demo
	ImportPathMap  map[string]bool //ImportPath列表
	ApiList        []Api           //api列表
}

func (conf *tplConfig) mergeImportPath(importPathList []string) {
	for _, importPath := range importPathList {
		conf.ImportPathMap[importPath] = true
	}
}

type Api struct {
	Name        string                 //在这个系统里面的名字
	InArgsList  []ArgumentNameTypePair //输入变量列表
	OutArgsList []ArgumentNameTypePair //输出变量列表
}

func (api Api) GetOutArgsListWithoutError() []ArgumentNameTypePair {
	out := make([]ArgumentNameTypePair, 0, len(api.OutArgsList))
	for _, pair := range api.OutArgsList {
		if pair.ObjectTypeStr == "error" {
			continue
		}
		out = append(out, pair)
	}
	return out
}
func (api Api) GetOutArgsNameListForAssign() string {
	nameList := []string{}
	for _, pair := range api.OutArgsList {
		nameList = append(nameList, pair.Name)
	}
	return strings.Join(nameList, ",")
}

func (api Api) HasReturnArgument() bool {
	return len(api.OutArgsList) > 0
}

func (api Api) GetClientOutArgument() []ArgumentNameTypePair {
	// 确保客户端的接口的返回变量一定有一个error,用来返回错误.
	for _, pair := range api.OutArgsList {
		if pair.ObjectTypeStr == "error" {
			return api.OutArgsList
		}
	}
	return append(api.OutArgsList, ArgumentNameTypePair{
		Name:          "Err",
		ObjectTypeStr: "error",
	})
}

// TODO 下一个版本不要这个hook了,复杂度太高
func (api Api) IsOutExpendToOneArgument() bool {
	return len(api.OutArgsList) == 2 &&
		api.OutArgsList[0].Name == "Response" &&
		api.OutArgsList[1].ObjectTypeStr == "error"
}

func (api Api) GetClientInArgsList() []ArgumentNameTypePair{
	// 删去特殊输入参数 Ctx *http.Context
	out:=[]ArgumentNameTypePair{}
	for _,pair:=range api.InArgsList{
		if pair.ObjectTypeStr == "*kmgHttp.Context"{
			continue
		}
		out = append(out,pair)
	}
	return out
}

func (api Api) HasHttpContextArgument() bool{
	for _,pair:=range api.InArgsList{
		if pair.ObjectTypeStr == "*kmgHttp.Context"{
			return true
		}
	}
	return false
}

// 服务器端,函数调用的那个括号里面的东西
func (api Api) serverCallArgumentStr() string{
	out:=""
	for _,pair:=range api.InArgsList{
		if pair.ObjectTypeStr == "*kmgHttp.Context"{
			out+="Ctx,"
		}else{
			out+="reqData."+pair.Name+","
		}
	}
	return out
}

type ArgumentNameTypePair struct {
	Name          string
	ObjectTypeStr string
}
