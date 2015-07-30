package gitCmd

import (
	"flag"

	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/third/kmgGit"
)

func gitPull() {
	flag.Parse()
	remote := flag.Arg(0)
	if remote == "" {
		remote = "origin"
	}
	branchName := kmgGit.DefaultRepository().MustGetCurrentBranchName()
	kmgCmd.ProxyRun("git pull " + remote + " " + branchName)
}
