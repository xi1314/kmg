package kmgBootstrap

import (
	"bytes"
	"github.com/sipin/gorazor/gorazor"
)

func tplNavTabList(conf NavTabList) string {
	var _buffer bytes.Buffer
	_buffer.WriteString("\n<ul class=\"nav nav-pills\">\n    ")
	for _, opt := range conf.OptionList {

		_buffer.WriteString("<li\n        ")
		if opt.Name == conf.ActiveName {

			_buffer.WriteString("class=\"active\"")

		}
		_buffer.WriteString("\n            >\n        <a href=\"")
		_buffer.WriteString(gorazor.HTMLEscape(opt.Url))
		_buffer.WriteString("\">")
		_buffer.WriteString(gorazor.HTMLEscape(opt.Name))
		_buffer.WriteString("</a>\n    </li>\n    ")
	}
	_buffer.WriteString("\n</ul>")

	return _buffer.String()
}
