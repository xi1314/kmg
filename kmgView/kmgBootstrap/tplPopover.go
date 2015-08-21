package kmgBootstrap

import (
	"bytes"
	"github.com/bronze1man/kmg/kmgXss"
)

func tplPopover(popover Popover) string {
	var _buf bytes.Buffer
	_buf.WriteString(`data-toggle="`)
	_buf.WriteString(kmgXss.H(string(popover.Type)))
	_buf.WriteString(`"
data-placement="`)
	_buf.WriteString(kmgXss.H(string(popover.Placement)))
	_buf.WriteString(`"
title="`)
	_buf.WriteString(popover.Title)
	_buf.WriteString(`"
data-content="`)
	_buf.WriteString(popover.Content)
	_buf.WriteString(`"
data-html="true"`)
	return _buf.String()
}
