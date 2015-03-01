package goCommand

import (
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgConsole"
)

func init() {
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "GitPush",
		Desc:   "add,commit,push this git resp",
		Runner: gitPush,
	})

}

func gitPush(){
    kmgCmd.NewOsStdioCmdString("git add -A").Run()
    kmgCmd.NewOsStdioCmdString("git commit -am'save'").Run()
    kmgCmd.NewOsStdioCmdString("git push origin master").Run()
}
