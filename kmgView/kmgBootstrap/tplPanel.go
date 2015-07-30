package kmgBootstrap

import (
	"bytes"
	"github.com/bronze1man/kmg/kmgXss"
)

func tplPanel(panel Panel) string {
	var _buf bytes.Buffer
	_buf.WriteString(`<div class="panel panel-default">
    <div class="panel-heading">
        `)
	_buf.WriteString(kmgXss.H(panel.Title))
	_buf.WriteString(`    </div>
    <div class="panel-body">
        `)
	_buf.WriteString(panel.Body.HtmlRender())
	_buf.WriteString(`    </div>
</div>`)
	return _buf.String()
}
