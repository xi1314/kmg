package kmgCache

import (
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
	"time"
)

func TestFileChangeCacheOneFile(t *testing.T) {
	//可以递归 遍历文件
	//缓存数据文件不存在,没有问题
	//指定的文件不存在,也没有问题
	callLog := make([]string, 32)
	//指定的文件不存在
	kmgFile.MustDeleteFile(getFileChangeCachePath("test_file_change_cache"))
	kmgFile.MustDeleteFile("testFile/notExist")

	MustFileChangeCache("test_file_change_cache", []string{
		"testFile/notExist",
	}, func() {
		callLog[1] = "notExist"
		kmgFile.MustWriteFileWithMkdir("testFile/notExist", []byte("1"))
	})
	kmgTest.Equal(callLog[1], "notExist")

	MustFileChangeCache("test_file_change_cache", []string{
		"testFile/notExist",
	}, func() {
		callLog[2] = "notExist"
	})
	kmgTest.Equal(callLog[2], "")

	time.Sleep(time.Second * 1)

	kmgFile.MustWriteFile("testFile/notExist", []byte("2"))
	MustFileChangeCache("test_file_change_cache", []string{
		"testFile/notExist",
	}, func() {
		callLog[3] = "notExist"
	})
	kmgTest.Equal(callLog[3], "notExist")
}
func TestFileChangeCacheOneDir(t *testing.T) {
	callLog := make([]string, 32)
	//递归可用
	kmgFile.MustDeleteFile(getFileChangeCachePath("test_file_change_cache"))
	kmgFile.MustMkdirAll("testFile/d1/d2")
	kmgFile.MustWriteFile("testFile/d1/d2/f3", []byte("1"))
	MustFileChangeCache("test_file_change_cache", []string{
		"testFile/d1",
	}, func() {
		callLog[3] = "f3"
	})
	kmgTest.Equal(callLog[3], "f3")

	MustFileChangeCache("test_file_change_cache", []string{
		"testFile/d1",
	}, func() {
		callLog[4] = "f3"
	})
	kmgTest.Equal(callLog[4], "")

	time.Sleep(time.Second * 1)
	kmgFile.MustWriteFile("testFile/d1/d2/f3", []byte("2"))
	MustFileChangeCache("test_file_change_cache", []string{
		"testFile/d1",
	}, func() {
		callLog[5] = "f3"
	})
	kmgTest.Equal(callLog[5], "f3")
}
