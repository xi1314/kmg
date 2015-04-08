package kmgQiniu

import (
	"github.com/bronze1man/kmg/kmgErr"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/qiniu/api/rs"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

var Download = DownloadDir

//下载单个文件,到本地,会覆盖本地已经存在的文件,会创建所有父级目录,会使用hash检查文件是否存在.
func DownloadFileWithHash(ctx *Context, remotePath string, localPath string, hash string) (err error) {
	fhash, err := ComputeHashFromFile(localPath)
	if err == nil && fhash == hash {
		return
	}
	if err != nil && !os.IsNotExist(err) {
		kmgErr.LogErrorWithStack(err)
		return
	}
	if fhash == hash {
		return
	}
	return DownloadFile(ctx, remotePath, localPath)
}

func DownloadFile(ctx *Context, remotePath string, localPath string) (err error) {
	var downloadUrl string
	if ctx.isPrivate {
		baseUrl := rs.MakeBaseUrl(ctx.domain, remotePath)
		policy := rs.GetPolicy{}
		downloadUrl = policy.MakeRequest(baseUrl, nil)
	} else {
		downloadUrl = rs.MakeBaseUrl(ctx.domain, remotePath)
	}
	resp, err := http.Get(downloadUrl)
	if err != nil {
		kmgErr.LogErrorWithStack(err)
		return
	}
	defer resp.Body.Close()
	err = kmgFile.MkdirForFile(localPath)
	if err != nil {
		kmgErr.LogErrorWithStack(err)
		return
	}
	f, err := os.Create(localPath)
	if err != nil {
		kmgErr.LogErrorWithStack(err)
		return
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		kmgErr.LogErrorWithStack(err)
		return
	}
	return
}

func DownloadDir(ctx *Context, remoteRoot string, localRoot string) (err error) {
	entries, err := ListPrefix(ctx, remoteRoot)
	if err != nil {
		kmgErr.LogErrorWithStack(err)
		return err
	}
	if len(entries) == 0 {
		return ErrNoFile
	}
	for _, entry := range entries {
		refPath, err := filepath.Rel(remoteRoot, entry.Key)
		if err != nil {
			kmgErr.LogErrorWithStack(err)
			return err
		}
		err = DownloadFileWithHash(ctx, entry.Key, filepath.Join(localRoot, refPath), entry.Hash)
		if err != nil {
			return err
		}
	}
	return nil
}
