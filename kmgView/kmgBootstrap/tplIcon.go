package kmgBootstrap

import (
	"bytes"
	"github.com/bronze1man/kmg/kmgXss"
)

func tplIcon(icon Icon) string {
	var _buf bytes.Buffer
	_buf.WriteString(` <i data-toggle="tooltip" data-original-title="`)
	_buf.WriteString(kmgXss.H(icon.Title))
	_buf.WriteString(`" class="`)
	_buf.WriteString(kmgXss.H(string(icon.IconColor)))
	_buf.WriteString(` fa fa-`)
	_buf.WriteString(kmgXss.H(icon.IconName))
	_buf.WriteString(`"></i>`)
	return _buf.String()
}
