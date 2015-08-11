package kmgGoParser

import (
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
	//"github.com/bronze1man/kmg/kmgDebug"
	"github.com/bronze1man/kmg/kmgFile"
	"path/filepath"
	"strings"
)

func TestMustParsePackage(ot *testing.T) {
	pkg := MustParsePackage(kmgConfig.DefaultEnv().GOPATHToString(), "github.com/bronze1man/kmg/kmgGoSource/kmgGoParser/testPackage")
	kmgTest.Equal(pkg.GetImportList(), []string{"errors", "bytes"})
}

func TestMustParsePackageFunc(ot *testing.T) {
	pkg := MustParsePackage(kmgConfig.DefaultEnv().GOPATHToString(), "github.com/bronze1man/kmg/kmgGoSource/kmgGoParser/testPackage/testFunc")
	kmgTest.Equal(len(pkg.FuncList), 7)
}

func TestParseGoSrc(ot *testing.T) {
	gopath := "/usr/local/go"
	goSourcePath := filepath.Join(gopath, "src")
	dirList := kmgFile.MustGetAllDir(goSourcePath)
	for _, dir := range dirList {
		if strings.Contains(dir, "cmd/go/testdata") {
			continue
		}
		dir, err := filepath.Rel(goSourcePath, dir)
		if err != nil {
			panic(err)
		}
		MustParsePackage(gopath, dir)
	}
}

func TestParseCurrentProject(ot *testing.T) {
	gopath := kmgConfig.DefaultEnv().GOPATHToString()
	goSourcePath := filepath.Join(gopath, "src")
	dirList := kmgFile.MustGetAllDir(goSourcePath)
	for _, dir := range dirList {
		if strings.Contains(dir, "go/loader/testdata") {
			continue
		}
		dir, err := filepath.Rel(goSourcePath, dir)
		if err != nil {
			panic(err)
		}
		MustParsePackage(gopath, dir)
	}
}
