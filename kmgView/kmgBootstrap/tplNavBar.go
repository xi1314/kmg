package kmgBootstrap

import (
	"bytes"
	"github.com/bronze1man/kmg/kmgXss"
)

func tplNavBar(n NavBar) string {
	var _buf bytes.Buffer
	_buf.WriteString(`<style>
   body {
       padding-top: 71px!important;
   }
</style>
<nav class="navbar navbar-inverse navbar-fixed-top">
    <div class="container-fluid">
        <!-- Brand and toggle get grouped for better mobile display -->
        <div class="navbar-header">
            <button type="button" class="navbar-toggle collapsed" data-toggle="collapse" data-target="#kmgBootstrapNavBarMainContent" aria-expanded="false">
                <span class="sr-only">Toggle navigation</span>
                <span class="icon-bar"></span>
                <span class="icon-bar"></span>
                <span class="icon-bar"></span>
            </button>
            <span class="navbar-brand">`)
	_buf.WriteString(n.Title.HtmlRender())
	_buf.WriteString(`</span>
        </div>

        <!-- Collect the nav links, forms, and other content for toggling -->
        <div class="collapse navbar-collapse" id="kmgBootstrapNavBarMainContent">
            <ul class="nav navbar-nav">
                `)
	for _, o := range n.OptionList {
		_buf.WriteString(`                <li `)
		if o.Name == n.ActiveName {
			_buf.WriteString(`class="active"`)
		}
		_buf.WriteString(` >
                    <a href="`)
		_buf.WriteString(kmgXss.H(o.Url))
		_buf.WriteString(`">`)
		_buf.WriteString(kmgXss.H(o.Name))
		_buf.WriteString(`</a>
                </li>
                `)
	}
	_buf.WriteString(`            </ul>
            <ul class="nav navbar-nav navbar-right">
                `)
	for _, o := range n.RightOptionList {
		_buf.WriteString(`                <li `)
		if o.Name == n.ActiveName {
			_buf.WriteString(`class="active"`)
		}
		_buf.WriteString(` >
                    <a href="`)
		_buf.WriteString(kmgXss.H(o.Url))
		_buf.WriteString(`">`)
		_buf.WriteString(kmgXss.H(o.Name))
		_buf.WriteString(`</a>
                </li>
                `)
	}
	_buf.WriteString(`            </ul>
        </div><!-- /.navbar-collapse -->
    </div><!-- /.container-fluid -->
</nav>`)
	return _buf.String()
}
