package serviceCmd

import (
	"os"
	"os/exec"
	"strings"

	"github.com/bronze1man/kmg/errors"
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgTextTemplate"
)

var ErrServiceExist = errors.New("Service already exist")

//这个是服务的定义数据结构,会序列化到具体实现的配置文件里面去
type Service struct {
	Name             string //Required
	WorkingDirectory string
	CommandLineSlice []string //Required 运行命令需要的slice数组.不要包含空格
}

func (s *Service) init() (err error) {
	if s.Name == "" {
		return errors.New("[Service] require Service.Name")
	}
	if s.WorkingDirectory == "" {
		s.WorkingDirectory = "/"
	}
	if len(s.CommandLineSlice) == 0 {
		return errors.New("[Service] require Service.CommandLineSlice")
	}
	s.CommandLineSlice[0], err = exec.LookPath(s.CommandLineSlice[0])
	return err
}

func (s *Service) GetCommandLineBashString() string {
	return kmgCmd.BashEscape(strings.Join(s.CommandLineSlice, " "))
}

func getConfigPath(name string) string {
	return "/etc/init.d/" + name
}

func Install(s *Service) (err error) {
	err = s.init()
	if err != nil {
		return err
	}
	confPath := getConfigPath(s.Name)
	_, err = os.Stat(confPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil {
		return ErrServiceExist
	}
	out := kmgTextTemplate.MustRender(sysvScript, s)
	err = kmgFile.WriteFile(confPath, out)
	if err != nil {
		return err
	}
	err = os.Chmod(confPath, 0755)
	if err != nil {
		return err
	}
	for _, i := range [...]string{"2", "3", "4", "5"} {
		if err = os.Symlink(confPath, "/etc/rc"+i+".d/S50"+s.Name); err != nil {
			continue
		}
	}
	for _, i := range [...]string{"0", "1", "6"} {
		if err = os.Symlink(confPath, "/etc/rc"+i+".d/K02"+s.Name); err != nil {
			continue
		}
	}
	return nil
}

func Uninstall(name string) (err error) {
	Stop(name)
	// TODO Stop 里面关闭不是错误.
	confPath := getConfigPath(name)
	pathList := []string{
		confPath,
	}
	for _, i := range [...]string{"2", "3", "4", "5"} {
		pathList = append(pathList, "/etc/rc"+i+".d/S50"+name)
	}
	for _, i := range [...]string{"0", "1", "6"} {
		pathList = append(pathList, "/etc/rc"+i+".d/K02"+name)
	}
	for _, path := range pathList {
		kmgFile.MustDelete(path)
	}
	return nil
}

func Start(name string) (err error) {
	c := waitRpcRespond()
	err = kmgCmd.CmdSlice([]string{"service", name, "start"}).Run()
	if err != nil {
		return err
	}
	return <-c
}

//TODO 已经关闭不是错误
func Stop(name string) (err error) {
	return kmgCmd.CmdSlice([]string{"service", name, "stop"}).Run()
}

func Restart(name string) (err error) {
	c := waitRpcRespond()
	err = kmgCmd.CmdSlice([]string{"service", name, "restart"}).Run()
	if err != nil {
		return err
	}
	return <-c
}

const sysvScript = `#!/bin/bash
# For RedHat and cousins:
# chkconfig: - 99 01

### BEGIN INIT INFO
# Required-Start:
# Required-Stop:
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
### END INIT INFO

ulimit -n 1048576

cmd={{.GetCommandLineBashString}}

name='{{.Name}}'
pid_file="/var/run/$name.pid"
stdout_log="/var/log/$name.log"

get_pid() {
    cat "$pid_file"
}

is_running() {
    [ -f "$pid_file" ] && ps $(get_pid) > /dev/null 2>&1
}

case "$1" in
    start)
        if is_running; then
            echo "Already started"
        else
            echo "Starting $name"
            cd '{{.WorkingDirectory}}'
            $cmd &>> "$stdout_log" &
            echo $! > "$pid_file"
            if ! is_running; then
                echo "Unable to start, see $stdout_log and $stderr_log"
                exit 1
            fi
        fi
    ;;
    stop)
        if is_running; then
            echo -n "Stopping $name.."
            kill $(get_pid)
            for i in {1..10}
            do
                if ! is_running; then
                    break
                fi
                echo -n "."
                sleep 1
            done
            echo
            if is_running; then
                echo "Not stopped; may still be shutting down or shutdown may have failed"
                exit 1
            else
                echo "Stopped"
                if [ -f "$pid_file" ]; then
                    rm "$pid_file"
                fi
            fi
        else
            echo "Not running"
        fi
    ;;
    restart)
        $0 stop
        if is_running; then
            echo "Unable to stop, will not attempt to start"
            exit 1
        fi
        $0 start
    ;;
    status)
        if is_running; then
            echo "Running"
        else
            echo "Stopped"
            exit 1
        fi
    ;;
    *)
    echo "Usage: $0 {start|stop|restart|status}"
    exit 1
    ;;
esac
exit 0
`
