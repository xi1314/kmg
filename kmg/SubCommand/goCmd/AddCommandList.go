package goCmd

import "github.com/bronze1man/kmg/kmgConsole"

func AddCommandList() {
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "Go",
		Desc:   "run go command in current project",
		Runner: GoCommand,
	})
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "GoCrossCompile",
		Desc:   "cross compile target in current project",
		Runner: runGoCrossCompile,
	})
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "GoCrossCompileInit",
		Desc:   "cross compile init target in current project",
		Runner: runGoCrossCompileInit,
	})
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "GoRun",
		Desc:   "run go run in current project and use go install to speed up build",
		Runner: GoRunCmd,
	})
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "GoTest",
		Desc:   "递归目录的go test",
		Runner: runGoTest,
	})
	kmgConsole.AddCommandWithName("GoEnv", GoEnvCmd)
}
