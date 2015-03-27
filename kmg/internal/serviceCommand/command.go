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
		Name:   "Service.SetAndStart",
		Desc:   "manage system service more easy",
		Runner: setAndStartCmd,
	})
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "Service.Install",
		Desc:   "manage system service more easy",
		Runner: installCmd,
	})
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "Service.Uninstall",
		Desc:   "manage system service more easy",
		Runner: newNameCmd(func(s service.Service) error { return s.Uninstall() }),
	})
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "Service.Start",
		Desc:   "manage system service more easy",
		Runner: newNameCmd(func(s service.Service) error { return s.Start() }),
	})
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "Service.Stop",
		Desc:   "manage system service more easy",
		Runner: newNameCmd(func(s service.Service) error { return s.Stop() }),
	})
	//TODO linux restart bug,
	//TODO bug1: you have to first start,then restart.
	//TODO bug2: can not see the real reason.
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "Service.Restart",
		Desc:   "manage system service more easy",
		Runner: newNameCmd(func(s service.Service) error { return s.Restart() }),
	})
}

type installRequest struct {
	Name             string   //名字
	ExecuteArgs      []string //执行的命令,第一个是执行命令的进程地址
	WorkingDirectory string   //工作目录(默认是当前目录)
}

func setAndStartCmd() {
	req, err := parseInstallRequest()
	svcConfig := &service.Config{
		Name:             req.Name,
		Executable:       req.ExecuteArgs[0],
		Arguments:        req.ExecuteArgs[1:],
		WorkingDirectory: req.WorkingDirectory,
	}
	s, err := service.New(nil, svcConfig)
	kmgConsole.ExitOnErr(err)
	err = s.Install()
	if err == nil {
		return
	}
	err = s.Uninstall()
	kmgConsole.ExitOnErr(err)
	err = s.Install()
	kmgConsole.ExitOnErr(err)
	err = s.Start()
	kmgConsole.ExitOnErr(err)
}

func installCmd() {
	req, err := parseInstallRequest()
	svcConfig := &service.Config{
		Name:             req.Name,
		Executable:       req.ExecuteArgs[0],
		Arguments:        req.ExecuteArgs[1:],
		WorkingDirectory: req.WorkingDirectory,
	}
	s, err := service.New(nil, svcConfig)
	kmgConsole.ExitOnErr(err)
	err = s.Install()
	kmgConsole.ExitOnErr(err)
}

func parseInstallRequest() (req *installRequest, err error) {
	req = &installRequest{}
	flag.StringVar(&req.Name, "name", "", "name of the service(require)")
	var executeString string
	flag.StringVar(&executeString, "exec", "", "command to run(require,use ' ' to separate args)")
	flag.StringVar(&req.WorkingDirectory, "cd", "", "working directory(optional),default to currend directory")
	flag.Parse()
	if req.Name == "" {
		return nil, fmt.Errorf("require name args")
	}
	if executeString == "" {
		return nil, fmt.Errorf("require exec args")
	}
	req.ExecuteArgs = strings.Split(executeString, " ")
	req.ExecuteArgs[0], err = exec.LookPath(req.ExecuteArgs[0])
	kmgConsole.ExitOnErr(err)
	if req.WorkingDirectory == "" {
		req.WorkingDirectory, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}
	return req, nil
}

func newNameCmd(fn func(s service.Service) error) func() {
	return func() {
		req := installRequest{}
		flag.StringVar(&req.Name, "name", "", "name of the service(require)")
		flag.Parse()
		name := ""
		switch {
		case req.Name != "":
			name = req.Name
		case flag.Arg(0) != "":
			name = flag.Arg(0)
		default:
			fmt.Println("require name args")
			flag.Usage()
			return
		}
		svcConfig := &service.Config{
			Name: name,
		}
		s, err := service.New(nil, svcConfig)
		kmgConsole.ExitOnErr(err)
		err = fn(s)
		kmgConsole.ExitOnErr(err)
	}
}
