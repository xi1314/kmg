package kmgSsh

import (
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgTime"
	"os/exec"
	"strings"
	"time"
)

func MustRpcSshCmd(ip string, cmd ...string) []byte {
	if len(cmd) == 0 {
		return []byte{}
	}
	if ip == "" {
		return []byte{}
	}
	cmdCombine := strings.Join(cmd, "&&")
	out, err := kmgCmd.CmdString("ssh -o StrictHostKeyChecking=no -o PasswordAuthentication=no root@" + ip + " " + cmdCombine).RunAndReturnOutput()
	logPath := "log/rpcSshCmd-" + ip
	kmgFile.MustAppendFile(logPath, []byte(strings.Join([]string{cmdCombine, kmgTime.DefaultFormat(time.Now())}, "\n")))
	kmgFile.MustAppendFile(logPath, out)
	if err != nil {
		_, ok := err.(*exec.ExitError)
		if ok {
			return out
		}
		kmgFile.MustAppendFile(logPath, []byte(err.Error()))
		panic(err)
	}
	return out
}
