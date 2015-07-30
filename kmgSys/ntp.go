package kmgSys

import (
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgPlatform"
)

//目前只支持 Ubuntu
func SyncTime() {
	currentPlatform := kmgPlatform.GetCompiledPlatform()
	if currentPlatform.Os != kmgPlatform.LinuxAmd64.Os {
		return
	}
	if !kmgCmd.Exist("ntpdate") {
		kmgCmd.MustRun("apt-get install -y ntpdate")
	}
	kmgCmd.MustRun("ntpdate -u ntp.api.bz")
}
