package kmgNet

import (
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgPlatform"
	"strings"
)

func GetDefaultGateway() string {
	if !kmgPlatform.IsDarwin() {
		panic("not support platform")
	}
	output := kmgCmd.MustCombinedOutput("netstat -nr")
	for _, line := range strings.Split(string(output), "\n") {
		if strings.Contains(line, "default") {
			return strings.Fields(line)[1]
		}
	}
	return ""
}

func SetDnsServerAddr(ip string) {
	if !kmgPlatform.IsDarwin() {
		panic("not support platform")
	}
	kmgCmd.MustRun("networksetup -setdnsservers Wi-Fi " + ip)
}
