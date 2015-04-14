package kmgQiniu

import (
	"github.com/qiniu/api/conf"
	"github.com/qiniu/api/rs"
	"github.com/qiniu/api/rsf"
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
