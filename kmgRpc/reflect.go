package kmgRpc

import (
	"fmt"
	"path"

	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgGoSource/kmgGoParser"
	"github.com/bronze1man/kmg/kmgStrings"
)

func reflectToTplConfig(req *GenerateRequest) *tplConfig {
	config := &tplConfig{
		ObjectName:     req.ObjectName,
		OutPackageName: path.Base(req.OutPackageImportPath),
		ImportPathMap: map[string]bool{
			"encoding/json": true,
			"errors":        true,
			"fmt":           true,
			"github.com/bronze1man/kmg/kmgCrypto":      true,
			"github.com/bronze1man/kmg/kmgLog":         true,
			"github.com/bronze1man/kmg/kmgNet/kmgHttp": true,
			"net/http": true,
			"bytes":    true,
		},
	}
	if req.ApiNameFilterCb == nil {
		req.ApiNameFilterCb = func(name string) bool {
			return true
		}
	}

	pkg := kmgGoParser.MustParsePackage(kmgConfig.DefaultEnv().GOPATHToString(), req.ObjectPkgPath)
	namedTyp := pkg.LookupNamedType(req.ObjectName)
	if namedTyp == nil {
		panic(fmt.Errorf("can not find this object. [%s]", req.ObjectName))
	}
	objTyp := kmgGoParser.Type(namedTyp)
	if req.ObjectIsPointer {
		objTyp = kmgGoParser.NewPointer(objTyp)
	}
	var importPathList []string
	config.ObjectTypeStr, importPathList = kmgGoParser.MustWriteGoTypes(req.OutPackageImportPath, objTyp)
	config.mergeImportPath(importPathList)

	//获取 object的 上面所有的方法
	methodList := pkg.GetNamedTypeMethodSet(namedTyp)
	for _, methodObj := range methodList {
		if !methodObj.IsExport() {
			continue
		}
		if !req.ApiNameFilterCb(methodObj.Name) {
			continue
		}
		api := Api{
			Name: methodObj.Name,
		}
		for _, pairObj := range methodObj.InParameter {
			pair := ArgumentNameTypePair{
				Name: kmgStrings.FirstLetterToUpper(pairObj.Name),
			}
			pair.ObjectTypeStr, importPathList = kmgGoParser.MustWriteGoTypes(req.OutPackageImportPath, pairObj.Type)
			config.mergeImportPath(importPathList)
			api.InArgsList = append(api.InArgsList, pair)
		}
		for i, pairObj := range methodObj.OutParameter {
			name := kmgStrings.FirstLetterToUpper(pairObj.Name)
			if name == "" {
				builtintyp, ok := pairObj.Type.(kmgGoParser.BuiltinType)
				if ok && string(builtintyp) == "error" { //TODO 不要特例
					name = "Err"
				} else {
					name = fmt.Sprintf("Out_%d", i)
				}
			}
			pair := ArgumentNameTypePair{
				Name: name,
			}
			pair.ObjectTypeStr, importPathList = kmgGoParser.MustWriteGoTypes(req.OutPackageImportPath, pairObj.Type)
			config.mergeImportPath(importPathList)
			api.OutArgsList = append(api.OutArgsList, pair)
		}
		config.ApiList = append(config.ApiList, api)
	}
	return config
}
