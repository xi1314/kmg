package goCmd

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgConsole"
	//"github.com/bronze1man/kmg/kmgFile"
	//"github.com/bronze1man/kmg/kmgRand"
	"github.com/bronze1man/kmg/kmgReflect"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

func gorun() {
	kmgc, err := kmgConfig.LoadEnvFromWd()
	kmgConsole.ExitOnErr(err)
	goPath := kmgc.GOPATHToString()

	//假设第一个是文件或者package名称,后面是传给命令行的参数
	if len(os.Args) < 2 {
		kmgConsole.ExitOnErr(fmt.Errorf("you need pass in running filename"))
		return
	}
	pathOrPkg := os.Args[1]
	_, err = os.Stat(pathOrPkg)
	switch {
	case os.IsNotExist(err): //package名称 TODO 用户在src下面怎么办?
		//这个方案的实现缓存是正常的,但是也只能更新本GOPATH里面的pkg,不能更多多个GOPATH里面其他GOPATH的pkg缓存.
		//是package
		//build
		runCmdSliceWithGoPath(goPath, []string{"go", "install", pathOrPkg})

		//run
		outPath := filepath.Join(goPath, "bin", filepath.Base(pathOrPkg))
		runCmdSliceWithGoPath(goPath, append([]string{outPath}, os.Args[2:]...))

		return
	case err != nil: //其他错误
		kmgConsole.ExitOnErr(err)
		return

	default: //文件或目录
		//  找出指向的这个文件的所有import的包,全部install一遍,再go run
		//靠谱实现这个东西的复杂度太高,目前已有的方案不能达到目标,暂时先使用go run
		// 如果有需要使用请把这个文件放到package里面,或者运行前删除pkg目录.
		//已经证实不行的方案:
		// 1.在临时目录建一个package,并且使用GOPATH指向那个临时目录,缓存会出现问题,并且效果和 go build -i 没有区别
		// 2.使用go build -i 效果和直接go run没有区别(缓存还是会出现问题)
		//有可能的方案是:
		//

		//找出这个文件所有的 import然后install 一遍
		//先解决文件的问题
		fset := token.NewFileSet()
		pkgs, err := parser.ParseFile(fset, pathOrPkg, nil, parser.ImportsOnly)
		kmgConsole.ExitOnErr(err)
		for _, thisImport := range pkgs.Imports {
			//目前没有找到反序列化golang的双引号的方法,暂时使用简单的办法
			pkgName, err := kmgReflect.UnquoteGolangDoubleQuote(thisImport.Path.Value)
			kmgConsole.ExitOnErr(err)
			runCmdSliceWithGoPath(goPath, []string{"go", "install", pkgName})
		}
		runCmdSliceWithGoPath(goPath, append([]string{"go", "run"}, os.Args[1:]...))
		return
	}
	kmgConsole.ExitOnErr(fmt.Errorf("unexpected run path"))
}

//不回显命令
func runCmdSliceWithGoPath(gopath string, cmdSlice []string) {
	err := kmgCmd.CmdSlice(cmdSlice).
		MustSetEnv("GOPATH", gopath).StdioRun()
	kmgConsole.ExitOnErr(err)
}
