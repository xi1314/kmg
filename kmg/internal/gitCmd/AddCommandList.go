package gitCmd

import "github.com/bronze1man/kmg/kmgConsole"

func AddCommandList() {
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "GitPull",
		Desc:   "git pull origin master",
		Runner: gitPull,
	})
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "GitPush",
		Desc:   "add,commit,push this git resp",
		Runner: gitPush,
	})
}
