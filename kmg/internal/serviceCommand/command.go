package serviceCommand

import (
	"flag"
	"fmt"
	"github.com/bronze1man/kmg/kmgConsole"
	"github.com/kardianos/service"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

//TODO kmg gorun 可以在service运行的进程里面
//TODO 完整描述使用过程
//TODO 在osx和linux上达到一致的行为
func init() {
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "Service.SetAndRestart",
		Desc:   "install the service,and restart the service,uninstall and stop if need",
		Runner: setAndRestartCmd,
	})
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "Service.Install",
		Desc:   "install the service",
		Runner: installCmd,
	})
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "Service.Uninstall",
		Desc:   "uninstall the serivce",
		Runner: newNameCmd(func(s service.Service) error { return s.Uninstall() }),
	})
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "Service.Start",
		Desc:   "start the service",
		Runner: newNameCmd(func(s service.Service) error { return s.Start() }),
	})
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "Service.Stop",
		Desc:   "start the service",
		Runner: newNameCmd(func(s service.Service) error { return s.Stop() }),
	})
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "Service.Restart",
		Desc:   "restart the service,if the service is not running,start it.",
		Runner: newNameCmd(kmgRestart),
	})
}

type installRequest struct {
	Name             string   //名字
	ExecuteArgs      []string //执行的命令,第一个是执行命令的进程地址
	WorkingDirectory string   //工作目录(默认是当前目录) 不能在upstart系统上面运行
	//有下列选项
	// 'Darwin Launchd'
	// 'systemd'
	// 'Upstart'
	// 'System-V'
	// 'Windows Service'
	SystemName string //使用的系统(linux上面可以进行选择)
}

func setAndRestartCmd() {
	s, err := parseInstallRequest()
	kmgConsole.ExitOnErr(err)
	err = s.Install()
	if err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			kmgConsole.ExitOnErr(err)
		}
		err = s.Uninstall()
		kmgConsole.ExitOnErr(err)
		err = s.Install()
		kmgConsole.ExitOnErr(err)
	}
	err = kmgRestart(s)
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
		Option: service.KeyValue{
			"KeepAlive": false,
		},
	}
	var system service.System
	if runtime.GOOS == "linux" && req.SystemName == "" {
		req.SystemName = "System-V"
	}
	if req.SystemName == "" {
		system = service.ChosenSystem()
		return service.New(nil, svcConfig)
	} else {
		avaliableListS := ""
		for _, thisSystem := range service.AvailableSystems() {
			if !thisSystem.Detect() {
				continue
			}
			avaliableListS += thisSystem.String() + ","
			if thisSystem.String() == req.SystemName {
				system = thisSystem
			}
		}
		if system == nil {
			return nil, fmt.Errorf("system [%s] not exist,avaliable:[%s]", req.SystemName, avaliableListS)
		}
	}
	return system.New(nil, svcConfig)

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

// 允许 stop restart 这种用法
func kmgRestart(s service.Service) (err error) {
	err = s.Stop()
	if err != nil {
		errS := err.Error()
		if !(strings.Contains(errS, "Unknown instance")) {
			return err
		}
	}
	err = s.Start()
	return err
}
