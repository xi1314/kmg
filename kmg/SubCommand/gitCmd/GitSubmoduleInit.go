package gitCmd

import (
	"github.com/bronze1man/kmg/kmgConsole"
	"github.com/bronze1man/kmg/third/kmgGit"
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
}
