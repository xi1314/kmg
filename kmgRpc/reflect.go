package kmgRpc

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgGoSource"
	"github.com/bronze1man/kmg/kmgStrings"
	"golang.org/x/tools/go/types"
	"path"
)

func reflectToTplConfig(req GenerateRequest) tplConfig {

	config := tplConfig{
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
	OutKeyByteList := fmt.Sprintf("%#v", req.Key[:])
	config.OutKeyByteList = OutKeyByteList[7 : len(OutKeyByteList)-1]

	ObjTyp := kmgGoSource.MustGetGoTypeFromPkgPathAndTypeName(req.ObjectPkgPath, req.ObjectName)
	if req.ObjectIsPointer {
		ObjTyp = types.NewPointer(ObjTyp)
	}
	var importPathList []string
	config.ObjectTypeStr, importPathList = kmgGoSource.MustWriteGoTypes(req.OutPackageImportPath, ObjTyp)
	config.mergeImportPath(importPathList)

	//获取 object的 上面所有的方法
	methodList := kmgGoSource.MustGetMethodListFromGoTypes(ObjTyp)
	for _, methodObj := range methodList {
		if !methodObj.Obj().Exported() {
			continue
		}
		api := Api{
			Name: methodObj.Obj().Name(),
		}
		methodTyp := methodObj.Type().(*types.Signature)
		for i := 0; i < methodTyp.Params().Len(); i++ {
			pairObj := methodTyp.Params().At(i)
			pair := ArgumentNameTypePair{
				Name: kmgStrings.FirstLetterToUpper(pairObj.Name()),
			}
			pair.ObjectTypeStr, importPathList = kmgGoSource.MustWriteGoTypes(req.OutPackageImportPath, pairObj.Type())
			config.mergeImportPath(importPathList)
			api.InArgsList = append(api.InArgsList, pair)
		}
		for i := 0; i < methodTyp.Results().Len(); i++ {
			pairObj := methodTyp.Results().At(i)
			name := kmgStrings.FirstLetterToUpper(pairObj.Name())
			if name == "" {
				if pairObj.Type().String() == "error" { //TODO 不要特例
					name = "Err"
				} else {
					name = fmt.Sprintf("out_%d", i)
				}
			}
			pair := ArgumentNameTypePair{
				Name: name,
			}
			pair.ObjectTypeStr, importPathList = kmgGoSource.MustWriteGoTypes(req.OutPackageImportPath, pairObj.Type())
			config.mergeImportPath(importPathList)
			api.OutArgsList = append(api.OutArgsList, pair)
		}
		config.ApiList = append(config.ApiList, api)
	}
	return config
}
