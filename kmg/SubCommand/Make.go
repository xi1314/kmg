package SubCommand

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"errors"
	"fmt"
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgConsole"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgPlatform"
	"github.com/bronze1man/kmg/kmgTime"
)

func makeCmd() {
	kmgc, err := kmgConfig.LoadEnvFromWd()
	kmgConsole.ExitOnErr(err)
	if kmgc.Make == "" {
		kmgConsole.ExitOnStdErr(errors.New("Please defined a Make command in .kmg.yml file to use kmg make"))
		return
	}
	if len(os.Args) >= 2 && kmgc.MakeSubCommandMap != nil {
		for cmdName, cmdString := range kmgc.MakeSubCommandMap {
			if strings.EqualFold(cmdName, os.Args[1]) {
				args := strings.Split(cmdString, " ")
				os.Args = os.Args[1:]
				runCommand(kmgc, args)
				return
			}
		}
	}
	args := strings.Split(kmgc.Make, " ")
	runCommand(kmgc, args)
}

func runCommand(kmgc *kmgConfig.Env, args []string) {
	os.Chdir(kmgc.ProjectPath)
	logDir := filepath.Join(kmgc.LogPath, "run")
	kmgFile.MustMkdirAll(logDir)
	thisLogFilePath := filepath.Join(logDir, time.Now().Format(kmgTime.FormatFileName)+".log")
	kmgFile.MustWriteFile(thisLogFilePath, []byte{})
	if !kmgPlatform.GetCompiledPlatform().Compatible(kmgPlatform.WindowsAmd64) {
		lastLogPath := filepath.Join(logDir, "last.log")
		if kmgFile.MustFileExist(lastLogPath) {
			kmgFile.MustSymlink(kmgFile.MustReadSymbolLink(lastLogPath), filepath.Join(logDir, "last2.log"))
		}
		kmgFile.MustSymlink(filepath.Base(thisLogFilePath), lastLogPath)
	}
	//TODO 大部分命令是 kmg gorun xxx 在这个地方可以直接调用gorun解决问题,这样可以少开一个进程加快了一些速度
	// 问题: 上述做法不靠谱,会导致last.log没有用处.
	//if len(args) >= 2 && args[0] == "kmg" && strings.EqualFold(args[1], "gorun") {
	//	os.Args = append(args[1:], os.Args[1:]...)
	//	goCmd.GoRunCmd()
	//	return
	//}
	// 下面的做法不靠谱,会让signle无法传递
	//err := kmgCmd.CmdSlice(append(args, os.Args[1:]...)).
	//	SetDir(kmgc.ProjectPath).
	//	RunAndTeeOutputToFile(thisLogFilePath)
	// TODO bash转义
	bashCmd := strings.Join(append(args, os.Args[1:]...), " ")
	bashCmdStr:=bashCmd + " 2>&1 | tee -i " + thisLogFilePath + " ; test ${PIPESTATUS[0]} -eq 0"
	if kmgPlatform.IsWindows(){
		err := kmgCmd.CmdString(bashCmdStr).SetDir(kmgc.ProjectPath).StdioRun()
		if err != nil {
			err = fmt.Errorf("kmg make: %s", err)
			kmgConsole.ExitOnErr(err)
		}
		return
	}else {
		err := kmgCmd.CmdBash(bashCmdStr).SetDir(kmgc.ProjectPath).StdioRun()
		if err != nil {
			err = fmt.Errorf("kmg make: %s", err)
			kmgConsole.ExitOnErr(err)
		}
		return
	}
}
