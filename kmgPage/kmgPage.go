package kmgPage

import (
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	"github.com/bronze1man/kmg/kmgSql"
	"github.com/bronze1man/kmg/kmgSql/MysqlAst"
	"math"
	"strconv"
)

type KmgPage struct {
	ItemPerPage int                 // 每页数据量
	ShowPageNum int                 // 显示分页的页面的时候，显示出来的页数 如 10，11，12，13，14，15 就是显示了6个页
	PageKeyName string              // 页面参数的名字
	CurrentPage int                 // 当前页数，默认为1
	TotalItem   int                 // 总数据量
	BaseUrl     string              // 不包含页码参数的url,具体渲染url的时候会把页面参数加上
	Data        []map[string]string // 本次分页查询到的数据
	TotalPage   int                 // 共有多少个页面
}

type CreateFromSelectCommandRequest struct {
	Select      *MysqlAst.SelectCommand
	Url         string
	ItemPerPage int
	CurrentPage int
}

func CreateFromSelectCommand(req CreateFromSelectCommandRequest) *KmgPage {
	page := &KmgPage{}
	page.BaseUrl = req.Url
	page.ItemPerPage = req.ItemPerPage
	page.CurrentPage = req.CurrentPage
	page.init()
	return page.runSelectCommand(req.Select)
}

func (page *KmgPage) CreateFromData(data []map[string]string, baseUrl string, itemPerPage int) *KmgPage {
	page.init()
	page.BaseUrl = baseUrl
	page.ItemPerPage = itemPerPage
	page.TotalItem = len(data)
	return page
}

func (page *KmgPage) runSelectCommand(selectCommand *MysqlAst.SelectCommand) *KmgPage {
	if page.BaseUrl == "" {
		panic("runSelectCommand need baseUrl parameter")
	}
	output, parameterList := selectCommand.GetPrepareParameter()
	countData, err := kmgSql.QueryOne("SELECT COUNT(*) AS c FROM ("+output+") as View", parameterList...)
	if err != nil {
		panic(err)
	}
	page.TotalItem, err = strconv.Atoi(countData["c"])
	if err != nil {
		panic(err)
	}
	dataSelect := selectCommand.Copy()
	dataSelect = dataSelect.Limit(strconv.Itoa(page.GetMysqlOffset()) + "," + strconv.Itoa(page.ItemPerPage))
	page.Data = kmgSql.MustRunSelectCommand(dataSelect)
	return page
}

func (page *KmgPage) init() {
	if page.ItemPerPage == 0 {
		page.ItemPerPage = 10
	}
	page.ShowPageNum = 10
	page.PageKeyName = "page"
	if page.CurrentPage == 0 {
		page.CurrentPage = 1
	}
	if page.CurrentPage < 1 {
		page.CurrentPage = 1
	}
	if page.CurrentPage > page.GetTotalPage() {
		page.CurrentPage = page.GetTotalPage()
	}
}

func (page *KmgPage) GetTotalPage() int {
	totalPage := int(math.Ceil(float64(page.TotalItem) / float64(page.ItemPerPage)))
	if totalPage <= 0 {
		totalPage = 1
	}
	return totalPage
}
func (page *KmgPage) HtmlRender() string {
	return tplPager(page)
}

// 一共有多少项
func (page *KmgPage) GetTotalItem() int {
	return page.TotalItem
}

// 是否有向前的按钮
func (page *KmgPage) IsBeforePageActive() bool {
	return page.CurrentPage-1 >= 1
}

// 是否有后向的按钮
func (page *KmgPage) IsAfterPageActive() bool {
	return page.CurrentPage+1 <= page.GetTotalPage()
}

// MySQL 数据库的偏移量
func (page *KmgPage) GetMysqlOffset() int {
	ret := (page.CurrentPage - 1) * page.ItemPerPage
	return ret
}

// 前一页的 Url
func (page *KmgPage) GetBeforePageUrl() string {
	pageNumber := page.CurrentPage - 1
	if pageNumber < 1 {
		pageNumber = 1
	}
	return page.GetUrlWithPage(pageNumber)
}

// 后一页的 Url
func (page *KmgPage) GetAfterPageUrl() string {
	pageNumber := page.CurrentPage + 1
	if pageNumber > page.GetTotalPage() {
		pageNumber = page.GetTotalPage()
	}
	return page.GetUrlWithPage(pageNumber)
}

// 中间显示的分页的数据
func (page *KmgPage) GetShowPageArray() []UrlParam {
	// 页面比显示数据量还少，快速通道
	if page.GetTotalPage() <= page.ShowPageNum {
		return page.getShowPageArrayFromNum(1, page.GetTotalPage())
	}
	start := page.CurrentPage - page.ShowPageNum/2
	end := page.CurrentPage + page.ShowPageNum/2 - 1
	if start < 1 {
		end += 1 - start
		start = 1
	}
	if end > page.GetTotalPage() {
		start = start - (end - page.GetTotalPage()) + 1
		end = page.GetTotalPage()
	}
	return page.getShowPageArrayFromNum(start, end)
}

type UrlParam struct {
	IsCurrent bool
	PageNum   int
	Url       string
}

func (page *KmgPage) getShowPageArrayFromNum(start int, end int) []UrlParam {
	var output []UrlParam
	var param UrlParam
	for i := start; i < end; i++ {
		url := page.GetUrlWithPage(i)
		param.IsCurrent = (i == page.CurrentPage)
		param.PageNum = i
		param.Url = url
		output[i] = param
	}
	return output
}

// 获取页面的 Url
func (page *KmgPage) GetUrlWithPage(pageNum int) string {
	return kmgHttp.MustSetParameterToUrl(page.BaseUrl, page.PageKeyName, strconv.Itoa(pageNum))
}
