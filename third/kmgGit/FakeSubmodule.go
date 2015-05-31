package kmgGit

import (
	"fmt"
	"github.com/bronze1man/kmg/encoding/kmgJson"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgRand"
	"os"
	"path/filepath"
)

type SubRepositoryInfo struct {
	RemoteUrl string
	CommitId  string
}

//把当前项目里面的所有的子项目都变成伪submodule,并且保存版本信息到.gitFakeSubmodule 中
func (repo *Repository) MustFakeSubmoduleCommit() {
	//把当前目录下面的所有叫.git的目录都翻出来
	rootPath := repo.gitPath
	SubmoduleList := map[string]SubRepositoryInfo{}
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if filepath.Base(path) != ".git" {
			return nil
		}
		path = filepath.Dir(path)
		if path == rootPath {
			return nil
		}
		relPath, err := filepath.Rel(rootPath, path)
		if err != nil {
			panic(err)
		}
		ret := repo.MustFakeSubmoduleAdd(relPath)
		if !ret {
			return nil
		}
		//记录版本号
		//TODO 确认origin/master有这个commitid,如果没有报一个warning.
		subRepo := MustGetRepositoryFromPath(path)
		commitId := subRepo.MustGetHeadCommitId()
		if !subRepo.MustIsInParent("origin/master", commitId) {
			fmt.Printf("warning: [%s] HEAD is not in origin/master\n", relPath)
		}
		SubmoduleList[relPath] = SubRepositoryInfo{
			RemoteUrl: subRepo.MustGetRemoteUrl("origin"),
			CommitId:  commitId,
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	kmgJson.MustWriteFileIndent(filepath.Join(rootPath, ".gitFakeSubmodule"), SubmoduleList)
	return
}

//从 .gitFakeSubmodule 中还原旧的fakeSubmoudle,并且将所有子项目都切换到该文件里面写的分支(使用reset切分支,保证不掉数据)
func (repo *Repository) MustFakeSubmoduleUpdate() {
	rootPath := repo.gitPath
	SubmoduleList := map[string]SubRepositoryInfo{}
	kmgJson.MustReadFile(filepath.Join(rootPath, ".gitFakeSubmodule"), &SubmoduleList)
	for repoPath, SubmoduleInfo := range SubmoduleList {
		repoRealPath := filepath.Join(rootPath, repoPath)
		//子项目存在?
		if !MustIsRepositoryAtPath(repoRealPath) {
			tmpPath := filepath.Join(repoRealPath, kmgRand.MustCryptoRandToHex(8))
			kmgFile.MustMkdir(tmpPath)
			MustGitClone(SubmoduleInfo.RemoteUrl, tmpPath)
			kmgFile.MustRename(filepath.Join(tmpPath, ".git"), filepath.Join(repoRealPath, ".git"))
			kmgFile.MustDelete(tmpPath)
		}
		//子项目的远程路径正确?
		subRepo := MustGetRepositoryFromPath(repoRealPath)
		if subRepo.MustGetRemoteUrl("origin") != SubmoduleInfo.RemoteUrl {
			subRepo.MustSetRemoteUrl("origin", SubmoduleInfo.RemoteUrl)
		}
		//子项目的版本号正确?
		if subRepo.MustGetHeadCommitId() != SubmoduleInfo.CommitId {
			subRepo.MustResetToCommitId(SubmoduleInfo.CommitId)
		}
	}
}

//这个函数返回 经过处理后该路径是否是一个伪submodule
// 请调用者保证这个path里面有.git文件夹
func (repo *Repository) MustFakeSubmoduleAdd(path string) bool {
	//已经是伪submodule
	if repo.MustHasFilesInDirInIndex(path) {
		return true
	}
	//1.被忽略,此处也忽略,什么也不做
	if repo.MustIsFileIgnore(path) {
		return false
	}
	//2.是真submodule
	if repo.MustIsFileInIndex(path) {
		repo.MustIndexRemoveByPath(path)
	}
	//3.没有加入Index里面
	subrepo := MustGetRepositoryFromPath(filepath.Join(repo.gitPath, path))
	if subrepo.gitPath == repo.gitPath {
		panic("[MustFakeSubmoduleAdd] input path do not have .git file")
	}
	fileList := subrepo.MustGetIndexFileList()
	if len(fileList) == 0 {
		panic("[MustFakeSubmoduleAdd] submodule do not have any file in index")
	}

	//加入该submodule里面任意一个文件(速度快)
	for _, file := range fileList {
		repo.MustIndexAddFile(filepath.Join(path, file))
		return true
	}
	return true
}
