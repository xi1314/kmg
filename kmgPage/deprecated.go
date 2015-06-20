package kmgPage

import (
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	"github.com/bronze1man/kmg/kmgSql/MysqlAst"
)

// @deprecated
type CreateFromSelectCommandRequest struct {
	Select      *MysqlAst.SelectCommand //必填
	Ctx         kmgHttp.Context         //可选
	Url         string                  //可选
	ItemPerPage int                     //可选
	CurrentPage int                     //可选
	PageKeyName string                  //可选
}

// @deprecated
func CreateFromSelectCommand(req CreateFromSelectCommandRequest) *KmgPage {
	page := &KmgPage{}
	page.BaseUrl = req.Url
	page.ItemPerPage = req.ItemPerPage
	page.CurrentPage = req.CurrentPage
	page.PageKeyName = req.PageKeyName
	page.init()
	return page.runSelectCommand(req.Select)
}
