package kmgConsole

import (
	"fmt"
	"os"
)

var VERSION = "1.0"

type Command struct {
	Name   string
	Runner func()
}

var actionList = []Command{
	{
		Name:   "version",
		Runner: version,
	},
}

func Main() {
	actionName := ""
	if len(os.Args) >= 2 {
		actionName = os.Args[1]
	}
	var action Command
	for i := 0; i < len(actionList); i++ {
		if actionList[i].Name == actionName {
			action = actionList[i]
			break
		}
	}
	if action.Name == "" {
		fmt.Println("command not found")
		help()
		return
	}
	os.Args = os.Args[1:]
	action.Runner()
}

func AddAction(action Command) {
	actionList = append(actionList, action)
}

//avoid initialization loop
func init() {
	AddAction(Command{
		Name:   "help",
		Runner: help,
	})
}

func help() {
	fmt.Println("Usage: ")
	for i := 0; i < len(actionList); i++ {
		fmt.Println("\t", actionList[i].Name)
	}
}

func version() {
	fmt.Println(VERSION)
}
