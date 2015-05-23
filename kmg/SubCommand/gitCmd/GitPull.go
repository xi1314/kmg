package gitCmd

import (
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/third/kmgGit"
)

func gitPull() {
	branchName := kmgGit.DefaultRepository().MustGetCurrentBranchName()
	kmgCmd.ProxyRun("git pull origin " + branchName)
}
