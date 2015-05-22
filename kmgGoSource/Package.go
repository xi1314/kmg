package kmgGoSource

import (
	"github.com/bronze1man/kmg/kmgConfig"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"strings"
)

// 这个东西目前 有2个问题,
// 1.从ast生成回Type.
// 2.运行必须能反射.
// 3.import似乎很难处理.
// 4.把各种类型都加入复杂度比较高(生成struct的时候还是要使用的.)

// golang.org/x/tools/go/types 和 golang.org/x/tools/go/ssa 的问题
// 1.接口过度复杂和奇葩, 学习成本高.
// 2.写回类型的时候没有import信息

//表示一个golang里面的package
type Package struct {
	docPkg *doc.Package
}

func (pkg *Package) getDocTypeByName(name string) *doc.Type {
	for _, docType := range pkg.docPkg.Types {
		if docType.Name == name {
			return docType
		}
	}
	return nil
}

// 填在import的那个路径的值
func (pkg *Package) PkgPath() string {
	return pkg.docPkg.ImportPath
}

// package 的名字,填在 package 那个地方的名字
// package kmgReflect
func (pkg *Package) Name() string {
	return pkg.docPkg.Name
}

//使用importPath获取一个Package,
// 仅支持导入一个目录作为一个package
// 只导入一个package里面的主package,xxx_test不导入.
// 注意: 需要源代码,需要使用kmg配置GOPATH
func MustNewPackageFromImportPath(importPath string) *Package {
	astPkg, _ := MustNewMainAstPackageFromImportPath(importPath)
	docPkg := doc.New(astPkg, importPath, doc.AllMethods)
	return &Package{
		docPkg: docPkg,
	}
}

func GetImportPathListFromFile(filepath string) (importPathList []string, err error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFile(fset, filepath, nil, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}
	for _, thisImport := range pkgs.Imports {
		//目前没有找到反序列化golang的双引号的方法,暂时使用简单的办法
		pkgName, err := UnquoteGolangDoubleQuote(thisImport.Path.Value)
		if err != nil {
			return nil, err
		}
		importPathList = append(importPathList, pkgName)
	}
	return importPathList, nil
}

func MustNewMainAstPackageFromImportPath(importPath string) (pkg *ast.Package, fset *token.FileSet) {
	pkgDir := kmgConfig.DefaultEnv().MustGetPathFromImportPath(importPath)
	fset = token.NewFileSet()
	astPkgMap, err := parser.ParseDir(fset, pkgDir, nil, 0)
	if err != nil {
		panic(err)
	}
	if len(astPkgMap) == 0 {
		panic("can not found package")
	}
	for name, astPkg := range astPkgMap {
		if strings.HasSuffix(name, "_test") {
			continue
		}
		return astPkg, fset
	}
	panic("impossible execute path")
}
