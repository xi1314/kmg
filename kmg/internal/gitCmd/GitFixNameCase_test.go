// +build darwin

package gitCmd

import (
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgTest"
	"os"
	"path/filepath"
	"testing"
)

func TestGitFixNameCaseWithFile(ot *testing.T) {
	oldWd, err := os.Getwd()
	kmgTest.Equal(err, nil)
	kmgFile.MustDelete("testfile")
	kmgFile.Mkdir("testfile")
	os.Chdir("testfile")
	defer os.Chdir(oldWd)
	kmgFile.MustWriteFile("a.txt", []byte("abc"))
	kmgCmd.ProxyRun("git init")
	kmgCmd.ProxyRun("git add -A")
	kmgCmd.ProxyRun("git commit -am'save'")
	err = os.Rename("a.txt", "A.txt")
	kmgTest.Equal(err, nil)

	err = GitFixNameCase(filepath.Join(oldWd, "testfile"))
	kmgTest.Equal(err, nil)

	kmgCmd.ProxyRun("git status")
	kmgCmd.ProxyRun("git add -A")
	kmgCmd.ProxyRun("git commit -am'save'")
}

func TestGitFixNameCaseWithDirectory(ot *testing.T) {
	oldWd, err := os.Getwd()
	kmgTest.Equal(err, nil)
	kmgFile.MustDelete("testfile")
	kmgFile.MustWriteFileWithMkdir("testfile/a/a.txt", []byte("abc"))
	os.Chdir("testfile")
	defer os.Chdir(oldWd)
	kmgCmd.ProxyRun("git init")
	kmgCmd.ProxyRun("git add -A")
	kmgCmd.ProxyRun("git commit -am'save'")
	err = os.Rename("a", "A")
	kmgTest.Equal(err, nil)

	err = GitFixNameCase(filepath.Join(oldWd, "testfile"))
	kmgTest.Equal(err, nil)

	kmgCmd.ProxyRun("git status")
	kmgCmd.ProxyRun("git add -A")
	kmgCmd.ProxyRun("git commit -am'save'")
}
