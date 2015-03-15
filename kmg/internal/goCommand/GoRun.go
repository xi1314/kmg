package goCommand

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgConsole"
	//"github.com/bronze1man/kmg/kmgFile"
	//"github.com/bronze1man/kmg/kmgRand"
	"os"
	"path/filepath"
)

func init() {
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "GoRun",
		Desc:   "run go run in current project and use go install to speed up build",
		Runner: gorun,
	})
}

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
	case os.IsNotExist(err):
		//这个方案的实现缓存是正常的,但是也只能更新本GOPATH里面的pkg,不能更多多个GOPATH里面其他GOPATH的pkg缓存.
		//是package
		//build
		args := []string{"install", pathOrPkg}
		cmd := kmgCmd.NewOsStdioCmd("go", args...)
		err = kmgCmd.SetCmdEnv(cmd, "GOPATH", goPath)
		kmgConsole.ExitOnErr(err)
		err = cmd.Run()
		kmgConsole.ExitOnErr(err)
		//run
		outPath := filepath.Join(goPath, "bin", filepath.Base(pathOrPkg))
		cmd = kmgCmd.NewOsStdioCmd(outPath, os.Args[2:]...)
		err = kmgCmd.SetCmdEnv(cmd, "GOPATH", goPath)
		kmgConsole.ExitOnErr(err)
		err = cmd.Run()
		kmgConsole.ExitOnErr(err)
		return
	case err != nil:
		kmgConsole.ExitOnErr(err)
		return
	default:
		//靠谱实现这个东西的复杂度太高,目前已有的方案不能达到目标,暂时先使用go run
		// 如果有需要使用请go run一个package
		//已经证实不行的方案:
		// 在临时目录建一个package,并且使用GOPATH指向那个临时目录,缓存会出现问题,并且效果和 go build -i 没有区别
		// 使用go build -i 效果和直接go run没有区别(缓存还是会出现问题)
		//有可能的方案是:
		// 取出指向的这个文件的所有import的包,全部install一遍,再go build,并且run(如何找到某个文件的import项?)
		cmd := kmgCmd.NewOsStdioCmd("go", append([]string{"run"}, os.Args[1:]...)...)
		err = kmgCmd.SetCmdEnv(cmd, "GOPATH", goPath)
		kmgConsole.ExitOnErr(err)
		err = cmd.Run()
		kmgConsole.ExitOnErr(err)
		return
	}
	kmgConsole.ExitOnErr(fmt.Errorf("unexpected run path"))
}
