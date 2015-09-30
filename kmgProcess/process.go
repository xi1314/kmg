package kmgProcess

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgErr"
	"github.com/bronze1man/kmg/kmgPlatform"
	"github.com/bronze1man/kmg/kmgStrconv"
	"github.com/bronze1man/kmg/kmgSys"
	"os"
	"strconv"
	"strings"
)

type Process struct {
	Id       int
	Command  string
	StartCmd string
}

//杀不死就算了
func (p *Process) Kill() {
	sysProcess, err := os.FindProcess(p.Id)
	kmgErr.LogErrorWithStack(err)
	err = sysProcess.Kill()
	kmgErr.LogErrorWithStack(err)
}

func (p *Process) CmdKill() string {
	return "kill " + strconv.Itoa(p.Id)
}

//psOutput 样例
// 1234 /bin/hi -l :1024
//12345 /bin/world -n test
func Extract(psOutput string) []*Process {
	lines := strings.Split(psOutput, "\n")
	out := []*Process{}
	for _, l := range lines {
		ls := strings.Fields(l)
		if len(ls) == 0 {
			continue
		}
		p := &Process{
			Id:      kmgStrconv.AtoIDefault0(ls[0]),
			Command: strings.Join(ls[1:], " "),
		}
		p.StartCmd = "setsid " + p.Command
		out = append(out, p)
	}
	return out
}

// 时间复杂度 2n^2
func Diff(expect, running []*Process) (notExpect, notRunning []*Process) {
	aHaveBNotHave := func(a, b []*Process) []*Process {
		bNotHave := []*Process{}
		for _, ap := range a {
			isMatch := false
			for _, bp := range b {
				if bp.Command == ap.Command {
					isMatch = true
					continue
				}
			}
			if isMatch {
				continue
			}
			bNotHave = append(bNotHave, ap)
		}
		return bNotHave
	}
	return aHaveBNotHave(running, expect), aHaveBNotHave(expect, running)
}

//只兼容Linux，请用下面的那个命令
func CmdProcessByBinName(binName string) string {
	return fmt.Sprintf("ps -C %s -o pid=,cmd=", binName)
}

//以特定格式化的形式列出所有进程，方便提取进程信息
//第一列是进程 ID，第二列是进程完整的执行命令（执行文件或命令+全部执行参数）
//兼容：Linux, OS X
func CmdAllProcess() string {
	return "ps ax -o pid=,args="
}

func LinuxGetAllProcessByBinName(binName string) []*Process {
	if !kmgPlatform.LinuxAmd64.Compatible(kmgPlatform.GetCompiledPlatform()) {
		panic(kmgSys.ErrPlatformNotSupport)
	}
	b, err := kmgCmd.CmdBash(CmdProcessByBinName(binName)).GetExecCmd().CombinedOutput()
	if err != err {
		fmt.Println(err)
	}
	return Extract(string(b))
}
