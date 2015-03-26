package serviceCommand

import (
	"flag"
	"fmt"
	"github.com/bronze1man/kmg/kmgConsole"
	"github.com/kardianos/service"
	"os"
	"os/exec"
	"strings"
)

func init() {
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "Service.Install",
		Desc:   "manage system service more easy",
		Runner: installCmd,
	})
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "Service.Uninstall",
		Desc:   "manage system service more easy",
		Runner: uninstallCmd,
	})
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "Service.Start",
		Desc:   "manage system service more easy",
		Runner: startCmd,
	})
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "Service.Stop",
		Desc:   "manage system service more easy",
		Runner: stopCmd,
	})
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "Service.Restart",
		Desc:   "manage system service more easy",
		Runner: restartCmd,
	})
}

type installRequest struct {
	Name             string //名字
	ExecuteString    string //执行命令(暂使用' '模式,如果不够用再说)
	WorkingDirectory string //工作目录(默认是当前目录)
}

func installCmd() {
	req := installRequest{}
	flag.StringVar(&req.Name, "name", "", "name of the service(require)")
	flag.StringVar(&req.ExecuteString, "exec", "", "command to run(require,use ' ' to separate args)")
	flag.StringVar(&req.WorkingDirectory, "cd", "", "working directory(optional),default to currend directory")
	flag.Parse()
	if req.Name == "" {
		fmt.Println("require name args")
		flag.Usage()
		return
	}
	if req.ExecuteString == "" {
		fmt.Println("require ExecuteString args")
		flag.Usage()
		return
	}
	executeArgs := strings.Split(req.ExecuteString, " ")
	executeFilePath, err := exec.LookPath(executeArgs[0])
	kmgConsole.ExitOnErr(err)
	if req.WorkingDirectory == "" {
		req.WorkingDirectory, err = os.Getwd()
		kmgConsole.ExitOnErr(err)
	}
	svcConfig := &service.Config{
		Name:             req.Name,
		Executable:       executeFilePath,
		Arguments:        executeArgs[1:],
		WorkingDirectory: req.WorkingDirectory,
	}
	s, err := service.New(nil, svcConfig)
	kmgConsole.ExitOnErr(err)
	err = s.Install()
	kmgConsole.ExitOnErr(err)
}

func uninstallCmd() {
	req := installRequest{}
	flag.StringVar(&req.Name, "name", "", "name of the service(require)")
	flag.Parse()
	if req.Name == "" {
		fmt.Println("require name args")
		flag.Usage()
		return
	}
	svcConfig := &service.Config{
		Name: req.Name,
	}
	s, err := service.New(nil, svcConfig)
	kmgConsole.ExitOnErr(err)
	err = s.Uninstall()
	kmgConsole.ExitOnErr(err)
}
func startCmd() {
	name := getNameFromArgs()
	svcConfig := &service.Config{
		Name: name,
	}
	s, err := service.New(nil, svcConfig)
	kmgConsole.ExitOnErr(err)
	err = s.Start()
	kmgConsole.ExitOnErr(err)
}
func stopCmd() {
	name := getNameFromArgs()
	svcConfig := &service.Config{
		Name: name,
	}
	s, err := service.New(nil, svcConfig)
	kmgConsole.ExitOnErr(err)
	err = s.Stop()
	kmgConsole.ExitOnErr(err)
}
func restartCmd() {
	name := getNameFromArgs()
	svcConfig := &service.Config{
		Name: name,
	}
	s, err := service.New(nil, svcConfig)
	kmgConsole.ExitOnErr(err)
	err = s.Restart()
	kmgConsole.ExitOnErr(err)
}

func getNameFromArgs() string {
	req := installRequest{}
	flag.StringVar(&req.Name, "name", "", "name of the service(require)")
	flag.Parse()
	if req.Name == "" {
		fmt.Println("require name args")
		flag.Usage()
		os.Exit(1)
		return ""
	}
	return req.Name
}
