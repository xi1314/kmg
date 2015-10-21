package kmgBootstrap

import (
	"bytes"
	"github.com/bronze1man/kmg/kmgXss"
)

func tplNavBarNode(node NavBarNode, level int) string {
	var _buf bytes.Buffer
	_buf.WriteString(`    `)
	if len(node.ChildList) == 0 {
		_buf.WriteString(`        <li><a href="`)
		_buf.WriteString(kmgXss.H(node.Url))
		_buf.WriteString(`">`)
		_buf.WriteString(kmgXss.H(node.Name))
		_buf.WriteString(`</a></li>
    `)
	} else if level == 0 {
		_buf.WriteString(`        <li class="dropdown">
            <a href="javascript:" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-haspopup="true" aria-expanded="false">`)
		_buf.WriteString(kmgXss.H(node.Name))
		_buf.WriteString(` <span class="caret"></span></a>
            <ul class="dropdown-menu">
                `)
		for _, subNode := range node.ChildList {
			_buf.WriteString(`                    `)
			_buf.WriteString(tplNavBarNode(subNode, level+1))
			_buf.WriteString(`                `)
		}
		_buf.WriteString(`            </ul>
        </li>
    `)
	} else {
		_buf.WriteString(`        <li class="dropdown-submenu">
            <a href="javascript:">`)
		_buf.WriteString(kmgXss.H(node.Name))
		_buf.WriteString(`</a>
            <ul class="dropdown-menu">
                `)
		for _, subNode := range node.ChildList {
			_buf.WriteString(`                `)
			_buf.WriteString(tplNavBarNode(subNode, level+1))
			_buf.WriteString(`                `)
		}
		_buf.WriteString(`            </ul>
        </li>`)
	}
	return _buf.String()
}
