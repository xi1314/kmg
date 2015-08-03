package goCmd

import (
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgTest"
	"os"
	"path/filepath"
	"testing"
)

func TestGoRunFile(ot *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	defer os.Chdir(wd)
	root := kmgConfig.DefaultEnv().ProjectPath
	os.Chdir(root)
	kmgPath := filepath.Join(root, "bin/kmg")
	kmgFile.MustDelete(kmgPath)
	//kmgCmd.MustRun("kmg go install github.com/bronze1man/kmg/kmg")
	gopath := filepath.Join(wd, "testWorkspace")
	os.Chdir(gopath)
	kmgFile.MustDelete(filepath.Join(gopath, "bin"))
	kmgFile.MustDelete(filepath.Join(gopath, "pkg"))
	kmgFile.MustDelete(filepath.Join(gopath, "tmp"))
	kmgFile.MustDelete(filepath.Join(gopath, "src", "kmgTestMain", "other.go"))
	goRunInstall(gopath, "kmgTestMain")
	output := kmgCmd.MustRunAndReturnOutput(filepath.Join(gopath, "bin", "kmgTestMain"))
	kmgTest.Equal(string(output), "1\n")
	kmgFile.MustWriteFile(filepath.Join(wd, "testWorkspace", "src", "kmgTestMain", "other.go"), []byte(`
package main

func init(){
	a=2
}
`))

	goRunInstall(gopath, "kmgTestMain")
	output = kmgCmd.MustRunAndReturnOutput(filepath.Join(gopath, "bin", "kmgTestMain"))
	kmgTest.Equal(string(output), "2\n")
	kmgFile.MustDelete(filepath.Join(wd, "testWorkspace", "src", "kmgTestMain", "other.go"))

	goRunInstall(gopath, "kmgTestMain")
	output = kmgCmd.MustRunAndReturnOutput(filepath.Join(gopath, "bin", "kmgTestMain"))
	kmgTest.Equal(string(output), "1\n")

	goRunInstall(gopath, "kmgTestMain/l2Main")
	output = kmgCmd.MustRunAndReturnOutput(filepath.Join(gopath, "bin", "l2Main"))
	kmgTest.Equal(string(output), "l2Main\n")

	kmgFile.MustWriteFile(filepath.Join(gopath, "src", "replaceBin", "replaceBin.go"), []byte(`
package main

import "fmt"

func main(){
	fmt.Println("replaceBin")
}
`))
	goRunInstall(gopath, "replaceBin")
	output = kmgCmd.MustRunAndReturnOutput(filepath.Join(gopath, "bin", "replaceBin"))
	kmgTest.Equal(string(output), "replaceBin\n")

	kmgFile.MustWriteFile(filepath.Join(gopath, "src", "replaceBin", "replaceBin.go"), []byte(`
package replaceBin

var A = "1"
`))
	goRunInstall(gopath, "replaceBin")
	// 应该会说这个不是main的package
}

/*
func TestGoRunPackageName(ot *testing.T){
	wd:=kmgFile.MustGetWd()
	projectPath:=filepath.Join(wd,"testWorkspace")
	kmgFile.MustDelete(filepath.Join(wd,"testWrokspace","pkg"))
	kmgFile.MustDelete(filepath.Join(wd,"testWorkspace","src","kmgTestMain","other.go"))
	goRunPackageName(projectPath,"kmgTestMain")
}
*/
