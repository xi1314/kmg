package kmgBootstrap

import (
	"bytes"
	"github.com/bronze1man/kmg/kmgXss"
)

func tplNavTabList(conf NavTabList) string {
	var _buf bytes.Buffer
	if conf.CustomClass == "" {
		conf.CustomClass = "nav-pills"
	}
	_buf.WriteString(`    <ul class="nav `)
	_buf.WriteString(kmgXss.H(conf.CustomClass))
	_buf.WriteString(`">
   `)
	for _, opt := range conf.OptionList {
		_buf.WriteString(`    <li
        `)
		if opt.Name == conf.ActiveName {
			_buf.WriteString(`           class="active")
        `)
		}
		_buf.WriteString(`            >
        <a href="`)
		_buf.WriteString(kmgXss.H(opt.Url))
		_buf.WriteString(`">`)
		_buf.WriteString(kmgXss.H(opt.Name))
		_buf.WriteString(`</a>
    </li>
    `)
	}
	_buf.WriteString(`</ul>`)
	return _buf.String()
}
