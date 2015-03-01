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

var actionMap = map[string]Command{
	"version": Command{
		Name:   "Version",
		Runner: version,
	},
}

func Main() {
	actionName := ""
	if len(os.Args) >= 2 {
		actionName = os.Args[1]
	}
	lowerActionName := strings.ToLower(actionName)

	action, exist := actionMap[lowerActionName]
	if !exist {
		fmt.Println("command [" + actionName + "] not found.(case insensitive)")
		help()
		return
	}

	os.Args = os.Args[1:]
	action.Runner()
}

func AddAction(action Command) {
	name := strings.ToLower(action.Name)
	_, exist := actionMap[name]
	if exist {
		panic("command " + action.Name + " already defined.(case insensitive)")
	}
	actionMap[name] = action
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
	actionList := make([]Command, 0, len(actionMap))
	for _, command := range actionMap {
		actionList = append(actionList, command)
	}
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
