package kmgFile

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgRand"
)

func IsDotFile(path string) bool {
	if path == "./" {
		return false
	}
	base := filepath.Base(path)
	if strings.HasPrefix(base, ".") {
		return true
	}
	return false
}

func GetFileBaseWithoutExt(p string) string {
	return filepath.Base(p[:len(p)-len(filepath.Ext(p))])
}

func WriteFile(path string, content []byte) (err error) {
	return ioutil.WriteFile(path, content, os.FileMode(0777))
}
func MustWriteFile(path string, content []byte) {
	err := ioutil.WriteFile(path, content, os.FileMode(0777))
	if err != nil {
		panic(err)
	}
}

func MustWriteFileWithMkdir(path string, content []byte) {
	MustMkdirForFile(path)
	MustWriteFile(path, content)
}

func ReadFile(path string) (content []byte, err error) {
	return ioutil.ReadFile(path)
}

func MustReadFile(path string) (content []byte) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return content
}

//如果这个目录已经创建过了,不报错
func Mkdir(path string) (err error) {
	return os.MkdirAll(path, os.FileMode(0777))
}

//如果这个目录已经创建过了,不报错
func MustMkdir(dirname string) {
	err := os.MkdirAll(dirname, os.FileMode(0777))
	if err != nil {
		panic(err)
	}
}

//保证一个文件的路径可以写入,如果这个目录已经创建过了,不报错
func MkdirForFile(path string) (err error) {
	path = filepath.Dir(path)
	return os.MkdirAll(path, os.FileMode(0777))
}

func MustMkdirForFile(path string) {
	path = filepath.Dir(path)
	err := os.MkdirAll(path, os.FileMode(0777))
	if err != nil {
		panic(err)
	}
}

func AppendFile(path string, content []byte) (err error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0777))
	if err != nil {
		return
	}
	defer f.Close()
	_, err = f.Write(content)
	return
}

func MustAppendFile(path string, content []byte) {
	err := AppendFile(path, content)
	if err != nil {
		panic(err)
	}
}

func FileExist(path string) (exist bool, err error) {
	_, err = os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, err
}

func MustFileExist(path string) bool {
	exist, err := FileExist(path)
	if err != nil {
		panic(err)
	}
	return exist
}

//from http://stackoverflow.com/a/13027975/1586797
func RemoveExtFromFilePath(path string) string {
	return path[:len(path)-len(filepath.Ext(path))]
}

//just some Knowledge,you can direct call ioutil.ReadDir
func ReadDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirname)
}

//delete file or directory,ignore file not exist err
func MustDelete(path string) {
	err := os.RemoveAll(path)
	if os.IsNotExist(err) {
		return
	}
	if err != nil {
		panic(err)
	}
	return
}

// copy file
// can not copy directory
// * override dst file if it exist,
// * mkdir if base dir not exist
//from http://stackoverflow.com/a/21067803/1586797
func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(filepath.Dir(dst), os.FileMode(0777))
			if err != nil {
				return err
			}
			out, err = os.Create(dst)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("[CopyFile] createDst err[%s]", err.Error())
		}
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	//why this?
	//err = out.Sync()
	return
}

func MustCopyFile(src, dst string) {
	err := CopyFile(src, dst)
	if err != nil {
		panic(err)
	}
}

//拷贝文件,把文件从src拷贝到dst
// 如果源文件不存在,不报错
func MustCopyFileIgnoreNotExist(src, dst string) {
	err := CopyFile(src, dst)
	if os.IsNotExist(err) {
		return
	}
	if err != nil {
		panic(err)
	}
}

func MustChangeToTmpPath() string {
	folder := "/tmp/kmg/" + kmgRand.MustCryptoRandToHex(6)
	MustMkdirAll(folder)
	err := os.Chdir(folder)
	if err != nil {
		panic(err)
	}
	return folder
}

func MustRename(oldpath string, newpath string) {
	err := os.Rename(oldpath, newpath)
	if err != nil {
		panic(err)
	}
}

func MustSymlink(fromPath string, toPath string) {
	kmgCmd.CmdSlice([]string{"ln", "-sf", fromPath, toPath}).MustStdioRun()
}
