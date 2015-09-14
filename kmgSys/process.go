package kmgSys

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgPlatform"
	"github.com/bronze1man/kmg/kmgStrconv"
	"os"
	"strings"
)

type Process struct {
	Id      int
	Command string
}

func (p *Process) Kill() {
	sp, err := os.FindProcess(p.Id)
	handleErr(err)
	err = sp.Kill()
	handleErr(err)
}

func GetAllProcessByBinName(binName string) []*Process {
	if !kmgPlatform.LinuxAmd64.Compatible(kmgPlatform.GetCompiledPlatform()) {
		panic(ErrPlatformNotSupport)
	}
	b, err := kmgCmd.CmdBash(CmdProcessListByBinName(binName)).GetExecCmd().CombinedOutput()
	if err != err {
		fmt.Println(err)
	}
	return ExtractProcessListFromString(string(b))
}

func CmdKillProcessByPid(pid int) string {
	return fmt.Sprintf("kill %v", pid)
}

func CmdProcessListByBinName(binName string) string {
	return fmt.Sprintf("ps -C %s -o pid=,cmd=", binName)
}

func ExtractProcessListFromString(output string) []*Process {
	lines := strings.Split(output, "\n")
	out := []*Process{}
	for _, l := range lines {
		ls := strings.Fields(l)
		if len(ls) == 0 {
			continue
		}
		out = append(out, &Process{
			Id:      kmgStrconv.AtoIDefault0(ls[0]),
			Command: strings.Join(ls[1:], " "),
		})
	}
	return out
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
