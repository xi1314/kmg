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

//TODO kmg gorun 可以在service运行的进程里面
//TODO 完整描述使用过程
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
	WorkingDirectory string   //工作目录(默认是当前目录) 不能在upstart系统上面运行
	//有下列选项
	// 'Darwin Launchd'
	// 'Linux systemd'
	// 'Linux Upstart'
	// 'Linux System-V'
	// 'Windows Service'
	SystemName string //使用的系统(linux上面可以进行选择)
}

func setAndStartCmd() {
	s, err := parseInstallRequest()
	kmgConsole.ExitOnErr(err)
	err = s.Install()
	if err != nil {
		err = s.Uninstall()
		kmgConsole.ExitOnErr(err)
		err = s.Install()
		kmgConsole.ExitOnErr(err)
	}
	err = s.Start()
	kmgConsole.ExitOnErr(err)
}

func installCmd() {
	s, err := parseInstallRequest()
	kmgConsole.ExitOnErr(err)
	err = s.Install()
	kmgConsole.ExitOnErr(err)
}

func parseInstallRequest() (s service.Service, err error) {
	req := &installRequest{}
	flag.StringVar(&req.Name, "name", "", "name of the service(require)")
	var executeString string
	flag.StringVar(&executeString, "exec", "", "command to run(require,use ' ' to separate args)")
	flag.StringVar(&req.WorkingDirectory, "cd", "", "working directory(optional),default to currend directory")
	flag.StringVar(&req.SystemName, "system", "", "system name")
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
	svcConfig := &service.Config{
		Name:             req.Name,
		Executable:       req.ExecuteArgs[0],
		Arguments:        req.ExecuteArgs[1:],
		WorkingDirectory: req.WorkingDirectory,
	}
	if req.SystemName == "" {
		return service.New(nil, svcConfig)
	} else {
		avaliableListS := ""
		for _, system := range service.AvailableSystems() {
			avaliableListS += system.String() + ","
			if system.String()==req.SystemName {
				return system.New(nil, svcConfig)
			}
		}
		return nil, fmt.Errorf("system [%s] not exist,avaliable:[%s]", req.SystemName, avaliableListS)
	}

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
