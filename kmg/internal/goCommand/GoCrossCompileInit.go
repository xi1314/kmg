package goCommand

import (
	"fmt"
	"path/filepath"

	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgConsole"
	"os/exec"
	"runtime"
	"strings"
)

func init() {
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "GoCrossCompileInit",
		Desc:   "cross compile init target in current project",
		Runner: runGoCrossCompileInit,
	})
}

func runGoCrossCompileInit() {
	kmgc, err := kmgConfig.LoadEnvFromWd()
	kmgConsole.ExitOnErr(err)
	GOROOT := kmgc.GOROOT
	if GOROOT == "" {
		//guess GOROOT
		out, err := exec.Command("go", "env", "GOROOT").CombinedOutput()
		kmgConsole.ExitOnErr(err)
		GOROOT = strings.TrimSpace(string(out))
		if GOROOT == "" {
			kmgConsole.ExitOnErr(fmt.Errorf("you must set $GOROOT in environment to use GoCrossComplieInit"))
		}
	}
	var makeShellArgs []string
	var makeShellName string
	runCmdPath := filepath.Join(GOROOT, "src")
	if runtime.GOOS == "windows" {
		makeShellName = "cmd"
		makeShellArgs = []string{"/C", filepath.Join(GOROOT, "src", "make.bat"), "--no-clean"}
	} else {
		makeShellName = filepath.Join(GOROOT, "src", "make.bash")
		makeShellArgs = []string{"--no-clean"}
	}
	for _, target := range kmgc.CrossCompileTarget {
		cmd := kmgCmd.CmdSlice(append([]string{makeShellName}, makeShellArgs...)).GetExecCmd()
		kmgCmd.SetCmdEnv(cmd, "GOOS", target.GetGOOS())
		kmgCmd.SetCmdEnv(cmd, "GOARCH", target.GetGOARCH())
		cmd.Dir = runCmdPath
		err = cmd.Run()
		kmgConsole.ExitOnErr(err)
	}
	return
}
