package kmgQiniu

import (
	"io"

	"github.com/bronze1man/kmg/kmgErr"
	"github.com/qiniu/api/rsf"
)

// 列出所有前缀是xxx的数据,
// 已处理1000个限制
func ListPrefix(ctx *Context, prefix string) (entries []rsf.ListItem, err error) {
	var marker = ""
	for {
		var thisEntries []rsf.ListItem
		thisEntries, marker, err = ctx.rsfClient.ListPrefix(nil, ctx.bucket,
			prefix, marker, 1000)
		entries = append(entries, thisEntries...)
		if err == io.EOF {
			return entries, nil
		}
		if err != nil {
			kmgErr.LogErrorWithStack(err)
			return entries, err
		}
		if len(thisEntries) < 1000 {
			break
		}
	}
	return
}
