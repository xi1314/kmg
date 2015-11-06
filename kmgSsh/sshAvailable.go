package kmgSsh

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgCmd"
	"strings"
	"time"
)

//每次尝试连接，5秒超时，超时后，重试12次，最多等1分钟
//若 isReachable = false，则 havePermission 没有有意义
func AvailableCheck(remote *RemoteServer) (isReachable, havePermission bool) {
	retry := 0
	for {
		cmd := kmgCmd.CmdString("ssh -o StrictHostKeyChecking=no -o PasswordAuthentication=no -o ConnectTimeout=5 " + remote.String() + " echo ok")
		start := time.Now()
		b, e := cmd.CombinedOutput()
		delta := time.Now().Sub(start)
		fmt.Println("[kmgSsh AvailableCheck] ssh -o StrictHostKeyChecking=no -o PasswordAuthentication=no -o ConnectTimeout=5 " + remote.String() + " echo ok")
		fmt.Println("[kmgSsh AvailableCheck]", remote.Address, string(b))
		if e == nil && strings.HasPrefix(string(b), "ok") {
			return true, true
		}
		if e != nil && strings.Contains(string(b), "Permission denied") {
			return true, false
		}
		delta -= time.Second * 5
		if delta < 0 {
			time.Sleep(-1 * delta)
		}
		retry++
		fmt.Println("[kmgSsh AvailableCheck]", remote.Address, "retry", retry)
		if retry == 24 {
			return false, false
		}
	}
	return false, false
}
