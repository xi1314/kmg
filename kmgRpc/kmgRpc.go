package kmgRpc

import (
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgGoSource/kmgFormat"
	"strings"
)

type GenerateRequest struct {
	//Object               interface{} //需要分析的对象
	ObjectPkgPath        string //TODO 对此处进行封装,解决描述对象问题.
	ObjectName           string
	ObjectIsPointer      bool
	OutFilePath          string   //生成的文件路径
	OutPackageImportPath string   //生成的package的importPath
	Key                  [32]byte //密钥
}

//生成代码
// 会把Object上面所有的公开函数都拿去生成一遍
func MustGenerateCode(req GenerateRequest) {
	config := reflectToTplConfig(req)
	outBs := tplGenerateCode(config)
	//outBs=strings.Replace(outBs,"\n    \n","\n",-1)
	outB := []byte(outBs)
	outB1, err := kmgFormat.Source(outB)
	if err == nil {
		outB = outB1
	}
	kmgFile.MustWriteFileWithMkdir(req.OutFilePath, outB)
	return
}

type tplConfig struct {
	OutPackageName string          //生成的package的名字 testPackage
	OutKeyByteList string          //生成的key的base64的值
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

type ArgumentNameTypePair struct {
	Name          string
	ObjectTypeStr string
}
