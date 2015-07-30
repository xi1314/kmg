package kmgBootstrap

import (
	"bytes"

	"github.com/sipin/gorazor/gorazor"
)

func tplImage(conf Image) string {
	var _buffer bytes.Buffer
	_buffer.WriteString("\n<img src=\"")
	_buffer.WriteString(gorazor.HTMLEscape(conf.Src))
	_buffer.WriteString("\" style=\"")
	_buffer.WriteString(gorazor.HTMLEscape(conf.getStyle()))
	_buffer.WriteString("\"/>")

	return _buffer.String()
}
