package kmgView

import (
	"bytes"
	"github.com/bronze1man/kmg/kmgXss"
)

type HtmlRenderer interface {
	HtmlRender() string
}

type HtmlRendererList []HtmlRenderer

func (l HtmlRendererList) HtmlRender() string {
	_buffer := &bytes.Buffer{}
	for _, renderer := range l {
		_buffer.WriteString(renderer.HtmlRender())
	}
	return _buffer.String()
}

type String string

func (s String) HtmlRender() string {
	return kmgXss.H(string(s))
}

type Html string

func (s Html) HtmlRender() string {
	return string(s)
}
