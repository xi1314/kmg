package kmgBootstrap

import (
	"bytes"
	"github.com/sipin/gorazor/gorazor"
)

func tplPanel(panel Panel) string {
	var _buffer bytes.Buffer
	_buffer.WriteString("\n<div class=\"panel panel-default\">\n    <div class=\"panel-heading\">\n        <h4>")
	_buffer.WriteString(gorazor.HTMLEscape(panel.Title))
	_buffer.WriteString("</h4>\n    </div>\n    <div class=\"panel-body\">\n        ")
	_buffer.WriteString((panel.Body.HtmlRender()))
	_buffer.WriteString("\n    </div>\n</div>")

	return _buffer.String()
}
