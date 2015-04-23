package kmgCache

import (
	"github.com/bronze1man/kmg/encoding/kmgGob"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgFile"
	"os"
	"path/filepath"
	"time"
)

func getFileChangeCachePath(key string) string {
	return filepath.Join(kmgConfig.DefaultEnv().TmpPath, "FileChangeCache", key)
}

// 根据文件变化,对f这个请求进行缓存
// key表示这件事情的缓存key
// pathList表示需要监控的目录
// 文件列表里面如果有文件不存在,会运行代码
// TODO 有bug还不能用
func MustFileChangeCache(key string, pathList []string, f func()) {
	//读取文件修改时间缓存信息
	toChange := false
	cacheInfo := map[string]time.Time{}
	cacheFilePath := getFileChangeCachePath(key)
	kmgFile.MustMkdirForFile(cacheFilePath)
	err := kmgGob.ReadFile(cacheFilePath, &cacheInfo)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	for _, path := range pathList {
		statList, err := kmgFile.GetAllFileAndDirectoryStat(path)
		if err != nil {
			if os.IsNotExist(err) {
				toChange = true
				break
			}
			panic(err)
		}
		for _, stat := range statList {
			cacheTime := cacheInfo[stat.FullPath]
			if cacheTime.IsZero() {
				toChange = true
				break
			}
			if stat.Fi.ModTime() != cacheTime {
				toChange = true
				break
			}
		}
		if toChange {
			break
		}
	}
	if !toChange {
		return
	}
	f()
	cacheInfo = map[string]time.Time{}
	for _, path := range pathList {
		statList, err := kmgFile.GetAllFileAndDirectoryStat(path)
		if err != nil {
			if os.IsNotExist(err) {
				toChange = true
				break
			}
			panic(err)
		}
		for _, stat := range statList {
			cacheInfo[stat.FullPath] = stat.Fi.ModTime()
			cacheTime := cacheInfo[stat.FullPath]
			if cacheTime.IsZero() {
				toChange = true
				break
			}
			if stat.Fi.ModTime() != cacheTime {
				toChange = true
				break
			}
		}
		if toChange {
			break
		}
	}
	kmgGob.MustWriteFile(cacheFilePath, cacheInfo)
	//保存文件缓存信息
	return
}
