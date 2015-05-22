package kmgRpc

import (
	"github.com/bronze1man/kmg/encoding/kmgBase64"
	"github.com/bronze1man/kmg/kmgGoSource"
	"golang.org/x/tools/go/types"
	"path"
	"reflect"
)

func reflectToTplConfig(req GenerateRequest) tplConfig {
	config := tplConfig{
		ObjectName:     req.ObjectName,
		OutPackageName: path.Base(req.OutPackageImportPath),
		OutKeyBase64:   kmgBase64.Base64EncodeByteToString(req.Key[:]),
		ImportPathMap: map[string]bool{
			"encoding/json": true,
			"errors":        true,
			"fmt":           true,
			"github.com/bronze1man/kmg/encoding/kmgBase64": true,
			"github.com/bronze1man/kmg/kmgCrypto":          true,
			"github.com/bronze1man/kmg/kmgLog":             true,
			"github.com/bronze1man/kmg/kmgNet/kmgHttp":     true,
			"net/http": true,
			"bytes":    true,
		},
	}
	ObjTyp := kmgGoSource.MustGetGoTypesFromReflect(reflect.TypeOf(req.Object))
	var importPathList []string
	config.ObjectTypeStr, importPathList = kmgGoSource.MustWriteGoTypes(req.OutPackageImportPath, ObjTyp)
	config.mergeImportPath(importPathList)

	//获取 object的 上面所有的方法
	methodList := kmgGoSource.MustGetMethodListFromGoTypes(ObjTyp)
	for _, methodObj := range methodList {
		api := Api{
			Name: methodObj.Obj().Name(),
		}
		methodTyp := methodObj.Type().(*types.Signature)
		for i := 0; i < methodTyp.Params().Len(); i++ {
			pairObj := methodTyp.Params().At(i)
			pair := ArgumentNameTypePair{
				Name: pairObj.Name(),
			}
			pair.ObjectTypeStr, importPathList = kmgGoSource.MustWriteGoTypes(req.OutPackageImportPath, pairObj.Type())
			config.mergeImportPath(importPathList)
			api.InArgsList = append(api.InArgsList, pair)
		}
		for i := 0; i < methodTyp.Results().Len(); i++ {
			pairObj := methodTyp.Results().At(i)
			pair := ArgumentNameTypePair{
				Name: pairObj.Name(),
			}
			pair.ObjectTypeStr, importPathList = kmgGoSource.MustWriteGoTypes(req.OutPackageImportPath, pairObj.Type())
			config.mergeImportPath(importPathList)
			api.OutArgsList = append(api.OutArgsList, pair)
		}
		config.ApiList = append(config.ApiList, api)
	}
	return config
}
