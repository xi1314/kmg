package kmgConsole

import (
	"os"

	"github.com/bronze1man/kmg/kmgCmd"
)

func ProxyCommand(args ...string) func() {
	return func() {
		kmgCmd.CmdSlice(append(args, os.Args...)).ProxyRun()
	}
}
