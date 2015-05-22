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

func (repo *Repository) MustCurrentBranchName() string {
	output := kmgCmd.CmdString("git rev-parse --abbrev-ref HEAD").SetDir(repo.gitPath).MustRunAndReturnOutput()
	return strings.TrimSpace(string(output))
}
