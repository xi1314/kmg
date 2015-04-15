package internal

import (
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgConsole"
	"os"
	"strings"
)

func init() {
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "Make",
		Desc:   "run a project defined command",
		Runner: makeCmd,
	})

}

func makeCmd() {
	kmgc, err := kmgConfig.LoadEnvFromWd()
	kmgConsole.ExitOnErr(err)
	if kmgc.Make==""{
		kmgConsole.ExitOnStderr("Please defined a Make in .kmg.yml file to use kmg make")
		return
	}
	args:=strings.Split(kmgc.Make," ")
	err=kmgCmd.CmdSlice(append(args,os.Args[1:]...)).Run()
	kmgConsole.ExitOnErr(err)
}
