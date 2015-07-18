package kmgBootstrap

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgView"
	"github.com/bronze1man/kmg/kmgXss"
)

type Panel struct {
	Title string
	Body  kmgView.HtmlRenderer
}

func (p Panel) HtmlRender() string {
	return tplPanel(p)
}

type HorizontalDescription struct {
	Key   string
	Value string
}

func (p HorizontalDescription) HtmlRender() string {
	return tplHorizontalDescription(p)
}

type Table struct {
	Caption   kmgView.HtmlRenderer
	TitleList []kmgView.HtmlRenderer
	DataList  [][]kmgView.HtmlRenderer
}

func (p Table) HtmlRender() string {
	return tplTable(p)
}

func NewTable(TitleList []string, DataList [][]string) Table {
	t := Table{}
	for _, title := range TitleList {
		t.TitleList = append(t.TitleList, kmgView.String(title))
	}
	for _, row := range DataList {
		renderRow := []kmgView.HtmlRenderer{}
		for _, value := range row {
			renderRow = append(renderRow, kmgView.String(value))
		}
		t.DataList = append(t.DataList, renderRow)
	}
	return t
}

// 1.这个函数顺序会变,基本不靠谱.
// 2.直接从数据库取的数据,字段信息没有翻译.
func NewTableFromMapList(m []map[string]string) Table {
	t := Table{}
	if len(m) == 0 {
		return t
	}
	i := 0
	for k := range m[0] {
		// TODO 排序
		t.TitleList = append(t.TitleList, kmgView.String(k))
		i++
	}
	for _, row := range m {
		viewRow := []kmgView.HtmlRenderer{}
		for _, k := range t.TitleList {
			viewRow = append(viewRow, kmgView.String(row[string(k.(kmgView.String))]))
		}
		t.DataList = append(t.DataList, viewRow)
	}
	return t
}

type HorizontalHtmlRenderList []kmgView.HtmlRenderer

func (p HorizontalHtmlRenderList) HtmlRender() string {
	return tplHorizontalHtmlRenderList(p)
}

type NavTabList struct {
	ActiveName string //选中的名称
	OptionList []NavTabOption
}

type NavTabOption struct {
	Name string
	Url  string
}

func (p NavTabList) HtmlRender() string {
	return tplNavTabList(p)
}

type Image struct {
	Src       string
	MaxHeight int
	MaxWidth  int
}

func (p Image) HtmlRender() string {
	return tplImage(p)
}

func (p Image) getStyle() string {
	if p.MaxHeight == 0 {
		p.MaxHeight = 100
	}
	if p.MaxWidth == 0 {
		p.MaxWidth = 100
	}
	return fmt.Sprintf("max-height: %dpx;max-width: %dpx;", p.MaxHeight, p.MaxWidth)
}

type A struct {
	Href        string
	Title       string
	IsNewWindow bool
}

func (p A) HtmlRender() string {
	s := `<a href="` + kmgXss.H(p.Href) + `"`
	if p.IsNewWindow {
		s += ` target="_blank" `
	}
	return s + `>` + kmgXss.H(p.Title) + `</a>`
}
