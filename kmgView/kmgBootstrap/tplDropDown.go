package kmgBootstrap

import (
	"bytes"
	"github.com/bronze1man/kmg/kmgXss"
)

func tplDropDown(d DropDown) string {
	var _buf bytes.Buffer
	_buf.WriteString(`<div style="display: inline-block" class="dropdown">
    <span id="d-`)
	_buf.WriteString(kmgXss.H(d.id))
	_buf.WriteString(`" data-toggle="dropdown" aria-haspopup="true" aria-expanded="true">
        `)
	_buf.WriteString(d.Title.HtmlRender())
	_buf.WriteString(`    </span>
    <ul class="dropdown-menu" aria-labelledby="d-`)
	_buf.WriteString(kmgXss.H(d.id))
	_buf.WriteString(`">
        `)
	for _, o := range d.OptionList {
		_buf.WriteString(`        `)
		_buf.WriteString(tplNavBarNode(o, 1))
		_buf.WriteString(`        `)
	}
	_buf.WriteString(`    </ul>
</div>`)
	return _buf.String()
}
