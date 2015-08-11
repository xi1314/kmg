package kmgGoParser

import (
	"fmt"
	"path"
)

// 表示一个package,注意: 一个目录下最多会有2个package.暂时忽略xxx_test 这种package
type Package struct {
	PkgPath       string
	ImportMap     map[string]bool //这个package第一层import
	FuncList      []*FuncOrMethodDeclaration
	MethodList    []*FuncOrMethodDeclaration
	NamedTypeList []*NamedType
}

func (pkg *Package) GetImportList() []string {
	output := make([]string, 0, len(pkg.ImportMap))
	for imp := range pkg.ImportMap {
		output = append(output, imp)
	}
	return output
}

func (pkg *Package) AddImport(pkgPath string) {
	if pkgPath != "" {
		pkg.ImportMap[pkgPath] = true
	}
}

func (pkg *Package) GetNamedTypeMethodSet(typ *NamedType) (output []*FuncOrMethodDeclaration) {
	if typ.PackagePath != pkg.PkgPath {
		panic(fmt.Errorf("can not get MethodSet on diff pacakge typ[%s] pkg[%s]", typ.PackagePath, pkg.PkgPath))
	}
	for _, decl := range pkg.MethodList {
		recvier := decl.ReceiverType
		if recvier.GetKind() == Ptr {
			recvier = recvier.(PointerType).Elem
		}
		if recvier.GetKind() != Named {
			panic(fmt.Errorf("[GetNamedTypeMethodSet] reciver is not a named type %T %s", recvier, recvier.GetKind()))
		}
		if recvier.(*NamedType).Name == typ.Name {
			output = append(output, decl)
		}
	}
	return output
}

// 表示一个go的file
type File struct {
	PackageName    string //最开头的package上面写的东西.
	PackagePath    string
	ImportMap      map[string]bool //这个文件的导入表
	AliasImportMap map[string]string
	FuncList       []*FuncOrMethodDeclaration
	MethodList     []*FuncOrMethodDeclaration
	NamedTypeList  []*NamedType
}

func (pkg *File) AddImport(pkgPath string, aliasPath string) {
	if pkgPath == "" {
		return
	}
	pkg.ImportMap[pkgPath] = true
	if aliasPath == "" {
		aliasPath = path.Base(pkgPath)
	}
	pkg.AliasImportMap[aliasPath] = pkgPath
}

func (gofile *File) LookupFullPackagePath(pkgAliasPath string) (string, error) {
	pkgPath := gofile.AliasImportMap[pkgAliasPath]
	if pkgPath == "" {
		return pkgAliasPath, fmt.Errorf("unable to find import alias %s", pkgAliasPath)
	}
	return pkgPath, nil
}
