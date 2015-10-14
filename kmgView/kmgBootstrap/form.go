package kmgBootstrap
import "github.com/bronze1man/kmg/kmgView"

// 自带一个提交按钮的From表单.
// TODO 允许调用者配置,去掉自带的按钮.
type Form struct {
	Url       string
	IsGet     bool // 默认是POST
	InputList []kmgView.HtmlRenderer
}

func (f Form) HtmlRender() string {
	return tplForm(f)
}

type InputString struct{
	Name     string
	Value    string
	ShowName string
	Comment  string
	Need     bool
	ReadOnly bool
}

func (f InputString) HtmlRender() string {
	return tplInputString(f)
}