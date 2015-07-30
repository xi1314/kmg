package kmgCache

import (
	"encoding/hex"
	"path/filepath"
	"time"

	"github.com/bronze1man/kmg/encoding/kmgGob"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgFile"
)

type ttlCacheEntryV2 struct {
	Value   []byte
	Timeout time.Time
}

func getFileTtlCachePath(key string) string {
	return filepath.Join(kmgConfig.DefaultEnv().TmpPath, "FileTtlCache", hex.EncodeToString([]byte(key)))
}

// 文件ttl缓存实现,每次都会读取文件,由于没有泛型,此处需要调用者自行解决序列化问题.
// 2.存储在文件里面,以便重启之后可以使用.
// 3.当缓存存在的时候,并且ttl在时间范围内,使用ttl,当缓存不存在的时候,使用回调拉取数据.当缓存过期的时候,使用回调拉取数据.
// 4.每一次使用都会读取文件
// 5.当某一次缓存拉取出现错误的时候,直接返回错误给调用者
func FileTtlCache(key string, f func() (b []byte, ttl time.Duration, err error)) (b []byte, err error) {
	entry := ttlCacheEntryV2{}
	cacheFilePath := getFileTtlCachePath(key)
	now := time.Now()
	err = kmgGob.ReadFile(cacheFilePath, &entry)
	if err == nil && entry.Timeout.After(now) {
		return entry.Value, nil
	}
	b, ttl, err := f()
	if err != nil {
		return nil, err
	}
	entry.Value = b
	entry.Timeout = now.Add(ttl)
	err = kmgFile.MkdirForFile(cacheFilePath)
	if err != nil {
		return nil, err
	}
	err = kmgGob.WriteFile(cacheFilePath, entry)
	if err != nil {
		return nil, err
	}
	return b, nil
}
