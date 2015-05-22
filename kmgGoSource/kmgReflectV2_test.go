package kmgGoSource_test

import (
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgDebug"
	"github.com/bronze1man/kmg/kmgGoSource"
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	"github.com/bronze1man/kmg/kmgTest"
	"go/ast"
	"go/parser"
	"go/token"
	_ "golang.org/x/tools/go/gcimporter"
	"golang.org/x/tools/go/types"
	"reflect"
	"testing"
)

func TestMustNewPackageFromImportPath(ot *testing.T) {
	pkg := kmgGoSource.MustNewPackageFromImportPath("github.com/bronze1man/kmg/kmgReflect")
	kmgTest.Equal(pkg.PkgPath(), "github.com/bronze1man/kmg/kmgReflect")

	pkg = kmgGoSource.MustNewPackageFromImportPath("go/doc")
	kmgTest.Equal(pkg.PkgPath(), "go/doc")
}

func TestMustNewTypeFromReflectType(ot *testing.T) {
	typ := kmgGoSource.MustNewTypeFromReflectType(reflect.TypeOf(kmgGoSource.Package{}))
	kmgTest.Equal(typ.Kind(), kmgGoSource.Struct)
}

type TestB struct {
}

func FuncA(a *kmgHttp.Context, c int) (b error) {
	return nil
}

func TestDemo(ot *testing.T) {
	importPath := "github.com/bronze1man/kmg/kmgGoSource"
	kmgCmd.MustRun("kmg go test -i " + importPath)
	pkgDir := kmgConfig.DefaultEnv().MustGetPathFromImportPath(importPath)
	fset := token.NewFileSet()
	astPkgMap, err := parser.ParseDir(fset, pkgDir, nil, 0)
	if err != nil {
		panic(err)
	}
	astPkg := astPkgMap["kmgGoSource_test"]
	astFileList := []*ast.File{}
	for _, file := range astPkg.Files {
		astFileList = append(astFileList, file)
	}
	//os.Chdir(kmgConfig.DefaultEnv().ProjectPath)
	pkg, err := types.Check(pkgDir, fset, astFileList)
	if err != nil {
		panic(err)
	}
	funcA := pkg.Scope().Lookup("FuncA")
	recvPkg := types.NewPackage("github.com/bronze1man/kmg/kmgGoSource", "kmgGoSource")
	kmgDebug.Println(types.TypeString(recvPkg, funcA.Type()))
	funTypParams := funcA.Type().(*types.Signature).Params()
	for i := 0; i < funTypParams.Len(); i++ {
		kmgDebug.Println(funTypParams.At(i).Name())
		kmgDebug.Println(funTypParams.At(i).Type().String())
	}
	//for _,p:=range funcA.Type().(*types.Signature).Params().
	//kmgDebug.Println(funcA.Type().(*types.Signature).Params().String())
}
