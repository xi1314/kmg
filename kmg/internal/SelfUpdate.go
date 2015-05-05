package internal

import (
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgConsole"
)

func init() {
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "SelfUpdate",
		Desc:   "update kmg tool from our server",
		Runner: selfUpdate,
	})
}

func selfUpdate() {
	kmgCmd.MustRunInBash("curl http://kmgtools.qiniudn.com/v1/installKmg.bash?v=$RANDOM | bash")
}
