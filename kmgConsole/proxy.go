package kmgConsole

import (
	"github.com/bronze1man/kmg/kmgCmd"
	"os"
)

func ProxyCommand(args ...string) func() {
	return func() {
		kmgCmd.CmdSlice(append(args, os.Args...)).ProxyRun()
	}
}
