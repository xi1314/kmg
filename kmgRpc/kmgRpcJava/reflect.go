package kmgRpcJava

import (
	"fmt"

	"github.com/bronze1man/kmg/kmgGoSource/kmgGoParser"
	"github.com/bronze1man/kmg/kmgStrings"
	"strconv"
)

func reflectToTplConfig(req *GenerateRequest) *tplConfig {
	config := &tplConfig{
		OutPackageName: req.OutPackageName,
		ClassName:      req.OutClassName,
		innerClassMap:  map[string]*InnerClass{},
	}
	if req.ApiNameFilterCb == nil {
		req.ApiNameFilterCb = func(name string) bool {
			return true
		}
	}

	pkg := kmgGoParser.NewProgramFromDefault().GetPackage(req.ObjectPkgPath)
	//pkg := kmgGoParser.MustParsePackage(kmgConfig.DefaultEnv().GOPATHToString(), req.ObjectPkgPath)
	namedTyp := pkg.LookupNamedType(req.ObjectName)
	if namedTyp == nil {
		panic(fmt.Errorf("can not find this object. [%s]", req.ObjectName))
	}
	objTyp := kmgGoParser.Type(namedTyp)
	if req.ObjectIsPointer {
		objTyp = kmgGoParser.NewPointer(objTyp)
	}
	//获取 object 上面所有的方法
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
		RequestInnerClass := &InnerClass{
			Name:     api.Name + "RpcRequest",
			IsPublic: false,
		}
		for _, pairObj := range methodObj.InParameter {
			if isPointerToKmgHttpContext(pairObj.Type) {
				// 此处忽略掉特殊传入参数 *kmgHttp.Context
				continue
			}
			pair := NameTypePair{
				Name: kmgStrings.FirstLetterToUpper(pairObj.Name),
			}
			pair.TypeStr = config.addType(pairObj.Type)
			api.InArgsList = append(api.InArgsList, pair)
			RequestInnerClass.FieldList = append(RequestInnerClass.FieldList, pair)
		}
		config.addInnerClass(RequestInnerClass)
		/*
		 * 在java中,如果输出参数的类型是error,会转到exception里面
		 * 在java中,如果输出参数去除error后,没有参数了,java函数不会返回参数
		 * 在java中,如果输出参数除去error后,还有1个,直接返回这个类型,此时此参数依然按照前面的说法进行向前兼容.
		 * 在java中,如果输出参数除去error后,还有2个及2个以上,返回一个class,这个class包含这些个参数.
		 */
		notErrorOutParameter := []NameTypePair{}
		for i, pairObj := range methodObj.OutParameter {
			builtintyp, ok := pairObj.Type.(kmgGoParser.BuiltinType)
			if !(ok && string(builtintyp) == "error") {
				name := kmgStrings.FirstLetterToUpper(pairObj.Name)
				if name == "" {
					name = "Out_" + strconv.Itoa(i)
				}
				notErrorOutParameter = append(notErrorOutParameter, NameTypePair{
					Name:    name,
					TypeStr: config.addType(pairObj.Type),
				})
			}
		}
		if len(notErrorOutParameter) == 0 {
			api.OutTypeString = "void"
		} else {
			responseInnerClass := &InnerClass{
				Name:      api.Name + "RpcResponse",
				FieldList: notErrorOutParameter,
				IsPublic:  true,
			}
			config.addInnerClass(responseInnerClass)

			if len(notErrorOutParameter) == 1 {
				api.OutTypeString = notErrorOutParameter[0].TypeStr
				api.OutTypeFieldName = notErrorOutParameter[0].Name
			} else {
				api.OutTypeString = responseInnerClass.Name
			}
		}
		config.ApiList = append(config.ApiList, api)
	}
	return config
}

func isPointerToKmgHttpContext(typi kmgGoParser.Type) bool {
	typ1, ok := typi.(kmgGoParser.PointerType)
	if !ok {
		return false
	}
	typ2, ok := typ1.Elem.(*kmgGoParser.NamedType)
	if !ok {
		return false
	}
	return typ2.Name == "Context" && typ2.PackagePath == "github.com/bronze1man/kmg/kmgNet/kmgHttp"
}

func (config *tplConfig) addType(typi kmgGoParser.Type) string {
	switch typ := typi.(type) {
	//case kmgGoParser.FuncType:
	//	panic("TODO")
	case kmgGoParser.PointerType:
		return config.addType(typ.Elem)
	case *kmgGoParser.NamedType:
		if typ.Name == "Time" && typ.PackagePath == "time" {
			return "Date"
		}
		// TODO 命名重复?
		name := typ.Name
		ntyp := config.innerClassMap[name]
		if ntyp != nil {
			return name
		}
		switch utyp := typ.GetUnderType().(type) {
		case kmgGoParser.StructType:
			innerClass := &InnerClass{
				Name:     name,
				IsPublic: true,
			}
			config.addInnerClass(innerClass) // 避免循环引用导致出现bug.
			for i := range utyp.Field {
				innerClass.FieldList = append(innerClass.FieldList, NameTypePair{
					Name:    utyp.Field[i].Name,
					TypeStr: config.addType(utyp.Field[i].Elem),
				})
			}
			return name
		default:
			return config.addType(utyp)
		}
	case kmgGoParser.BuiltinType:
		s, ok := builtinTypeMap[string(typ)]
		if !ok {
			panic("not support BuiltinType " + string(typ))
		}
		return s
	case kmgGoParser.SliceType:
		return "List<" + config.addType(typ.Elem) + ">"
	default:
		panic(fmt.Errorf("TODO %T", typi))
	}
}

var builtinTypeMap = map[string]string{
	"string":  "String",
	"int":     "Integer",
	"int64":   "Long",
	"float64": "Float",
	"bool":    "Boolean",
	"byte":    "Byte",
}
