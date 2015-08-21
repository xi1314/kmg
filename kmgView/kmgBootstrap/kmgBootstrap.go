package kmgBootstrap

import (
	"fmt"

	"github.com/bronze1man/kmg/kmgView"
	"github.com/bronze1man/kmg/kmgXss"
	"strconv"
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
	Id          string
	ClassName   string
}

func (a A) HtmlRender() string {
	s := `<a href="` + kmgXss.H(a.Href) + `" class="` + kmgXss.H(a.ClassName) + `"`
	if a.IsNewWindow {
		s += ` target="_blank" `
	}
	if a.Id != "" {
		s += ` id="` + kmgXss.H(a.Id) + `" `
	}
	return s + `>` + kmgXss.H(a.Title) + `</a>`
}

type NavBar struct {
	Title           kmgView.HtmlRenderer //可以是 Log 或者文字
	ActiveName      string
	OptionList      []NavTabOption
	RightOptionList []NavTabOption
}

func (navBar NavBar) HtmlRender() string {
	return tplNavBar(navBar)
}

type TextColor string

const (
	TextMuted   TextColor = "text-muted"
	TextPrimary TextColor = "text-primary"
	TextSuccess TextColor = "text-success"
	TextInfo    TextColor = "text-info"
	TextWarning TextColor = "text-warning"
	TextDanger  TextColor = "text-danger"
)

type BgColor string

const (
	BgPrimary TextColor = "bg-primary"
	BgSuccess TextColor = "bg-success"
	BgInfo    TextColor = "bg-info"
	BgWarning TextColor = "bg-warning"
	BgDanger  TextColor = "bg-danger"
)

type Icon struct {
	IconName      string
	IconColor     TextColor
	AttributeNode kmgView.HtmlRenderer
}

func (icon Icon) HtmlRender() string {
	return tplIcon(icon)
}

type PlacementType string

const (
	PlacementTypeLeft   PlacementType = "left"
	PlacementTypeRight  PlacementType = "right"
	PlacementTypeTop    PlacementType = "top"
	PlacementTypeBottom PlacementType = "bottom"
	PlacementTypeAuto   PlacementType = "auto"
)

type PopoverType string

const (
	PopoverTypePopover PopoverType = "popover"
	PopoverTypeTooltip PopoverType = "tooltip"
)

//注意
//Bootstrap 有两种提示小窗口，一种叫 Tooltip，一种叫 Popover，两者数据结构基本一致，这里为了简单，直接合并为一种 kmgBootstrap 的结构 Popover
type Popover struct {
	Type      PopoverType
	Title     string
	Content   string
	Placement PlacementType
}

func (p Popover) HtmlRender() string {
	if p.Type == "" {
		p.Type = PopoverTypeTooltip
	}
	if p.Placement == "" {
		p.Placement = PlacementTypeAuto
	}
	return tplPopover(p)
}

type ButtonSize string

const (
	ButtonSizeLarge      ButtonSize = "btn-lg"
	ButtonSizeDefault    ButtonSize = ""
	ButtonSizeSmall      ButtonSize = "btn-sm"
	ButtonSizeExtraSmall ButtonSize = "btn-xs"
)

type ButtonColor string

const (
	ButtonColorDefault ButtonColor = "btn-default"
	ButtonColorPrimary ButtonColor = "btn-primary"
	ButtonColorSuccess ButtonColor = "btn-success"
	ButtonColorInfo    ButtonColor = "btn-info"
	ButtonColorWarning ButtonColor = "btn-warning"
	ButtonColorDanger  ButtonColor = "btn-danger"
	ButtonColorLink    ButtonColor = "btn-link"
)

type ButtonType string

const (
	ButtonTypeA      ButtonType = "a"
	ButtonTypeButton ButtonType = "button"
)

type Button struct {
	Type          ButtonType
	Url           string
	Color         ButtonColor
	Size          ButtonSize
	Content       kmgView.HtmlRenderer
	AttributeNode kmgView.HtmlRenderer
	ClassName     string
	Id            string
}

func (b Button) HtmlRender() string {
	if b.Color == "" {
		b.Color = ButtonColorDefault
	}
	if b.Type == "" {
		b.Type = ButtonTypeButton
	}
	return tplButton(b)
}

type LabelColor string

const (
	LabelColorDefault LabelColor = "label-default"
	LabelColorPrimary LabelColor = "label-primary"
	LabelColorSuccess LabelColor = "label-success"
	LabelColorInfo    LabelColor = "label-info"
	LabelColorWarning LabelColor = "label-warning"
	LabelColorDanger  LabelColor = "label-danger"
)

type Label struct {
	Color   LabelColor
	Content kmgView.HtmlRenderer
}

func (l Label) HtmlRender() string {
	if l.Color == "" {
		l.Color = LabelColorDefault
	}
	if l.Content == nil {
		return ""
	}
	return `<span class="label ` + kmgXss.H(string(l.Color)) + `">` + l.Content.HtmlRender() + `</span>`
}

func Blank(num int) kmgView.HtmlRenderer {
	b := kmgView.HtmlRendererList{}
	i := 0
	for {
		if i >= num {
			break
		}
		i++
		b = append(b, kmgView.Html("&nbsp;"))
	}
	return b
}

func H(index int, content kmgView.HtmlRenderer) kmgView.HtmlRenderer {
	if content == nil {
		return kmgView.String("")
	}
	if index == 0 {
		index = 1
	}
	iStr := strconv.Itoa(index)
	return kmgView.Html(`<h` + iStr + `>` + content.HtmlRender() + `</h` + iStr + `>`)
}
