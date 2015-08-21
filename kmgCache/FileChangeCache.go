package kmgCache

import (
	"os"
	"path/filepath"

	"github.com/bronze1man/kmg/encoding/kmgGob"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgCrypto"
	"github.com/bronze1man/kmg/kmgFile"
)

//这个测试要用到.
func getFileChangeCachePath(key string) string {
	// 此处key长度不可控,所以用md5.
	return filepath.Join(kmgConfig.DefaultEnv().TmpPath, "FileChangeCache", kmgCrypto.Md5HexFromString(key))
}

func MustMd5FileChangeCache(key string, pathList []string, f func()) {
	// 此处需要考虑,
	//   用户新添加了一个文件
	//   用户删除了一个文件
	//   用户编辑了一个文件
	//   用户在目录里面添加了一个文件
	//   用在在目录里面删除了一个文件
	//读取文件修改时间缓存信息
	toChange := false
	cacheInfo := map[string]string{}
	cacheFilePath := getFileChangeCachePath(key)
	err := kmgGob.ReadFile(cacheFilePath, &cacheInfo)
	if err != nil {
		//忽略缓存读取的任何错误
		cacheInfo = map[string]string{}
	}
	hasReadFileMap := map[string]bool{}
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
			hasReadFileMap[stat.FullPath] = true
			cacheMd5 := cacheInfo[stat.FullPath]
			if kmgCrypto.MustMd5File(stat.FullPath) != cacheMd5 {
				toChange = true
				//fmt.Printf("[MustMd5FileChangeCache] path:[%s] mod md5 not match save[%s] file[%s]\n", stat.FullPath,
				//	cacheMd5, kmgCrypto.MustMd5File(stat.FullPath))
				break
			}
		}
		if toChange {
			break
		}
	}
	//删除一个已经存储在缓存列表里面的文件,是一个修改.
	for fullPath := range cacheInfo {
		if !hasReadFileMap[fullPath] {
			toChange = true
			break
		}
	}
	if !toChange {
		return
	}
	f()
	cacheInfo = map[string]string{}
	for _, path := range pathList {
		statList, err := kmgFile.GetAllFileAndDirectoryStat(path)
		if err != nil {
			panic(err)
		}
		for _, stat := range statList {
			if stat.Fi.IsDir() {
				continue
			}
			cacheInfo[stat.FullPath] = kmgCrypto.MustMd5File(stat.FullPath)
		}
	}
	kmgFile.MustMkdirForFile(cacheFilePath)
	kmgGob.MustWriteFile(cacheFilePath, cacheInfo)
	//保存文件缓存信息
	return
}
