package kmgGoParser

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/golang/groupcache/singleflight"
	"path/filepath"
	"strings"
	"sync"
)

// 表示整个程序
type Program struct {
	PackageLookupPathList []string            // GOPATH, GOROOT
	CachedPackageMap      map[string]*Package // pkgPath
	groupCache            singleflight.Group
	locker                sync.Mutex
}

func NewProgramFromDefault() *Program {
	return &Program{
		PackageLookupPathList: append(kmgConfig.DefaultEnv().GOPATH, kmgConfig.DefaultEnv().GetGOROOT()),
		CachedPackageMap:      map[string]*Package{},
	}
}

func NewProgram(lookupPathList []string) *Program {
	return &Program{
		PackageLookupPathList: lookupPathList,
		CachedPackageMap:      map[string]*Package{},
	}
}

func (prog *Program) GetPackage(pkgPath string) *Package {
	pkg := prog.getCachedPackage(pkgPath)
	if pkg != nil {
		return pkg
	}
	_, err := prog.groupCache.Do(pkgPath, func() (interface{}, error) {
		pkg = prog.mustParsePackage(pkgPath)
		prog.setCachedPackage(pkg)
		return nil, nil
	})
	if err != nil {
		panic(err)
	}
	return prog.getCachedPackage(pkgPath)
}

func (prog *Program) GetNamedType(pkgPath string, name string) *NamedType {
	return prog.GetPackage(pkgPath).LookupNamedType(name)
}

func (prog *Program) mustParsePackage(pkgPath string) *Package {
	var dirPath string
	var found bool
	for _, lookupPath := range prog.PackageLookupPathList {
		dirPath = filepath.Join(lookupPath, "src", pkgPath)
		if kmgFile.MustFileExist(dirPath) {
			found = true
			break
		}
	}
	if !found {
		panic(fmt.Errorf("can not found pkgPath %s %#v", pkgPath, prog.PackageLookupPathList))
	}

	pkg := &Package{
		ImportMap: map[string]bool{},
		PkgPath:   pkgPath,
	}
	for _, path := range kmgFile.MustReadDirFileOneLevel(dirPath) {
		if strings.HasSuffix(path, ".go") {
			pkg.mustAddFile(filepath.Join(dirPath, path))
		}
	}
	pkg.Program = prog
	return pkg
}

func (prog *Program) getCachedPackage(pkgPath string) *Package {
	prog.locker.Lock()
	defer prog.locker.Unlock()
	return prog.CachedPackageMap[pkgPath]
}

func (prog *Program) setCachedPackage(pkg *Package) {
	prog.locker.Lock()
	defer prog.locker.Unlock()
	prog.CachedPackageMap[pkg.PkgPath] = pkg
}
