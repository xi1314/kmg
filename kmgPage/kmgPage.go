package kmgPage

import (
	"github.com/bronze1man/kmg/kmgSql"
	"github.com/bronze1man/kmg/kmgSql/MysqlAst"
	"math"
	"net/url"
	"strconv"
)

type kmgPage struct {
	ItemPerPage int                 // 每页数据量
	ShowPageNum int                 // 显示分页的页面的时候，显示出来的页数 如 10，11，12，13，14，15 就是显示了6个页
	PageKeyName string              // 页面参数的名字
	CurrentPage int                 // 当前页数，默认为1
	TotalItem   int                 // 总数据量
	BaseUrl     string              // 不包含页码参数的url,具体渲染url的时候会把页面参数加上
	Data        []map[string]string // 本次分页查询到的数据
	TotalPage   int                 // 共有多少个页面
}

func (page *kmgPage) CreateFromSelectCommand(selectCommand MysqlAst.SelectCommand, baseUrl string, itemPerPage int) *kmgPage {
	page.init()
	page.BaseUrl = baseUrl
	page.ItemPerPage = itemPerPage
	return page.runSelectCommand(selectCommand)
}

func (page *kmgPage) CreateFromData(data []map[string]string, baseUrl string, itemPerPage int) *kmgPage {
	page.init()
	page.BaseUrl = baseUrl
	page.ItemPerPage = itemPerPage
	page.TotalItem = len(data)
	return page
}

func (page *kmgPage) runSelectCommand(selectCommand MysqlAst.SelectCommand) *kmgPage {
	if page.BaseUrl == "" {
		panic("runSelectCommand need baseUrl parameter")
	}
	output, parameterList := selectCommand.GetPrepareParameter()
	countData, err := kmgSql.Query("SELECT COUNT(*) AS c FROM ("+output+") as View", parameterList[0])
	if err != nil {
		panic(err)
	}
	page.TotalItem, err = strconv.Atoi(countData[0]["c"])
	if err != nil {
		panic(err)
	}
	dataSelect := selectCommand.Copy()
	dataSelect = dataSelect.Limit(string(page.GetMysqlOffset()) + "," + string(page.ItemPerPage))
	page.Data = kmgSql.RunSelectCommand(dataSelect)
	return page
}

func (page *kmgPage) init() {
	page.ItemPerPage = 10
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

func (page *kmgPage) GetTotalPage() int {
	totalPage := int(math.Ceil(float64(page.TotalItem) / float64(page.ItemPerPage)))
	if totalPage <= 0 {
		totalPage = 1
	}
	return totalPage
}

// 一共有多少项
func (page *kmgPage) GetTotalItem() int {
	return page.TotalItem
}

// 是否有向前的按钮
func (page *kmgPage) IsBeforePageActive() bool {
	return page.CurrentPage-1 >= 1
}

// 是否有后向的按钮
func (page *kmgPage) IsAfterPageActive() bool {
	return page.CurrentPage+1 <= page.GetTotalPage()
}

// MySQL 数据库的偏移量
func (page *kmgPage) GetMysqlOffset() int {
	ret := (page.CurrentPage - 1) * page.ItemPerPage
	return ret
}

// 前一页的 Url
func (page *kmgPage) GetBeforePageUrl() string {
	pageNumber := page.CurrentPage - 1
	if pageNumber < 1 {
		pageNumber = 1
	}
	return page.GetUrlWithPage(pageNumber)
}

// 后一页的 Url
func (page *kmgPage) GetAfterPageUrl() string {
	pageNumber := page.CurrentPage + 1
	if pageNumber > page.GetTotalPage() {
		pageNumber = page.GetTotalPage()
	}
	return page.GetUrlWithPage(pageNumber)
}

// 中间显示的分页的数据
func (page *kmgPage) GetShowPageArray() []urlParam {
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

type urlParam struct {
	IsCurrent bool
	PageNum   int
	Url       string
}

func (page *kmgPage) getShowPageArrayFromNum(start int, end int) []urlParam {
	var output []urlParam
	var param urlParam
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
func (page *kmgPage) GetUrlWithPage(pageNum int) string {
	v := url.Values{}
	v.Set(page.PageKeyName, string(pageNum))
	return page.BaseUrl + "?" + v.Encode()
}
