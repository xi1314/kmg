package kmgSsh

import (
	"github.com/bronze1man/kmg/kmgCmd"
	"strings"
)

//每次尝试连接，5秒超时，超时后，重试12次，最多等1分钟
//若 isReachable = false，则 havePermission 没有有意义
func AvailableCheck(remote *RemoteServer) (isReachable, havePermission bool) {
	retry := 0
	for {
		cmd := kmgCmd.CmdString("ssh -o StrictHostKeyChecking=no -o PasswordAuthentication=no -o ConnectTimeout=5 " + remote.String() + " echo ok")
		b, e := cmd.CombinedOutput()
		if e == nil && strings.HasPrefix(string(b), "ok") {
			return true, true
		}
		if e != nil && strings.Contains(string(b), "Permission denied") {
			return true, false
		}
		retry++
		if retry == 12 {
			return false, false
		}
	}
	return false, false
}
