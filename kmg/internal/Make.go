package internal

import (
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgConsole"
	"os"
	"strings"
)

func init() {
	kmgConsole.AddAction(kmgConsole.Command{
		Name: "Make",
		Desc: `run a project defined command
保证在项目根目录下运行
使用普通空格分割方法定义命令`,
		Runner: makeCmd,
	})
}

func makeCmd() {
	kmgc, err := kmgConfig.LoadEnvFromWd()
	kmgConsole.ExitOnErr(err)
	if kmgc.Make == "" {
		kmgConsole.ExitOnStderr("Please defined a Make command in .kmg.yml file to use kmg make")
		return
	}
	args := strings.Split(kmgc.Make, " ")
	err = kmgCmd.CmdSlice(append(args, os.Args[1:]...)).
		SetDir(kmgc.ProjectPath).
		Run()
	kmgConsole.ExitOnErr(err)
}
