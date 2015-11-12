package kmgBootstrap

import (
	"bytes"
	"github.com/bronze1man/kmg/kmgXss"
)

func tplInputWrapVertical(config InputWrapVertical) string {
	var _buf bytes.Buffer
	_buf.WriteString(`<div class="form-group has-feedback">
    <label class="col-sm-2 control-label">`)
	_buf.WriteString(kmgXss.H(config.ShowName))
	_buf.WriteString(`        `)
	if config.Need {
		_buf.WriteString(`            <span style="color:red">*</span>
        `)
	}
	_buf.WriteString(`    </label>

    <div class="col-sm-8 `)
	if config.AppendTpl != nil {
		_buf.WriteString(` form-inline `)
	}
	_buf.WriteString(`">
        `)
	_buf.WriteString(config.Body.HtmlRender())
	_buf.WriteString(`        <span style="font-size:12px;color:red">
            `)
	if config.Comment != "" {
		_buf.WriteString(` 提示: `)
		_buf.WriteString(kmgXss.H(config.Comment))
		_buf.WriteString(` `)
	}
	_buf.WriteString(`        </span>
        `)
	if config.AppendTpl != nil {
		_buf.WriteString(`            `)
		_buf.WriteString(config.AppendTpl.HtmlRender())
		_buf.WriteString(`        `)
	}
	_buf.WriteString(`    </div>
</div>`)
	return _buf.String()
}
