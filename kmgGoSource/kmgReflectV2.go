package kmgGoSource

import (
	"fmt"
	"go/doc"
	"reflect"
)

type Kind string

const (
	Struct Kind = "Struct"
	Ptr    Kind = "Ptr"
)

//表示一个指针类型
type TypePointer struct {
	Elem Type
}

func (p TypePointer) Kind() Kind {
	return Ptr
}
func (p TypePointer) Method() []Method {
	return p.Elem.Method()
}
func (p TypePointer) WriteInPkgPath(currentPkgPath string) (str string, importPathList []string) {
	str, importPathList = p.Elem.WriteInPkgPath(currentPkgPath)
	return "*" + str, importPathList
}

//表示一个struct类型
type TypeStruct struct {
	docType     *doc.Type
	pkg         *Package
	reflectType reflect.Type
}

func (p TypeStruct) Kind() Kind {
	return Struct
}
func (p TypeStruct) Package() *Package {
	return p.pkg
}
func (p TypeStruct) Method() []Method {
	numMethod := p.reflectType.NumMethod()
	if numMethod == 0 {
		return nil
	}
	out := make([]Method, numMethod)
	for i := 0; i < numMethod; i++ {
		refletMethod := p.reflectType.Method(i)
		method := Method{
			Recv: p,
			Func: Func{
				Name: refletMethod.Name,
			},
		}

		out[i] = method
	}
	return out
}

func (p TypeStruct) WriteInPkgPath(currentPkgPath string) (str string, importPathList []string) {
	thisPkgPath := p.Package().PkgPath()
	if thisPkgPath == "" || currentPkgPath == thisPkgPath {
		return p.docType.Name, nil
	}
	return p.Package().Name() + "." + p.docType.Name, []string{thisPkgPath}
}

//表示一种golang的类型
type Type interface {
	Kind() Kind
	Method() []Method
	//在这个pkgPath里面的表示方法
	WriteInPkgPath(currentPkgPath string) (str string, importPathList []string)
}

type Method struct {
	Recv Type
	Func
}

type Func struct {
	Name       string
	InArgList  []FuncArgument
	OutArgList []FuncArgument
}

type FuncArgument struct {
	Name string
	Type Type
}

//从反射的类型获取一个Type
// 主要: 需要运行时有源代码,需要使用kmg配置GOPATH
func MustNewTypeFromReflectType(typ reflect.Type) Type {
	switch typ.Kind() {
	case reflect.Ptr:
		return TypePointer{
			Elem: MustNewTypeFromReflectType(typ.Elem()),
		}
	case reflect.Struct:
		pkgPath := typ.PkgPath()
		if pkgPath == "" {
			panic(`TODO Handle pkgPath==""`)
		}
		pkg := MustNewPackageFromImportPath(pkgPath)
		docType := pkg.getDocTypeByName(typ.Name())
		return TypeStruct{
			docType:     docType,
			pkg:         pkg,
			reflectType: typ,
		}
	default:
		panic(fmt.Errorf("not process kind:%s", typ.Kind()))
	}

}
