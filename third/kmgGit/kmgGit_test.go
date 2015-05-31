package kmgGit

import (
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
)

func TestMustIsFileIgnore(t *testing.T) {
	GitTestCb(func() {
		kmgCmd.MustRun("git init")
		kmgFile.MustWriteFile(".gitignore", []byte("/1.txt"))
		kmgFile.MustWriteFile("1.txt", []byte("1"))
		kmgFile.MustWriteFile("2.txt", []byte("1"))

		repo := MustGetRepositoryFromPath(".")
		kmgTest.Equal(repo.MustIsFileIgnore("1.txt"), true)
		kmgTest.Equal(repo.MustIsFileIgnore("notExist.txt"), false)
		kmgTest.Equal(repo.MustIsFileIgnore("2.txt"), false)
	})
}

func TestMustIsFileInIndex(t *testing.T) {
	GitTestCb(func() {
		kmgCmd.MustRun("git init")
		kmgFile.MustWriteFile("1.txt", []byte("1"))
		kmgFile.MustWriteFile("2.txt", []byte("1"))
		kmgCmd.MustRun("git add 1.txt")
		repo := MustGetRepositoryFromPath(".")
		kmgTest.Equal(repo.MustIsFileInIndex("1.txt"), true)
		kmgTest.Equal(repo.MustIsFileInIndex("2.txt"), false)
	})
}
