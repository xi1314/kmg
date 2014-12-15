package command

import (
	"github.com/bronze1man/kmg/kmgConsole"
)

func exitOnErr(err error) {
	kmgConsole.ExitOnErr(err)
	return
}
