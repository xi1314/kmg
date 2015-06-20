package kmgBootstrap

import (
	"bytes"
)

func tplHorizontalHtmlRenderList(conf HorizontalHtmlRenderList) string {
	var _buffer bytes.Buffer
	_buffer.WriteString("\n<table style=\"width: 100%;\">\n    <tr>\n        ")
	for _, elem := range conf {

		_buffer.WriteString("<td>\n                ")
		_buffer.WriteString((elem.HtmlRender()))
		_buffer.WriteString("\n            </td>")

	}
	_buffer.WriteString("\n    </tr>\n</table>")

	return _buffer.String()
}
