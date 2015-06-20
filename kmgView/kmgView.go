package kmgView

import (
	"bytes"
	"github.com/bronze1man/kmg/kmgXss"
	"github.com/bronze1man/kmg/typeTransform"
)

type HtmlRenderer interface {
	HtmlRender() string
}

type HtmlRendererList []HtmlRenderer

func (l HtmlRendererList) HtmlRender() string {
	var _buffer bytes.Buffer
	for _, renderer := range l {
		_buffer.WriteString(renderer.HtmlRender())
	}
	return _buffer.String()
}

//使用泛型帮你解决各种无聊的类型转换
func NewHtmlRendererListFromList(obj interface{}) HtmlRendererList {
	out := HtmlRendererList{}
	typeTransform.MustTransform(obj, &out)
	return out
}

type String string

func (s String) HtmlRender() string {
	return kmgXss.H(string(s))
}

type Html string

func (s Html) HtmlRender() string {
	return string(s)
}
