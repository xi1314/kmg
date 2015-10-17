package kmgBootstrap

import (
	"bytes"
	"github.com/bronze1man/kmg/kmgXss"
)

func tplMenu(menu Menu) string {
	var _buf bytes.Buffer
	_buf.WriteString(`<nav class="navbar navbar-inverse navbar-fixed-top">
    <div class="container">
        <div class="navbar-header">
            <button type="button" class="navbar-toggle collapsed" data-toggle="collapse" data-target="#navbar" aria-expanded="false" aria-controls="navbar">
                <span class="sr-only">Toggle navigation</span>
                <span class="icon-bar"></span>
                <span class="icon-bar"></span>
                <span class="icon-bar"></span>
            </button>
            <span class="navbar-brand">`)
	_buf.WriteString(kmgXss.H(menu.Title))
	_buf.WriteString(`</span>
        </div>
        <div id="navbar" class="navbar-collapse collapse">
            <ul class="nav navbar-nav">
                `)
	for _, subMenuNode := range menu.NodeList {
		_buf.WriteString(`                    `)
		_buf.WriteString(tplMenuNode(subMenuNode, 0))
		_buf.WriteString(`                `)
	}
	_buf.WriteString(`            </ul>
        </div><!--/.nav-collapse -->
    </div>
</nav>
<div style="height: 50px;"></div>`)
	return _buf.String()
}
