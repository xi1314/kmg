package kmgFile

import (
	"os"
	"path/filepath"
)

/*
func AllDirectory(root string)(out []string,err error){
    err=filepath.Walk(root,func(path string, info os.FileInfo, err error) error {

    })
}
*/

type StatAndFullPath struct {
	Fi       os.FileInfo
	FullPath string
}

// 获取这个路径的所有文件的状态和完整路径
//   如果输入是一个文件,则返回这个文件的完整路径
//   如果输入是一个目录,则返回这个目录和下面所有目录和文件的信息和完整路径
//   目前暂不明确symlink的文件会如何处理
func GetAllFileAndDirectoryStat(root string) (out []StatAndFullPath, err error) {
	root, err = Realpath(root)
	if err != nil {
		return nil, err
	}
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		out = append(out, StatAndFullPath{
			FullPath: path,
			Fi:       info,
		})
		return nil
	})
	return
}

//返回这个目录下面所有的文件,返回格式为完整文件名
func GetAllFiles(root string) (out []string, err error) {
	root, err = Realpath(root)
	if err != nil {
		return nil, err
	}
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			out = append(out, path)
		}
		return nil
	})
	return
}

// 获得所有文件,不包括目录
func MustGetAllFiles(root string) (out []string) {
	out, err := GetAllFiles(root)
	if err != nil {
		panic(err)
	}
	return out
}

/*
// 回调只会给文件名,回调返回true表示需要这个文件,返回false表示不需要这个文件
func MustGetAllFilesWithCallback(root string,cb func(absPath string)bool) (out []string) {
	root, err := Realpath(root)
	if err != nil {
		panic(err)
	}
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir(){
			return nil
		}
		if cb(path) {
			out = append(out,path)
		}
		return nil
	})
	if err!=nil{
		panic(err)
	}
	return
}
*/
