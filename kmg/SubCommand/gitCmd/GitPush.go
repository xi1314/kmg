package gitCmd

import (
	"github.com/bronze1man/kmg/kmgCmd"
)

func gitPush() {
	kmgCmd.ProxyRun("git add -A")
	kmgCmd.ProxyRun("git commit -am'save'")
	kmgCmd.ProxyRun("git push origin master")
}
