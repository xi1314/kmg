package kmgView

import "bytes"

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
