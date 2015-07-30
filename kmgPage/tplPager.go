package kmgPage

import (
	"bytes"
	"net/url"

	"github.com/sipin/gorazor/gorazor"
)

func tplPager(kmgPage *KmgPage) string {
	var _buffer bytes.Buffer
	_buffer.WriteString("\n<ul class=\"pagination\">\n    ")
	if !kmgPage.IsBeforePageActive() {

		_buffer.WriteString("<li class=\"disabled\" >\n            <a href=\"javascript:\">&laquo;</a>\n        </li>")

	} else {

		_buffer.WriteString("<li>\n            <a href=\"")
		_buffer.WriteString(gorazor.HTMLEscape(kmgPage.GetBeforePageUrl()))
		_buffer.WriteString("\">&laquo;</a>\n        </li>")

	}
	_buffer.WriteString("\n    ")
	for _, opt := range kmgPage.GetShowPageArray() {

		_buffer.WriteString("<li class=\"")
		if opt.IsCurrent {

			_buffer.WriteString("active")

		}
		_buffer.WriteString("\">\n        <a href=\"")
		if opt.IsCurrent {

			_buffer.WriteString("javascript:")

		} else {

			_buffer.WriteString(gorazor.HTMLEscape(opt.Url))

		}
		_buffer.WriteString("\">")
		_buffer.WriteString(gorazor.HTMLEscape(opt.PageNum))
		_buffer.WriteString("\n            <span class=\"sr-only\">(current)</span></a>\n        </li>")

	}
	_buffer.WriteString("\n    ")
	if !kmgPage.IsAfterPageActive() {

		_buffer.WriteString("<li class=\"disabled\" >\n            <a href=\"javascript:\">&raquo;</a>\n        </li>")

	} else {

		_buffer.WriteString("<li>\n            <a href=\"")
		_buffer.WriteString(gorazor.HTMLEscape(kmgPage.GetAfterPageUrl()))
		_buffer.WriteString("\">&raquo;</a>\n        </li>")

	}
	_buffer.WriteString("\n    <li>\n        <form action=\"")
	_buffer.WriteString(gorazor.HTMLEscape(kmgPage.BaseUrl))
	_buffer.WriteString("\" method=\"GET\" style=\"position: relative;margin-left:10px;float:left;\">\n            ")
	u, _ := url.ParseRequestURI(kmgPage.BaseUrl)

	_buffer.WriteString("\n            ")
	for key, valueList := range u.Query() {
		if key == kmgPage.PageKeyName {
			continue
		}

		_buffer.WriteString("<input type=\"hidden\" name=\"")
		_buffer.WriteString(gorazor.HTMLEscape(key))
		_buffer.WriteString("\" value=\"")
		_buffer.WriteString(gorazor.HTMLEscape(valueList[0]))
		_buffer.WriteString("\"/>")

	}
	_buffer.WriteString("\n            <input type=\"text\" class=\"form-control\"\n                   style=\"width:30px;height: 29px;padding: 2px 2px;display:inline;text-align:center;position: relative;top:1px;\"\n                   name=\"")
	_buffer.WriteString(gorazor.HTMLEscape(kmgPage.PageKeyName))
	_buffer.WriteString("\" value=\"")
	_buffer.WriteString(gorazor.HTMLEscape(kmgPage.CurrentPage))
	_buffer.WriteString("\"/>\n            <input type=\"submit\" class=\"btn btn-primary\" style=\"padding: 5px 12px\" value=\"跳转至该页\"/>\n        </form>\n    </li>\n</ul>")

	return _buffer.String()
}
