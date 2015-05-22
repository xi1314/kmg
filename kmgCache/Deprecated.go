package kmgCache

import (
	"github.com/bronze1man/kmg/encoding/kmgGob"
	"github.com/bronze1man/kmg/kmgFile"
	"os"
	"time"
)

// 这个函数bug太多，以废弃，请使用 MustMd5FileChangeCache
// 根据文件变化,对f这个请求进行缓存
// key表示这件事情的缓存key
// pathList表示需要监控的目录
// 文件列表里面如果有文件不存在,会运行代码
// 已知bug,小于1秒的修改不能被检测到.
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
				//fmt.Printf("[MustFileChangeCache] path:[%s] not exist\n", path)
				break
			}
			panic(err)
		}
		for _, stat := range statList {
			if stat.Fi.IsDir() {
				continue
			}
			cacheTime := cacheInfo[stat.FullPath]
			if cacheTime.IsZero() {
				toChange = true
				//fmt.Printf("[MustFileChangeCache] path:[%s] no save mod time\n", stat.FullPath)
				break
			}
			if stat.Fi.ModTime() != cacheTime {
				toChange = true
				//fmt.Printf("[MustFileChangeCache] path:[%s] mod time not match save[%s] file[%s]\n", stat.FullPath,
				//	cacheTime, stat.Fi.ModTime())
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
			panic(err)
		}
		for _, stat := range statList {
			if stat.Fi.IsDir() {
				continue
			}
			cacheInfo[stat.FullPath] = stat.Fi.ModTime()
		}
	}
	kmgGob.MustWriteFile(cacheFilePath, cacheInfo)
	//保存文件缓存信息
	return
}
