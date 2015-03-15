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
