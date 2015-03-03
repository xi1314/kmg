package goCommand

/*
import (
    "github.com/bronze1man/kmg/kmgConsole"
    "flag"
    "os"
    git2go "github.com/libgit2/git2go"
    "path/filepath"
    "fmt"
    "strings"
)

func init() {
    kmgConsole.AddAction(kmgConsole.Command{
        Name:   "GitFixNameCase",
        Desc:   "fix git name case problem on case insensitive opearate system(windows,osx)",
        Runner: GitFixNameCase,
    })

}

func GitFixNameCase() {
    //检查index里面的文件名大小写和当前的文件名大小写是否一致
    var basePath string
    var err error
    flag.StringVar(&basePath, "p", "", "base path of git directory")
    flag.Parse()
    if basePath == "" {
        basePath, err = os.Getwd()
        kmgConsole.ExitOnErr(err)
    }
    repo, err := git2go.OpenRepository(basePath)
    kmgConsole.ExitOnErr(err)
    index, err := repo.Index()
    kmgConsole.ExitOnErr(err)
    indexCount := index.EntryCount()
    caseDiffChangeArray := []caseDiffChange{}
    for i := uint(0); i < indexCount; i++ {
        ie, err := index.EntryByIndex(i)
        kmgConsole.ExitOnErr(err)
        fullPath := filepath.Join(basePath, ie.Path)
        exist, actualFileName := checkOneFileFoldDiff(fullPath)
        if !exist {
            continue
        }
        if actualFileName != filepath.Base(ie.Path) {
            //在index里面修复大小写错误
            caseDiffChangeArray = append(caseDiffChangeArray, caseDiffChange{
                gitPath:        ie.Path,
                fileSystemPath: filepath.Join(basePath, filepath.Dir(ie.Path), actualFileName),
            })
        }
    }

    if len(caseDiffChangeArray) > 0 {
        fmt.Println("file name diff in case:")
        for _, change := range caseDiffChangeArray {
            fmt.Println("\t", change.gitPath)
            index.RemoveByPath(change.gitPath)
        }
        err = index.Write()
        kmgConsole.ExitOnErr(err)
    }
}

type caseDiffChange struct {
    gitPath        string
    fileSystemPath string
}

var folderFileCache = map[string][]string{}

func checkOneFileFoldDiff(path string) (exist bool, actualFileName string) {
    fileName := filepath.Base(path)
    dirPath := filepath.Dir(path)
    names, ok := folderFileCache[dirPath]
    if !ok {
        dirFile, err := os.Open(filepath.Dir(path))
        defer dirFile.Close()
        if err != nil {
            if os.IsNotExist(err) {
                return false, ""
            }
            kmgConsole.ExitOnErr(err)
        }

        names, err = dirFile.Readdirnames(-1)
        kmgConsole.ExitOnErr(err)
        folderFileCache[dirPath] = names
    }
    CaseDiffFound := false
    actualFileName = ""
    for _, n := range names {
        if n == fileName {
            return true, n
        }
        if strings.EqualFold(n, fileName) {
            CaseDiffFound = true
            actualFileName = n
        }
    }
    if CaseDiffFound {
        return true, actualFileName
    }
    return false, ""
}
*/
