package kmgConsole

import (
	"fmt"
	"os"
	"sort"
	"strings"
	//	"runtime/debug"
	//	"github.com/bronze1man/kmg/kmgDebug"
)

var VERSION = ""

type Command struct {
	Name   string //名称,不区分大小写
	Desc   string //暂时没有用到
	Runner func() //运行这个命令的函数
	Hidden bool   //隐藏这个命令,使其在帮助列表里面不显示
}

//一组命令
type CommandGroup struct {
	commandMap map[string]Command
}

func NewCommandGroup() *CommandGroup {
	return &CommandGroup{commandMap: map[string]Command{}}
}

func (g *CommandGroup) Main() {
	actionName := ""
	if len(os.Args) >= 2 {
		actionName = os.Args[1]
	}
	lowerActionName := strings.ToLower(actionName)

	action, exist := g.commandMap[lowerActionName]
	if !exist {
		fmt.Println("command [" + actionName + "] not found.(case insensitive)")
		g.Help()
		return
	}

	os.Args = os.Args[1:]
	//搞这个是为了在panic的时候,仍然可以把输出搞到其他地方去,但是此处复杂度很高,很容易不稳定,具体效果有待确定.
	/*
		defer func(){
			r:=recover()
			if r!=nil{
				fmt.Println("panic:",r)
				os.Stderr.Write(kmgDebug.GetAllStack())
			}
			if len(exitActionList) > 0 {
				for _, action := range exitActionList {
					action()
				}
				exitActionList = nil
			}
		}()
	*/
	action.Runner()
	if len(exitActionList) > 0 {
		WaitForExit()
		for _, action := range exitActionList {
			action()
		}
		exitActionList = nil
	}
}

/*
var actionMap = map[string]Command{
	"version": Command{
		Name:   "Version",
		Runner: version,
	},
}
*/

func (g *CommandGroup) AddCommand(action Command) *CommandGroup {
	name := strings.ToLower(action.Name)
	_, exist := g.commandMap[name]
	if exist {
		panic("command " + action.Name + " already defined.(case insensitive)")
	}
	g.commandMap[name] = action
	return g
}

func (g *CommandGroup) Help() {
	fmt.Println("Usage: ")
	actionList := make([]Command, 0, len(g.commandMap))
	for _, command := range g.commandMap {
		actionList = append(actionList, command)
	}
	sort.Sort(tActionList(actionList))
	for i := 0; i < len(actionList); i++ {
		if actionList[i].Hidden {
			continue
		}
		fmt.Println("\t", actionList[i].Name)
	}
}

func (g *CommandGroup) AddCommandWithName(name string, runner func()) *CommandGroup {
	return g.AddCommand(Command{
		Name:   name,
		Runner: runner,
	})
}

var DefaultCommandGroup = NewCommandGroup()

func Main() {
	if VERSION != "" {
		DefaultCommandGroup.AddCommand(Command{
			Name:   "version",
			Runner: version,
		})
	}
	DefaultCommandGroup.AddCommand(Command{
		Name:   "help",
		Runner: DefaultCommandGroup.Help,
	})
	DefaultCommandGroup.Main()
}

func AddCommand(action Command) *CommandGroup {
	return DefaultCommandGroup.AddCommand(action)
}

func AddCommandWithName(name string, runner func()) *CommandGroup {
	return DefaultCommandGroup.AddCommandWithName(name, runner)
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
