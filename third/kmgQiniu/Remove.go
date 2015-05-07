package kmgQiniu

import "github.com/qiniu/api/rs"

//使用前缀删除文件,没有文件删,不报错
func RemovePrefix(ctx *Context, prefix string) (err error) {
	itemList, err := ListPrefix(ctx, prefix)
	if err != nil {
		return
	}
	if len(itemList) == 0 {
		return nil
	}
	deleteItemList := make([]rs.EntryPath, len(itemList))
	for i, item := range itemList {
		deleteItemList[i] = rs.EntryPath{
			Key:    item.Key,
			Bucket: ctx.bucket,
		}
	}
	_, err = ctx.client.BatchDelete(nil, deleteItemList)
	if err != nil {
		return err
	}
	return
}
