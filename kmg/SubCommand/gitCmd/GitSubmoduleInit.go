package gitCmd

import (
	"github.com/bronze1man/kmg/kmgConsole"
	"github.com/bronze1man/kmg/third/kmgGit"
	"fmt"
	"flag"
)

// fake     //伪造submodule
// add      //添加伪submodule
// recover  //恢复伪submodule(可能要拉项目)
// status
func init() {
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "GitSubmoduleCommit",
		Desc:   "git fake submodule Commit",
		Runner: func() { kmgGit.DefaultRepository().MustFakeSubmoduleCommit() },
	})
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "GitSubmoduleUpdate",
		Desc:   "git fake submodule Update",
		Runner: func() { kmgGit.DefaultRepository().MustFakeSubmoduleUpdate() },
	})
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "GitSmallestChange",
		Runner: func() {
			Local:=""
			Target:=""
			flag.StringVar(&Local,"local","","localcommit id (sha256 commit id or branch name or HEAD)")
			flag.StringVar(&Target,"target","","targetCommit id (sha256 commit id or branch name or HEAD)")
			flag.Parse()
			if Local=="" || Target==""{
				flag.Usage()
				return
			}
			result:=kmgGit.DefaultRepository().MustSmallestChange(Local,Target)
			fmt.Println( result )
			fmt.Println("#see diff: git diff --stat "+Local+" "+result)
		},
	})
}
