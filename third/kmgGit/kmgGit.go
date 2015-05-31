package kmgGit

import (
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgStrings"
	"os"
	"path/filepath"
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

func MustGetRepositoryFromPath(path string) (repo *Repository) {
	p, err := kmgFile.SearchFileInParentDir(path, ".git")
	if err != nil {
		panic(err)
	}
	return &Repository{
		gitPath: p,
	}
}

func MustIsRepositoryAtPath(path string) bool {
	return kmgFile.MustFileExist(filepath.Join(path, ".git"))
}

func MustGitClone(url string, path string) {
	kmgCmd.CmdSlice([]string{"git", "clone", url, path}).MustRun()
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

func (repo *Repository) GetGitPath() string {
	return repo.gitPath
}
func (repo *Repository) MustGetCurrentBranchName() string {
	output := kmgCmd.CmdString("git rev-parse --abbrev-ref HEAD").SetDir(repo.gitPath).MustCombinedOutput()
	return strings.TrimSpace(string(output))
}

func (repo *Repository) MustGetIndexFileList() []string {
	output := kmgCmd.CmdString("git ls-files").SetDir(repo.gitPath).MustCombinedOutput()
	outputSlice := []string{}
	for _, s := range strings.Split(string(output), "\n") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		outputSlice = append(outputSlice, s)
	}
	return outputSlice
}

func (repo *Repository) MustGetHeadCommitId() string {
	output := kmgCmd.CmdString("git rev-parse HEAD").SetDir(repo.gitPath).MustCombinedOutput()
	return strings.TrimSpace(string(output))
}
func (repo *Repository) MustResetToCommitId(commitId string) {
	kmgCmd.CmdSlice([]string{"git", "reset", commitId}).SetDir(repo.gitPath).MustStdioRun()
}
func (repo *Repository) MustIndexRemoveByPath(path string) {
	kmgCmd.CmdSlice([]string{"git", "rm", "--cached", "-r", path}).SetDir(repo.gitPath).MustStdioRun()
}

func (repo *Repository) MustIndexAddFile(path string) {
	kmgCmd.CmdSlice([]string{"git", "add", path}).SetDir(repo.gitPath).MustStdioRun()
}

func (repo *Repository) MustIsFileIgnore(path string) bool {
	return kmgCmd.CmdSlice([]string{"git", "check-ignore", path}).SetDir(repo.gitPath).
		MustHiddenRunAndGetExitStatus() == 0
}

func (repo *Repository) MustGetAllParentCommitId(commitId string) []string {
	output := kmgCmd.CmdSlice([]string{"git", "log", "--format=%H", commitId}).SetDir(repo.gitPath).MustCombinedOutput()
	outputSlice := []string{}
	for _, s := range strings.Split(string(output), "\n") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		outputSlice = append(outputSlice, s)
	}
	return outputSlice
}

func (repo *Repository) MustIsInParent(commitId string, parentCommitId string) bool {
	allParent := repo.MustGetAllParentCommitId(commitId)
	return kmgStrings.IsInSlice(allParent, parentCommitId)
}

//是否这个路径指向的东西在index里面存在(确切的文件)
func (repo *Repository) MustIsFileInIndex(path string) bool {
	output := kmgCmd.CmdSlice([]string{"git", "ls-files", path}).SetDir(repo.gitPath).MustCombinedOutput()
	if len(output) == 0 {
		return false
	}
	for _, s := range strings.Split(string(output), "\n") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if s == path {
			return true
		}
	}
	return false
}

//返回index里面某个文件路径下面是否有文件,如果它本身在index里面返回false
func (repo *Repository) MustHasFilesInDirInIndex(path string) bool {
	output := kmgCmd.CmdSlice([]string{"git", "ls-files", path}).SetDir(repo.gitPath).MustCombinedOutput()
	if len(output) == 0 {
		return false
	}
	hasElem := false
	for _, s := range strings.Split(string(output), "\n") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		hasElem = true
		if s == path {
			return false
		}
	}
	return hasElem
}

func (repo *Repository) MustGetRemoteUrl(branch string) string {
	output := kmgCmd.CmdSlice([]string{"git", "ls-remote", "--get-url", branch}).SetDir(repo.gitPath).MustCombinedOutput()
	return strings.TrimSpace(string(output))
}

func (repo *Repository) MustSetRemoteUrl(branch string, url string) {
	kmgCmd.CmdSlice([]string{"git", "remote", "add", branch, url}).SetDir(repo.gitPath).MustStdioRun()
}
