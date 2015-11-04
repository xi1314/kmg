package SubCommand

import (
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgConsole"
	"github.com/bronze1man/kmg/kmgFile"
)

func AddSelfInstallCommand() {
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "_SelfInstall",
		Desc:   "install kmg in this computer(should call in install bash)",
		Runner: selfInstallCmd,
		Hidden: true,
	})
}

//由于bash过度难用,直接安装kmg的时候又会遇到很复杂的情况,此处用于处理某些复杂情况
func selfInstallCmd() {
	//迁移代码,从/bin/kmg迁移到/usr/local/bin/kmg
	if kmgFile.MustFileExist("/bin/kmg") {
		kmgCmd.MustRunInBash("sudo rm /bin/kmg;hash -r")
	}
	kmgFile.MustEnsureBinPath("/usr/local/bin/kmg")
}
