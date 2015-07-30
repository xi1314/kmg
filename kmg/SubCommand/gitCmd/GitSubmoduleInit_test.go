package gitCmd

import (
	"testing"

	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgTest"
	"github.com/bronze1man/kmg/third/kmgGit"
)

func TestGitSubmoduleInitIgnore(ot *testing.T) {
	kmgGit.GitTestCb(func() {
		kmgCmd.MustRun("git init")
		kmgFile.MustWriteFile(".gitignore", []byte("/subIgnored"))
		kmgFile.MustWriteFileWithMkdir("subIgnored/1.txt", []byte("1"))
		kmgCmd.CmdString("git init").SetDir("subIgnored").MustRun()
		kmgCmd.CmdString("git add -A").SetDir("subIgnored").MustRun()
		kmgCmd.CmdString("git commit -am'save'").SetDir("subIgnored").MustRun()

		repo := kmgGit.MustGetRepositoryFromPath(".")
		GitSubmoduleInit(repo)
		kmgTest.Equal(repo.MustIsFileInIndex("subIgnored"), false)
	})
}

func TestGitSubmoduleInitNotInIndex(ot *testing.T) {
	kmgGit.GitTestCb(func() {
		kmgCmd.MustRun("git init")
		kmgFile.MustWriteFileWithMkdir("sub/1.txt", []byte("1"))
		kmgCmd.CmdString("git init").SetDir("sub").MustRun()
		kmgCmd.CmdString("git add -A").SetDir("sub").MustRun()
		kmgCmd.CmdString("git commit -am'save'").SetDir("sub").MustRun()

		repo := kmgGit.MustGetRepositoryFromPath(".")
		GitSubmoduleInit(repo)
		kmgTest.Equal(repo.MustIsFileInIndex("sub/1.txt"), true)
	})
}
