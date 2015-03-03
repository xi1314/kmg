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
	default:
		//在临时目录建一个package,并且使用GOPATH指向那个临时目录,缓存会出现问题,并且效果和 go build -i 没有区别
		//是文件,文件可以go build -i
		//优化项: 1.是src里面的文件,可以go install? 2.可以创建一个缓存package?
		//build
		outputPath := filepath.Join(goPath, "bin", filepath.Base(pathOrPkg))
		cmd := kmgCmd.NewOsStdioCmd("go", "build", "-i", "-o", outputPath, pathOrPkg)
		err = kmgCmd.SetCmdEnv(cmd, "GOPATH", goPath)
		kmgConsole.ExitOnErr(err)
		err = cmd.Run()
		kmgConsole.ExitOnErr(err)
		//run
		cmd = kmgCmd.NewOsStdioCmd(outputPath, os.Args[2:]...)
		err = kmgCmd.SetCmdEnv(cmd, "GOPATH", goPath)
		kmgConsole.ExitOnErr(err)
		err = cmd.Run()
		kmgConsole.ExitOnErr(err)

		//如果是 两个GOPATH,则只会有一个GOPATH被检查更新,另一个不检查,必须使用一个GOPATH
		//在临时目录建一个package,并且建立一个新的GOPATH指向那个临时目录 ,
		//addGoPath := "/tmp/kmg-gorun-gopath"
		//pkgName := "kmgFakePkg" + kmgRand.MustCryptoRandToHex(10)
		//err = kmgFile.Mkdir(addGoPath + "/src/" + pkgName)
		//kmgConsole.ExitOnErr(err)
		//err = kmgFile.CopyFile(pathOrPkg, addGoPath+"/src/"+pkgName+"/main.go")
		//kmgConsole.ExitOnErr(err)

		//build
		//args := []string{"install", pkgName}
		//cmd := kmgCmd.NewOsStdioCmd("go", args...)
		//err = kmgCmd.SetCmdEnv(cmd, "GOPATH", addGoPath+string(os.PathListSeparator)+goPath)
		//kmgConsole.ExitOnErr(err)
		//err = cmd.Run()
		//kmgConsole.ExitOnErr(err)
		//run
		//cmd = kmgCmd.NewOsStdioCmd(addGoPath+"/bin/"+pkgName, os.Args[2:]...)
		//err = kmgCmd.SetCmdEnv(cmd, "GOPATH", goPath)
		//kmgConsole.ExitOnErr(err)
		//err = cmd.Run()
		//kmgConsole.ExitOnErr(err)
	}
}
