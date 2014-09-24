package upyun

import (
	"github.com/bronze1man/kmg/kmgLog"
	"github.com/bronze1man/kmg/kmgTask"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

//批量upyun操作
type BulkUpyun struct {
	UpYun *UpYun
	Tm    *kmgTask.LimitThreadTaskManager
}

//批量上传接口
//upload a file
func (obj *BulkUpyun) UploadFile(upyun_path, local_path string) {
	obj.Tm.AddTask(kmgTask.TaskFunc(func() {
		kmgLog.Log("upyun", "upload file: "+upyun_path, nil)
		file, err := os.Open(local_path)
		if err != nil {
			kmgLog.Log("upyunError", err.Error(), err)
			return
		}
		defer file.Close()
		err = obj.UpYun.WriteFile(upyun_path, file, true)
		if err != nil {
			kmgLog.Log("upyunError", err.Error(), err)
			return
		}
		return
	}))
}

//upload a dir
func (obj *BulkUpyun) UploadDir(upyun_path, local_path string) {
	obj.Tm.AddTask(kmgTask.TaskFunc(func() {
		kmgLog.Log("upyun", "upload dir: "+upyun_path, nil)

		dir, err := os.Open(local_path)
		if err != nil {
			kmgLog.Log("upyunError", err.Error(), err)
			return
		}
		file_list, err := dir.Readdir(0)
		if err != nil {
			kmgLog.Log("upyunError", err.Error(), err)
			return
		}
		err = obj.UpYun.MkDir(upyun_path, true)
		if err != nil {
			kmgLog.Log("upyunError", err.Error(), err)
			return
		}
		for _, file_info := range file_list {
			file_name := file_info.Name()
			this_local_path := local_path + "/" + file_name
			this_upyun_path := upyun_path + "/" + file_name
			if file_info.IsDir() {
				obj.UploadDir(this_upyun_path, this_local_path)
			} else {
				obj.UploadFile(this_upyun_path, this_local_path)
			}
		}
		return
	}))
}

//download a file
func (obj *BulkUpyun) DownloadFile(upyun_path, local_path string) {
	obj.Tm.AddTask(kmgTask.TaskFunc(func() {
		kmgLog.Log("upyun", "download file: "+upyun_path, nil)
		err := os.MkdirAll(filepath.Dir(local_path), os.FileMode(0777))
		if err != nil {
			kmgLog.Log("upyunError", err.Error(), err)
			return
		}

		file, err := os.Create(local_path)
		if err != nil {
			kmgLog.Log("upyunError", err.Error(), err)
			return
		}
		defer file.Close()
		err = obj.UpYun.ReadFile(upyun_path, file)
		if err != nil {
			kmgLog.Log("upyunError", err.Error(), err)
			return
		}
		return
	}))
}

//resursive download a dir
func (obj *BulkUpyun) DownloadDir(upyun_path string, file_path string) {
	obj.Tm.AddTask(kmgTask.TaskFunc(func() {
		kmgLog.Log("upyun", "download dir: "+upyun_path, nil)
		file_list, err := obj.UpYun.ReadDir(upyun_path)
		file_mode := os.FileMode(0777)
		if err != nil {
			kmgLog.Log("upyunError", err.Error(), err)
			return
		}
		for _, file_info := range file_list {
			file_type := file_info.Type
			file_name := file_info.Name
			this_local_path := file_path + "/" + file_name
			this_upyun_path := upyun_path + "/" + file_name
			if file_type == "folder" {
				err := os.MkdirAll(this_local_path, file_mode)
				if err != nil {
					kmgLog.Log("upyunError", "os.MkdirAll fail!"+err.Error(), err)
					return
				}
				obj.DownloadDir(this_upyun_path, this_local_path)
			} else if file_type == "file" {
				obj.DownloadFile(this_upyun_path, this_local_path)
			} else {
				kmgLog.Log("upyunError", "unknow file type2:"+file_type, err)
				return
			}
		}
		return
	}))
}

//delete a file
func (obj *BulkUpyun) DeleteFile(upyun_path string) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	obj.deleteFile(upyun_path, wg)
	obj.Tm.AddTaskNewThread(kmgTask.TaskFunc(func() {
		wg.Wait()
	}))
}

//resursive delete a dir
func (obj *BulkUpyun) DeleteDir(upyun_path string) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	obj.deleteDir(upyun_path, wg)
	obj.Tm.AddTaskNewThread(kmgTask.TaskFunc(func() {
		wg.Wait()
	}))
}

//delete a file
func (obj *BulkUpyun) deleteFile(upyun_path string, finish_wg *sync.WaitGroup) {
	obj.Tm.AddTask(kmgTask.TaskFunc(func() {
		defer finish_wg.Done()
		kmgLog.Log("upyun", "delete file: "+upyun_path, nil)
		err := obj.UpYun.DeleteFile(upyun_path)
		if err != nil {
			kmgLog.Log("upyunError", "delete file failed!:"+upyun_path+":"+err.Error(), nil)
			return
		}
		return
	}))
}

//we need to know when is finish delete all file in it ,so we can delete the dir
func (obj *BulkUpyun) deleteDir(upyun_path string, finish_wg *sync.WaitGroup) {
	obj.Tm.AddTask(kmgTask.TaskFunc(func() {
		wg := &sync.WaitGroup{}
		defer obj.Tm.AddTaskNewThread(kmgTask.TaskFunc(func() {
			wg.Wait()
			wg.Add(1)
			obj.deleteFile(upyun_path, wg)
			wg.Wait()
			finish_wg.Done()
		}))
		kmgLog.Log("upyun", "delete dir: "+upyun_path, nil)
		file_list, err := obj.UpYun.ReadDir(upyun_path)
		if err != nil {
			kmgLog.Log("upyunError", err.Error(), err)
			return
		}
		for _, file_info := range file_list {
			file_type := file_info.Type
			file_name := file_info.Name
			this_upyun_path := upyun_path + "/" + file_name
			if file_type == "folder" {
				wg.Add(1)
				obj.deleteDir(this_upyun_path, wg)
			} else if file_type == "file" {
				wg.Add(1)
				obj.deleteFile(this_upyun_path, wg)
			} else {
				kmgLog.Log("upyunError", "unknow file type2:"+file_type, nil)
				return
			}
		}
		return
	}))
}

func (obj *BulkUpyun) GetFileType(upyun_path string) (file_type string, err error) {
	info, err := obj.UpYun.GetFileInfo(upyun_path)
	if err != nil {
		return
	}
	file_type = info["type"]
	return
}

func (obj *BulkUpyun) GetFileSize(upyun_path string) (size uint64, err error) {
	info, err := obj.UpYun.GetFileInfo(upyun_path)
	if err != nil {
		return
	}
	size, err = strconv.ParseUint(info["size"], 10, 64)
	if err != nil {
		return
	}
	return
}

func (obj *BulkUpyun) GetFileDate(upyun_path string) (date time.Time, err error) {
	info, err := obj.UpYun.GetFileInfo(upyun_path)
	if err != nil {
		return
	}
	date, err = time.Parse(time.RFC1123, info["date"])
	if err != nil {
		return
	}
	return
}
