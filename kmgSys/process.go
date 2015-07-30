package kmgSys

import (
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgStrconv"
	"os"
	"strings"
)

func KillProcessByBinName(binName string) {
	pidList := GetAllProcessIdByBinName(binName)
	for _, pid := range pidList {
		KillProcessByPid(pid)
	}
}

//只支持 Linux
func GetAllProcessIdByBinName(binName string) []int {
	b := kmgCmd.CmdSlice([]string{"ps", "-C", binName, "-o", "pid="}).MustRunAndReturnOutput()
	s := strings.Split(string(b), "\n")
	out := []int{}
	for _, l := range s {
		out = append(out, kmgStrconv.AtoIDefault0(strings.TrimSpace(l)))
	}
	return out
}

func KillProcessByPid(pid int) {
	p, err := os.FindProcess(pid)
	handleErr(err)
	err = p.Kill()
	handleErr(err)
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
