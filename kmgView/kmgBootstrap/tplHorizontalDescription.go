package kmgBootstrap

import (
	"bytes"

	"github.com/sipin/gorazor/gorazor"
)

func tplHorizontalDescription(conf HorizontalDescription) string {
	var _buffer bytes.Buffer
	_buffer.WriteString("\n<dl class=\"dl-horizontal\">\n    <dt>")
	_buffer.WriteString(gorazor.HTMLEscape(conf.Key))
	_buffer.WriteString("</dt>\n    <dd>")
	_buffer.WriteString(gorazor.HTMLEscape(conf.Value))
	_buffer.WriteString("</dd>\n</dl>")

	return _buffer.String()
}
