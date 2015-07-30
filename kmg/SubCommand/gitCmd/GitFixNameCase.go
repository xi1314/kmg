package gitCmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bronze1man/kmg/kmgConsole"
	"github.com/bronze1man/kmg/third/kmgGit"
)

func init() {
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "GitFixNameCase",
		Desc:   "fix git name case problem on case insensitive opearate system(windows,osx)",
		Runner: gitFixNameCaseCmd,
	})
}

func gitFixNameCaseCmd() {
	//检查index里面的文件名大小写和当前的文件名大小写是否一致
	var basePath string
	var err error
	flag.StringVar(&basePath, "p", "", "base path of git directory")
	flag.Parse()
	if basePath == "" {
		basePath, err = os.Getwd()
		kmgConsole.ExitOnErr(err)
	}
	err = GitFixNameCase(basePath)
	kmgConsole.ExitOnErr(err)
}

func GitFixNameCase(basePath string) (err error) {
	folderFileCache = map[string][]string{}
	repo, err := kmgGit.GetRepositoryFromPath(basePath)
	if err != nil {
		return err
	}
	caseDiffChangeArray := []string{}
	for _, indexPath := range repo.MustGetIndexFileList() {
		//fullPath := filepath.Join(basePath, ie.Path)
		isSameOrNotExist := checkOneFileFoldDiff(basePath, indexPath)
		if isSameOrNotExist {
			continue
		}
		//在index里面修复大小写错误
		caseDiffChangeArray = append(caseDiffChangeArray, indexPath)
	}

	if len(caseDiffChangeArray) > 0 {
		fmt.Println("file name diff in case:")
		for _, refPath := range caseDiffChangeArray {
			fmt.Println("\t", refPath)
			repo.MustIndexRemoveByPath(refPath)
			//index.RemoveByPath(refPath)
		}
	}
	return nil
}

type caseDiffChange struct {
	gitPath string
}

var folderFileCache = map[string][]string{}

func checkOneFileFoldDiff(basePath string, refPath string) (isSameOrNotExist bool) {
	//此处要检查文件的每一部分的fold都一致
	filePathPartList := strings.Split(refPath, "/")
	for i := range filePathPartList {
		ret := checkOneBasePathFoldDiff(filepath.Join(basePath, strings.Join(filePathPartList[:i+1], "/")))
		if !ret {
			return false
		}
	}
	return true
}
func checkOneBasePathFoldDiff(path string) (isSameOrNotExist bool) {
	fileName := filepath.Base(path)
	dirPath := filepath.Dir(path)
	names, ok := folderFileCache[dirPath]
	if !ok {
		dirFile, err := os.Open(filepath.Dir(path))
		defer dirFile.Close()
		if err != nil {
			if os.IsNotExist(err) {
				return true //dir not exist
			}
			kmgConsole.ExitOnErr(err)
		}

		names, err = dirFile.Readdirnames(-1)
		kmgConsole.ExitOnErr(err)
		folderFileCache[dirPath] = names
	}
	for _, n := range names {
		if n == fileName {
			return true //case same
		}
		if strings.EqualFold(n, fileName) {
			return false //case not same
		}
	}
	return true // base path not exist
}
