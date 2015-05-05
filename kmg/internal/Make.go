package internal

import (
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgConsole"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgPlatform"
	"github.com/bronze1man/kmg/kmgTime"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func init() {
	kmgConsole.AddAction(kmgConsole.Command{
		Name: "Make",
		Desc: `run a project defined command
保证在项目根目录下运行
使用普通空格分割方法定义命令
将命令输出结果log到文件中`,
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
	os.Chdir(kmgc.ProjectPath)
	logDir := filepath.Join(kmgc.LogPath, "run")
	kmgFile.MustMkdirAll(logDir)
	args := strings.Split(kmgc.Make, " ")
	thisLogFilePath := filepath.Join(logDir, time.Now().Format(kmgTime.FormatFileName)+".log")
	kmgFile.MustWriteFile(thisLogFilePath, []byte{})
	if !kmgPlatform.GetCompiledPlatform().Compatible(kmgPlatform.WindowsAmd64) {
		lastLogPath := filepath.Join(logDir, "last.log")
		kmgFile.MustDeleteFile(lastLogPath)
		kmgCmd.ProxyRun("ln -s " + filepath.Base(thisLogFilePath) + " " + lastLogPath)
	}
	err = kmgCmd.CmdSlice(append(args, os.Args[1:]...)).
		SetDir(kmgc.ProjectPath).
		RunAndTeeOutputToFile(thisLogFilePath)
	kmgConsole.ExitOnErr(err)
}
