package kmgGit

import (
	"testing"

	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgTest"
)

func TestGitSubmoduleAddIgnore(ot *testing.T) {
	GitTestCb(func() {
		kmgCmd.MustRun("git init")
		kmgFile.MustWriteFile(".gitignore", []byte("/subIgnored"))
		kmgFile.MustWriteFileWithMkdir("subIgnored/1.txt", []byte("1"))
		kmgCmd.CmdString("git init").SetDir("subIgnored").MustRun()
		kmgCmd.CmdString("git add -A").SetDir("subIgnored").MustRun()
		kmgCmd.CmdString("git commit -am'save'").SetDir("subIgnored").MustRun()

		repo := MustGetRepositoryFromPath(".")
		repo.MustFakeSubmoduleAdd("subIgnored")
		kmgTest.Equal(repo.MustIsFileInIndex("subIgnored"), false)
	})
}

func TestGitSubmoduleAddNotInIndex(ot *testing.T) {
	GitTestCb(func() {
		kmgCmd.MustRun("git init")
		kmgFile.MustWriteFileWithMkdir("sub/1.txt", []byte("1"))
		kmgCmd.CmdString("git init").SetDir("sub").MustRun()
		kmgCmd.CmdString("git add -A").SetDir("sub").MustRun()
		kmgCmd.CmdString("git commit -am'save'").SetDir("sub").MustRun()

		repo := MustGetRepositoryFromPath(".")
		repo.MustFakeSubmoduleAdd("sub")
		kmgTest.Equal(repo.MustIsFileInIndex("sub/1.txt"), true)
	})
}

func TestGitSubmoduleAddRealSubmodule(ot *testing.T) {
	GitTestCb(func() {
		kmgCmd.MustRun("git init")
		kmgFile.MustWriteFileWithMkdir("sub/1.txt", []byte("1"))
		kmgCmd.CmdString("git init").SetDir("sub").MustRun()
		kmgCmd.CmdString("git add -A").SetDir("sub").MustRun()
		kmgCmd.CmdString("git commit -am'save'").SetDir("sub").MustRun()

		kmgCmd.MustRun("git add -A")
		kmgCmd.MustRun("git commit -am'save'")

		repo := MustGetRepositoryFromPath(".")
		kmgTest.Equal(repo.MustIsFileInIndex("sub/1.txt"), false)
		kmgTest.Equal(repo.MustIsFileInIndex("sub"), true)

		repo.MustFakeSubmoduleAdd("sub")
		kmgTest.Equal(repo.MustIsFileInIndex("sub"), false)
		kmgTest.Equal(repo.MustIsFileInIndex("sub/1.txt"), true)
	})
}
