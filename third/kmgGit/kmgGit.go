package kmgGit

import (
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgFile"
	"os"
	"strings"
	"sync"
)

func GetRepositoryFromPath(path string) (repo *Repository, err error) {
	p, err := kmgFile.SearchFileInParentDir(path, ".git")
	if err != nil {
		return
	}
	return &Repository{
		gitPath: p,
	}, nil
}

var defaultRepoOnce sync.Once
var defaultRepo *Repository

func DefaultRepository() (repo *Repository) {
	defaultRepoOnce.Do(func() {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		defaultRepo, err = GetRepositoryFromPath(wd)
		if err != nil {
			panic(err)
		}
	})
	return defaultRepo
}

type Repository struct {
	gitPath string
}

func (repo *Repository) MustGetCurrentBranchName() string {
	output := kmgCmd.CmdString("git rev-parse --abbrev-ref HEAD").SetDir(repo.gitPath).MustCombinedOutput()
	return strings.TrimSpace(string(output))
}

func (repo *Repository) MustGetIndexFileList() []string {
	output := kmgCmd.CmdString("git ls-files").SetDir(repo.gitPath).MustCombinedOutput()
	outputSlice:=[]string{}
	for _,s:=range strings.Split(string(output),"\n"){
		s=strings.TrimSpace(s)
		if s==""{
			continue
		}
		outputSlice = append(outputSlice,s)
	}
	return outputSlice
}

func (repo *Repository) MustIndexRemoveByPath(path string) {
	kmgCmd.CmdSlice([]string{"git", "rm", "--cached","-r",path}).MustStdioRun()
}