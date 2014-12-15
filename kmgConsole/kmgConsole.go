package kmgConsole

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

var VERSION = "1.0"

type Command struct {
	Name   string
	Desc   string
	Runner func()
}

var actionList = []Command{
	{
		Name:   "Version",
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
		if strings.EqualFold(actionList[i].Name, actionName) {
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
		Name:   "Help",
		Runner: help,
	})
}

func help() {
	fmt.Println("Usage: ")
	sort.Sort(tActionList(actionList))
	for i := 0; i < len(actionList); i++ {
		fmt.Println("\t", actionList[i].Name)
	}
}

func version() {
	fmt.Println(VERSION)
}

type tActionList []Command

func (t tActionList) Len() int      { return len(t) }
func (t tActionList) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t tActionList) Less(i, j int) bool {
	return strings.ToLower(t[i].Name) < strings.ToLower(t[j].Name)
}
