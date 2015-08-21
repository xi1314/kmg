package kmgBootstrap

import (
	"bytes"
	"github.com/bronze1man/kmg/kmgXss"
)

func tplIcon(icon Icon) string {
	var _buf bytes.Buffer
	_buf.WriteString(`<i class="`)
	_buf.WriteString(kmgXss.H(string(icon.IconColor)))
	_buf.WriteString(` fa fa-`)
	_buf.WriteString(kmgXss.H(icon.IconName))
	_buf.WriteString(`"
    `)
	if icon.AttributeNode != nil {
		_buf.WriteString(`        `)
		_buf.WriteString(icon.AttributeNode.HtmlRender())
		_buf.WriteString(`    `)
	}
	_buf.WriteString(`></i>`)
	return _buf.String()
}
