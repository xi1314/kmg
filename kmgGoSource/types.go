package kmgGoSource

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgCmd"
	"go/ast"
	_ "golang.org/x/tools/go/gcimporter"
	"golang.org/x/tools/go/types"
	"path"
	"reflect"
)

func MustGetGoTypesFromReflect(typ reflect.Type) types.Type {
	switch typ.Kind() {
	case reflect.Ptr:
		return types.NewPointer(MustGetGoTypesFromReflect(typ.Elem()))
	case reflect.Struct:
		if typ.PkgPath() == "" {
			panic(fmt.Errorf(`[MustGetGoTypesFromReflect] Not implement typ.PkgPath=="" name[%s]`,
				typ.Name()))
		}
		//此处没有办法获取Package的实际Package名称
		pkg := MustNewGoTypesMainPackageFromImportPath(typ.PkgPath())
		typObj := pkg.Scope().Lookup(typ.Name())
		return typObj.Type()
	default:
		panic(fmt.Errorf("[MustGetGoTypesFromReflect] Not implement Kind [%s]",
			typ.Kind().String()))
	}
}

func MustWriteGoTypes(thisPackagePath string, typi types.Type) (s string, addPkgPathList []string) {
	switch typ := typi.(type) {
	case *types.Basic:
		return typ.String(), nil
	case *types.Named:
		if typ.Obj().Pkg() == nil {
			return typ.Obj().Name(), nil
		}
		typPkgPath := typ.Obj().Pkg().Path()
		if thisPackagePath == typPkgPath {
			return typ.Obj().Name(), nil
		}
		return path.Base(typPkgPath) + "." + typ.Obj().Name(), []string{typPkgPath}
	case *types.Pointer:
		s, addPkgPathList = MustWriteGoTypes(thisPackagePath, typ.Elem())
		return "*" + s, addPkgPathList
	case *types.Slice:
		s, addPkgPathList = MustWriteGoTypes(thisPackagePath, typ.Elem())
		return "[]" + s, addPkgPathList
	case *types.Interface:
		return typ.String(), nil
		//s, addPkgPathList = MustWriteGoTypes(thisPackagePath, typ.Elem())
		//return "[]" + s, addPkgPathList
	default:
		panic(fmt.Errorf("[MustWriteGoTypes] Not implement go/types [%T] [%s]",
			typi, typi.String()))
	}
	return "", nil
}

func MustGetMethodListFromGoTypes(typ types.Type) (output []*types.Selection) {
	methodSet := types.NewMethodSet(typ)
	num := methodSet.Len()
	if num == 0 {
		return nil
	}
	output = make([]*types.Selection, num)
	for i := range output {
		output[i] = methodSet.At(i)
	}
	return output
}

//返回这个导入路径的主Package的types.Package对象
//TODO 解决测试package的问题
func MustNewGoTypesMainPackageFromImportPath(importPath string) *types.Package {
	kmgCmd.CmdSlice([]string{"kmg", "go", "install", importPath}).MustRun()
	kmgCmd.CmdSlice([]string{"kmg", "go", "test", "-i", importPath}).MustRun()
	//TODO 解决需要预先创建pkg的问题.
	astPkg, fset := MustNewMainAstPackageFromImportPath(importPath)
	astFileList := []*ast.File{}
	for _, file := range astPkg.Files {
		astFileList = append(astFileList, file)
	}
	//os.Chdir(kmgConfig.DefaultEnv().ProjectPath)
	conf := &types.Config{
		IgnoreFuncBodies: true,
	}
	pkg, err := conf.Check(importPath, fset, astFileList, nil)
	if err != nil {
		panic(err)
	}
	return pkg
}
