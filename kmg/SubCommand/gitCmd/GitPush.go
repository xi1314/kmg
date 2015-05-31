package gitCmd

import (
	"flag"
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/third/kmgGit"
)

//把当前分支推到origin的当前分支,当前分支不一定是master
func gitPush() {
	flag.Parse()
	remote := flag.Arg(0)
	if remote == "" {
		remote = "origin"
	}
	branchName := kmgGit.DefaultRepository().MustGetCurrentBranchName()
	kmgCmd.MustRunNotExistStatusCheck("git add -A")
	kmgCmd.MustRunNotExistStatusCheck("git commit -am save")
	kmgCmd.MustRunNotExistStatusCheck("git push " + remote + " " + branchName)
}
