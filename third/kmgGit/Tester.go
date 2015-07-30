package kmgGit

import (
	"os"

	"github.com/bronze1man/kmg/kmgFile"
)

func GitTestCb(f func()) {
	oldwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	defer func() {
		os.Chdir(oldwd)
		kmgFile.MustDelete("testFile")
	}()
	kmgFile.MustDelete("testFile")
	kmgFile.MustMkdir("testFile")
	os.Chdir("testFile")
	f()
}
