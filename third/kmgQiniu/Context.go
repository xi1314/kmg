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
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Context struct {
	client    rs.Client
	rsfClient rsf.Client
	bucket    string //bucket名
	domain    string //下载域名
	isPrivate bool   //是否是私有bucket
}

type Bucket struct {
	Ak string
	Sk string

	Name      string //空间名
	Domain    string //下载使用的域名
	IsPrivate bool   // 是否是私有Api
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

//上传字节 remotePath 开头带 / 或不带 / 效果完全不一样. 正常情况应该是不带 /的
// 此处没有实现流式接口,这个接口的效果和 UploadFromBytes 没有什么差别,(依然会爆内存)
// 分片上传 功能似乎可以解决此类问题,可惜太过复杂了.
func (ctx *Context) UploadFromReader(remotePath string, reader io.Reader) (err error) {
	buf, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	return ctx.UploadFromBytes(remotePath, buf)
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

// 目录开头带 / 或不带 / 效果一致
func (ctx *Context) MustRemoveBatch(PathList []string) {
	ctx.singleContextCheck()
	if len(PathList) == 0 {
		return
	}
	//这个好像也有1000个文件的限制.
	deleteItemList := make([]rs.EntryPath, 0, len(PathList))
	length := len(PathList)
	for i := 0; i < length; i += 1000 {
		end := i + 1000
		if end > length {
			end = length
		}
		deleteItemList = deleteItemList[0:0]
		for j := i; j < end; j++ {
			path := strings.TrimPrefix(PathList[j], "/")
			deleteItemList = append(deleteItemList, rs.EntryPath{
				Key:    path,
				Bucket: ctx.bucket,
			})
		}
		_, err := ctx.client.BatchDelete(nil, deleteItemList)
		if err != nil {
			panic(err)
		}
	}
}

// 返回 scheme和domain ,结尾没有 /
// 例如: http://xxx.com
func (ctx *Context) GetSchemeAndDomain() string {
	return "http://" + ctx.domain
}

type FileInfo struct {
	Path    string //
	Hash    string
	Size    int64
	ModTime time.Time
	//还有几个字段暂时用不着.
}

// 返回的path前面不带 /
func (ctx *Context) MustListPrefix(prefix string) (output []FileInfo) {
	ctx.singleContextCheck()
	prefix = strings.TrimPrefix(prefix, "/")
	entries, err := ListPrefix(ctx, prefix)
	if err != nil {
		panic(err)
	}
	output = make([]FileInfo, len(entries))
	for i := range entries {
		output[i].Path = entries[i].Key
		output[i].Hash = entries[i].Hash
		output[i].Size = entries[i].Fsize
		output[i].ModTime = time.Unix(entries[i].PutTime/1e7, entries[i].PutTime%1e7*100)
	}
	return output
}

func (ctx *Context) singleContextCheck() {
	if ctx != currentContext {
		panic("同时只能有一个Context存在.")
	}
}
