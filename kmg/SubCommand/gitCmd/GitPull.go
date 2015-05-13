package gitCmd

import (
	"github.com/bronze1man/kmg/kmgCmd"
)

func gitPull() {
	kmgCmd.ProxyRun("git pull origin master")
}
