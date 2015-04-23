package internal

import (
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgConsole"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgTime"
	"os"
	"os/exec"
	"time"
)

func init() {
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "_SelfInstall",
		Desc:   "install kmg in this computer(should call in install bash)",
		Runner: selfInstallCmd,
		Hidden: true,
	})
}

//由于bash过度难用,直接安装kmg的时候又会遇到很复杂的情况,此处用于处理某些复杂情况
func selfInstallCmd() {
	//kmg 被装到了不是 /bin/kmg 的目录
	path, err := exec.LookPath("kmg")
	if err != nil && !os.IsNotExist(err) {
		kmgConsole.ExitOnErr(err)
	}
	if err == nil && path != "/bin/kmg" {
		backPathDir := "/var/backup/kmg/" + time.Now().Format(kmgTime.FormatFileName)
		kmgFile.MustMkdirAll(backPathDir)
		kmgCmd.MustRun("mv /bin/kmg " + backPathDir)
	}
}
