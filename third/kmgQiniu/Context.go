package kmgQiniu

import (
	"github.com/qiniu/api/conf"
	qiniuIo "github.com/qiniu/api/io"
	"github.com/qiniu/api/rs"
	"github.com/qiniu/api/rsf"

	"fmt"
	"github.com/qiniu/bytes"
	"hash/crc32"
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

//注意: 由于实现的问题,此处Ak和Sk只能使用同一个
func NewContext(bucket Bucket) *Context {
	conf.ACCESS_KEY = bucket.Ak
	conf.SECRET_KEY = bucket.Sk
	return &Context{
		client:    rs.New(nil),
		rsfClient: rsf.New(nil),
		bucket:    bucket.Name,
		domain:    bucket.Domain,
		isPrivate: bucket.IsPrivate,
	}
}

//可以下载文件或目录 remoteRoot 开头不要带 /
func (ctx *Context) DownloadToFile(remoteRoot string, localRoot string) (err error) {
	return DownloadDir(ctx, remoteRoot, localRoot)
}

func (ctx *Context) MustDownloadToFile(remoteRoot string, localRoot string) {
	err := DownloadDir(ctx, remoteRoot, localRoot)
	if err != nil {
		panic(err)
	}
	return
}

//可以上传文件或目录 remoteRoot 开头带 / 或不带 / 效果一致
func (ctx *Context) UploadFromFile(localRoot string, remoteRoot string) (err error) {
	return UploadDirMulitThread(ctx, localRoot, remoteRoot)
}

func (ctx *Context) MustUploadFromFile(localRoot string, remoteRoot string) {
	err := UploadDirMulitThread(ctx, localRoot, remoteRoot)
	if err != nil {
		panic(err)
	}
	return
}

//上传字节 remotePath 开头带 / 或不带 / 效果一致
func (ctx *Context) UploadFromBytes(remotePath string, b []byte) (err error) {
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
	err := ctx.UploadFromBytes(remotePath, context)
	if err != nil {
		panic(err)
	}
	return
}

//prefix 开头带 / 或不带 / 效果一致
func (ctx *Context) RemovePrefix(prefix string) (err error) {
	return RemovePrefix(ctx, prefix)
}
