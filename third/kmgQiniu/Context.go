package kmgQiniu

import (
	"github.com/qiniu/api/conf"
	qiniuIo "github.com/qiniu/api/io"
	"github.com/qiniu/api/rs"
	"github.com/qiniu/api/rsf"

	"bytes"
	"fmt"
	"hash/crc32"
	"io"
	"net/http"
	"strings"
)

type Context struct {
	client    rs.Client
	rsfClient rsf.Client
	bucket    string //bucket名
	domain    string //下载域名
	isPrivate bool   //是否是私有bucket
}

type Bucket struct {
	Name      string
	Domain    string
	IsPrivate bool
	Ak        string
	Sk        string
}

var currentContext *Context

//注意: 由于实现的问题,全局只能使用一个Context,
// TODO 解决全局只能使用一个Context的问题
func NewContext(bucket Bucket) *Context {
	conf.ACCESS_KEY = bucket.Ak
	conf.SECRET_KEY = bucket.Sk
	currentContext = &Context{
		client:    rs.New(nil),
		rsfClient: rsf.New(nil),
		bucket:    bucket.Name,
		domain:    bucket.Domain,
		isPrivate: bucket.IsPrivate,
	}
	return currentContext
}

//可以下载文件或目录 remoteRoot 开头带 / 或不带 / 效果一致
func (ctx *Context) DownloadToFile(remoteRoot string, localRoot string) (err error) {
	ctx.singleContextCheck()
	remoteRoot = strings.TrimPrefix(remoteRoot, "/")
	return DownloadDir(ctx, remoteRoot, localRoot)
}

func (ctx *Context) MustDownloadToFile(remoteRoot string, localRoot string) {
	ctx.singleContextCheck()
	remoteRoot = strings.TrimPrefix(remoteRoot, "/")
	err := DownloadDir(ctx, remoteRoot, localRoot)
	if err != nil {
		panic(err)
	}
	return
}

//下载到一个Writer里面
func (ctx *Context) DownloadToWriter(remotePath string, w io.Writer) (err error) {
	remotePath = strings.TrimPrefix(remotePath, "/")
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
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return err
	}
	return
}

func (ctx *Context) MustDownloadToBytes(remotePath string) (b []byte) {
	buf := &bytes.Buffer{}
	err := ctx.DownloadToWriter(remotePath, buf)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

//可以上传文件或目录 remoteRoot 开头带 / 或不带 / 效果一致
func (ctx *Context) UploadFromFile(localRoot string, remoteRoot string) (err error) {
	ctx.singleContextCheck()
	remoteRoot = strings.TrimPrefix(remoteRoot, "/")
	return UploadDirMulitThread(ctx, localRoot, remoteRoot)
}

func (ctx *Context) MustUploadFromFile(localRoot string, remoteRoot string) {
	ctx.singleContextCheck()
	remoteRoot = strings.TrimPrefix(remoteRoot, "/")
	err := UploadDirMulitThread(ctx, localRoot, remoteRoot)
	if err != nil {
		panic(err)
	}
	return
}

//上传字节 remotePath 开头带 / 或不带 / 效果完全不一样. 正常情况应该是不带 /的
func (ctx *Context) UploadFromBytes(remotePath string, b []byte) (err error) {
	ctx.singleContextCheck()
	remotePath = strings.TrimPrefix(remotePath, "/")
	h := crc32.NewIEEE()
	h.Write(b)
	crc := h.Sum32()
	var ret qiniuIo.PutRet
	var extra = &qiniuIo.PutExtra{
		Crc32:    crc,
		CheckCrc: 2,
	}
	putPolicy := rs.PutPolicy{
		Scope: ctx.bucket + ":" + remotePath,
	}
	uptoken := putPolicy.Token(nil)
	r := bytes.NewReader(b)
	err = qiniuIo.Put2(nil, &ret, uptoken, remotePath, r, int64(len(b)), extra)
	if err != nil {
		return
	}
	expectHash := ComputeHashFromBytes(b)
	if ret.Hash != expectHash {
		return fmt.Errorf("[UploadFileWithHash][remotePath:%s] ret.Hash:[%s]!=expectHash[%s] ", remotePath, ret.Hash, expectHash)
	}
	return
}

func (ctx *Context) MustUploadFromBytes(remotePath string, context []byte) {
	ctx.singleContextCheck()
	err := ctx.UploadFromBytes(remotePath, context)
	if err != nil {
		panic(err)
	}
	return
}

//prefix 开头带 / 或不带 / 效果一致
func (ctx *Context) RemovePrefix(prefix string) (err error) {
	ctx.singleContextCheck()
	prefix = strings.TrimPrefix(prefix, "/")
	return RemovePrefix(ctx, prefix)
}

func (ctx *Context) singleContextCheck() {
	if ctx != currentContext {
		panic("同时只能有一个Context存在.")
	}
}
