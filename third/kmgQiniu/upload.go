package kmgQiniu

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bronze1man/kmg/kmgTask"
	qiniuIo "github.com/qiniu/api/io"
	"github.com/qiniu/api/rs"
	"github.com/qiniu/rpc"
)

//上传文件或目录
var Upload = UploadDirMulitThread

//单线程上传目录
//检查同文件名和同内容的文件是否存在,如果内容相同便不上传
func UploadDir(ctx *Context, localRoot string, remoteRoot string) (err error) {
	return filepath.Walk(localRoot, func(path string, info os.FileInfo, inErr error) (err error) {
		if inErr != nil {
			return inErr
		}
		if info.IsDir() {
			return
		}
		relPath, err := filepath.Rel(localRoot, path)
		if err != nil {
			return fmt.Errorf("[qiniuUploadDir] filepath.Rel err:%s", err.Error())
		}
		remotePath := filepath.Join(remoteRoot, relPath)
		err = UploadFileCheckExist(ctx, path, remotePath)
		if err != nil {
			return fmt.Errorf("[qiniuUploadDir] qiniuUploadFile err:%s", err.Error())
		}
		return
	})
}

type uploadFileRequest struct {
	localPath  string
	remotePath string
	expectHash string
}

//多线程上传目录
//1.某个文件仅在一个线程中上传,
//2.检查同名和同内容的文件是否已经存在了,如果存在,且hash相同便不上传(断点续传)
// TODO 解决一边上传,一边修改的bug.
func UploadDirMulitThread(ctx *Context, localRoot string, remoteRoot string) (err error) {
	tm := kmgTask.NewLimitThreadErrorHandleTaskManager(ThreadNum, 3)
	defer tm.Close()
	requestList := []uploadFileRequest{}
	remotePathList :=[]string{}
	//dispatch task 分配任务
	err = filepath.Walk(localRoot, func(path string, info os.FileInfo, inErr error) (err error) {
		if inErr != nil {
			return inErr
		}
		if info.IsDir() {
			return
		}
		relPath, err := filepath.Rel(localRoot, path)
		if err != nil {
			return fmt.Errorf("[qiniuUploadDir] filepath.Rel err:%s", err.Error())
		}
		remotePath := NormalizeRemotePath(filepath.Join(remoteRoot, relPath))
		expectHash, err := ComputeHashFromFile(path)
		if err != nil {
			return
		}
		requestList = append(requestList, uploadFileRequest{
			localPath:  path,
			remotePath: remotePath,
			expectHash: expectHash,
		})
		remotePathList = append(remotePathList,remotePath)
		return
	})
	if len(requestList) == 0 {
		return ErrNoFile
	}
	batchRet,err:=ctx.BatchStat(remotePathList)
	if err!=nil{
		return
	}
	for i,ret:=range batchRet{
		i:=i
		if ret.IsExist() && ret.Hash == requestList[i].expectHash {
			continue
		}
		tm.AddTask(func() (err error) {
			return UploadFileWithHash(ctx, requestList[i].localPath, requestList[i].remotePath, requestList[i].expectHash)
		})
	}
	////群发状态询问消息减少网络连接数量,加快速度
	//entryPathList := make([]rs.EntryPath, len(requestList))
	//for i, req := range requestList {
	//	entryPathList[i].Bucket = ctx.bucket
	//	entryPathList[i].Key = req.remotePath
	//}
	//batchRet, err := ctx.client.BatchStat(nil, entryPathList)
	///*
	//   //此处返回的错误很奇怪,有大量文件不存在信息,应该是正常情况,此处最简单的解决方案就是假设没有错误
	//   if err != nil{
	//       fmt.Printf("%T %#v\n",err,err)
	//       err1,ok:=err.(*rpc.ErrorInfo)
	//       if !ok{
	//           return err
	//       }
	//       if err1.Code!=298{
	//           return err
	//       }
	//   }
	//*/
	//if len(batchRet) != len(entryPathList) {
	//	return fmt.Errorf("[UploadDirMulitThread] len(batchRet)[%d]!=len(entryPathList)[%d] err[%s]",
	//		len(batchRet), len(entryPathList),err)
	//}
	//for i, ret := range batchRet {
	//	i := i
	//	//验证hash,当文件不存在时,err是空
	//	if ret.Error != "" && ret.Error != "no such file or directory" {
	//		return fmt.Errorf("[UploadDirMulitThread] [remotePath:%s]ctx.client.BatchStat err[%s]",
	//			requestList[i].remotePath, ret.Error)
	//	}
	//	if ret.Data.Hash == requestList[i].expectHash {
	//		continue
	//	}
	//	tm.AddTask(func() (err error) {
	//		return UploadFileWithHash(ctx, requestList[i].localPath, requestList[i].remotePath, requestList[i].expectHash)
	//	})
	//}
	tm.Wait()
	if err != nil {
		return err
	}
	err = tm.GetError()
	return
}

//上传文件,检查同名和同内容文件
//先找cdn上是不是已经有一样的文件了,以便分文件断点续传,再上传
func UploadFileCheckExist(ctx *Context, localPath string, remotePath string) (err error) {
	remotePath = NormalizeRemotePath(remotePath)
	entry, err := ctx.client.Stat(nil, ctx.bucket, remotePath)
	if err != nil {
		if !(err.(*rpc.ErrorInfo) != nil && err.(*rpc.ErrorInfo).Err == "no such file or directory") {
			return err
		}
	}
	expectHash, err := ComputeHashFromFile(localPath)
	if err != nil {
		return
	}
	//already have a file with same context and same key,do nothing
	if entry.Hash == expectHash {
		return
	}
	return UploadFileWithHash(ctx, localPath, remotePath, expectHash)
}

//上传文件,检查返回的hash和需要的hash是否一致
func UploadFileWithHash(ctx *Context, localPath string, remotePath string, expectHash string) (err error) {
	var ret qiniuIo.PutRet
	var extra = &qiniuIo.PutExtra{
		CheckCrc: 1,
	}
	putPolicy := rs.PutPolicy{
		Scope: ctx.bucket + ":" + remotePath,
	}
	uptoken := putPolicy.Token(nil)
	err = qiniuIo.PutFile(nil, &ret, uptoken, remotePath, localPath, extra)
	//fmt.Println(localPath,remotePath,err)
	if err != nil {
		return
	}
	if ret.Hash != expectHash {
		return fmt.Errorf("[UploadFileWithHash][remotePath:%s] ret.Hash:[%s]!=expectHash[%s] ", remotePath, ret.Hash, expectHash)
	}

	return
}
