package kmgBootstrap

import (
	"bytes"
	"github.com/bronze1man/kmg/kmgXss"
)

func tplInput(config Input) string {
	var _buf bytes.Buffer
	_buf.WriteString(`        <input type="`)
	_buf.WriteString(kmgXss.H(config.Type))
	_buf.WriteString(`" autocomplete="off" class="form-control"
            `)
	if config.ReadOnly {
		_buf.WriteString(` readonly `)
	}
	_buf.WriteString(`               name="`)
	_buf.WriteString(kmgXss.H(config.Name))
	_buf.WriteString(`"
               value="`)
	_buf.WriteString(kmgXss.H(config.Value))
	_buf.WriteString(`"/>`)
	return _buf.String()
}
