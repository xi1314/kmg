package goCmd

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgConsole"
	"github.com/bronze1man/kmg/kmgGoSource"
	"os"
	"path/filepath"
)

// go install bug
func GoRunCmd() {
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
	case os.IsNotExist(err): //package名称
		goRunPackageName(goPath, pathOrPkg)
		return
	case err != nil: //其他错误
		kmgConsole.ExitOnErr(err)
		return

	default: //文件或目录
		wd, err := os.Getwd()
		kmgConsole.ExitOnErr(err)
		if wd == filepath.Join(goPath, "src") {
			//用户在src下
			goRunPackageName(goPath, pathOrPkg)
			return
		}

		//  找出指向的这个文件的所有import的包,全部install一遍,再go run
		//靠谱实现这个东西的复杂度太高,目前已有的方案不能达到目标,暂时先使用go run
		// 如果有需要使用请把这个文件放到package里面,或者运行前删除pkg目录.
		// TODO 速度比较慢.
		//已经证实不行的方案:
		// 1.在临时目录建一个package,并且使用GOPATH指向那个临时目录,缓存会出现问题,并且效果和 go build -i 没有区别
		// 2.使用go build -i 效果和直接go run没有区别(缓存还是会出现问题)

		//找出这个文件所有的 import然后install 一遍
		importPathList, err := kmgGoSource.GetImportPathListFromFile(pathOrPkg)
		kmgConsole.ExitOnErr(err)
		for _, pkgPath := range importPathList {
			runCmdSliceWithGoPath(goPath, []string{"go", "install", pkgPath})
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

func goRunPackageName(goPath string, pathOrPkg string) {
	//这个方案的实现缓存是正常的,但是也只能更新本GOPATH里面的pkg,不能更多多个GOPATH里面其他GOPATH的pkg缓存.
	//是package
	//build
	// TODO 已知bug1 删除某个package里面的部分文件,然后由于引用到了旧的实现的代码,不会报错.删除pkg解决问题.
	// TODO 已知bug2 如果一个package先是main,然后build了一个东西,然后又改成了非main,再gorun会使用旧的缓存/bin/里面的缓存.
	runCmdSliceWithGoPath(goPath, []string{"go", "install", pathOrPkg})

	//run
	outPath := filepath.Join(goPath, "bin", filepath.Base(pathOrPkg))
	runCmdSliceWithGoPath(goPath, append([]string{outPath}, os.Args[2:]...))
}
