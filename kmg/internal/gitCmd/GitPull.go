package goCommand

import (
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgConsole"
)

func init() {
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "GitPull",
		Desc:   "git pull origin master",
		Runner: gitPull,
	})

}

func gitPull() {
	kmgCmd.NewOsStdioCmdString("git pull origin master").Run()
}
