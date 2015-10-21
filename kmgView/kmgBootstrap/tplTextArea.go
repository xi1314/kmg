package kmgBootstrap

import (
	"bytes"
	"github.com/bronze1man/kmg/kmgXss"
)

func tplTextArea(config TextArea) string {
	var _buf bytes.Buffer
	_buf.WriteString(`        <textarea autocomplete="off" class="form-control"
                    `)
	if config.ReadOnly {
		_buf.WriteString(` readonly `)
	}
	_buf.WriteString(`                name="`)
	_buf.WriteString(kmgXss.H(config.Name))
	_buf.WriteString(`" cols="30" rows="10" >`)
	_buf.WriteString(kmgXss.H(config.Value))
	_buf.WriteString(`</textarea>`)
	return _buf.String()
}
