package internal

import (
	"github.com/bronze1man/kmg/kmgConsole"
)

func init() {
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "SelfUpdate",
		Desc:   "update kmg tool from our server",
		Runner: selfUpdate,
	})
}
