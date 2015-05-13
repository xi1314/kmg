package SubCommand

import (
	"github.com/bronze1man/kmg/kmgConsole"
)

func AddSelfUpdate() {
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "SelfUpdate",
		Desc:   "update kmg tool from our server",
		Runner: selfUpdate,
	})
}
