package goCmd

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"runtime"
	"sync"

	"github.com/bronze1man/kmg/encoding/kmgGob"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgConsole"
	"github.com/bronze1man/kmg/kmgCrypto"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgTask"
	"go/format"
)

func runGoFmt() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	wd, err := os.Getwd()
	kmgConsole.ExitOnErr(err)
	err = GoFmtDir(wd)
	kmgConsole.ExitOnErr(err)
}

func isGoFile(f os.FileInfo) bool {
	// ignore non-Go files
	name := f.Name()
	return !f.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go")
}

func GoFmtDir(path string) (outErr error) {
	dir := &goFmtDir{
		cacheMap:    map[string]bool{},
		newCacheMap: map[string]bool{},
		tasker:      kmgTask.NewLimitThreadTaskManager(runtime.NumCPU()),
	}
	if kmgConfig.IsWdHaveEnv() {
		kmgGob.ReadFile(kmgConfig.DefaultEnv().PathInTmp("gofmt.gob"), &dir.cacheMap)
		//此处故意忽略错误,如果失败cacheMap应该没有变化,表示没有缓存而已.
		dir.cacheAble = true
	}
	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			dir.addErr(err)
			return err
		}
		if !isGoFile(f) {
			return nil
		}
		dir.tasker.AddFunc(func() {
			err = dir.processFile(path)
			if err != nil {
				dir.addErr(err)
			}
		})
		return nil
	})
	dir.tasker.Close()
	err := kmgGob.WriteFile(kmgConfig.DefaultEnv().PathInTmp("gofmt.gob"), dir.newCacheMap)
	if err != nil {
		return err
	}
	if len(dir.outErrList) == 0 {
		return nil
	}
	s := ""
	for _, er := range dir.outErrList {
		s += er.Error() + "\n"
	}
	return errors.New(s)
}

type goFmtDir struct {
	cacheAble   bool
	cacheMap    map[string]bool
	newCacheMap map[string]bool
	outErrList  []error
	locker      sync.Mutex
	tasker      *kmgTask.LimitThreadTaskManager
}

func (dir *goFmtDir) addErr(err error) {
	dir.locker.Lock()
	dir.outErrList = append(dir.outErrList, err)
	dir.locker.Unlock()
}
func (dir *goFmtDir) processFile(filename string) (err error) {
	src, err := kmgFile.ReadFile(filename)
	if err != nil {
		return err
	}
	srcMd5 := kmgCrypto.Md5Hex(src)
	if dir.cacheMap[srcMd5] {
		dir.locker.Lock()
		dir.newCacheMap[srcMd5] = true
		dir.locker.Unlock()
		return nil
	}
	res, err := format.Source(src)
	if err != nil {
		return err
	}
	resMd5 := kmgCrypto.Md5Hex(res)
	dir.locker.Lock()
	dir.newCacheMap[resMd5] = true
	dir.locker.Unlock()
	if srcMd5 != resMd5 {
		// formatting has changed
		err = ioutil.WriteFile(filename, res, 0)
		if err != nil {
			return err
		}
	}

	return err
}

var (
// main operation modes
//options = &imports.Options{
//	TabWidth:  8,
//	TabIndent: true,
//	Comments:  true,
//	Fragment:  false,
//}
)

func processFile(filename string) error {
	//opt := options
	src, err := kmgFile.ReadFile(filename)
	if err != nil {
		return err
	}
	res, err := format.Source(src)
	if err != nil {
		return err
	}
	//res, err := imports.Process(filename, src, opt)
	//if err != nil {
	//	return err
	//}

	if !bytes.Equal(src, res) {
		// formatting has changed
		err = ioutil.WriteFile(filename, res, 0)
		if err != nil {
			return err
		}
	}

	return err
}
