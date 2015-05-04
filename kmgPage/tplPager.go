package kmgPage

import (
	"bytes"
	"github.com/sipin/gorazor/gorazor"
	"net/url"
)

func tplPager(kmgPage *KmgPage) string {
	var _buffer bytes.Buffer
	_buffer.WriteString("\n<ul class=\"pagination\">\n    <li class=\"disabled\">\n        <a href=\"\">\n            共 ")
	_buffer.WriteString(gorazor.HTMLEscape(kmgPage.TotalPage))
	_buffer.WriteString(" 页\n        </a>\n    </li>\n    <li class=\"")
	if !kmgPage.IsBeforePageActive() {

		_buffer.WriteString("disabled")

	}
	_buffer.WriteString("\" >\n        <a href=\"")
	_buffer.WriteString(gorazor.HTMLEscape(kmgPage.GetBeforePageUrl()))
	_buffer.WriteString("\">&laquo;</a>\n    </li>\n    ")
	for _, opt := range kmgPage.GetShowPageArray() {

		_buffer.WriteString("<li class=\"")
		if opt.IsCurrent {

			_buffer.WriteString("active")

		}
		_buffer.WriteString("\">\n        <a href=\"")
		_buffer.WriteString(gorazor.HTMLEscape(opt.Url))
		_buffer.WriteString("\">")
		_buffer.WriteString(gorazor.HTMLEscape(opt.PageNum))
		_buffer.WriteString("\n            <span class=\"sr-only\">(current)</span></a>\n        </li>")

	}
	_buffer.WriteString("\n    <li class=\"")
	if !kmgPage.IsAfterPageActive() {

		_buffer.WriteString("disabled")

	}
	_buffer.WriteString("\">\n        <a href=\"")
	_buffer.WriteString(gorazor.HTMLEscape(kmgPage.GetAfterPageUrl()))
	_buffer.WriteString("\">&raquo;</a>\n    </li>\n    <li>\n        <form action=\"")
	_buffer.WriteString(gorazor.HTMLEscape(kmgPage.BaseUrl))
	_buffer.WriteString("\" method=\"GET\" style=\"position: relative;margin-left:10px;float:left;\">\n            ")
	array, err := url.ParseRequestURI(kmgPage.BaseUrl)
	if err != nil {

	}

	_buffer.WriteString("\n            ")
	for key, value := range array.Query() {

		_buffer.WriteString("<input type=\"hidden\" name=\"")
		_buffer.WriteString(gorazor.HTMLEscape(key))
		_buffer.WriteString("\" value=\"")
		_buffer.WriteString(gorazor.HTMLEscape(value))
		_buffer.WriteString("\"/>")

	}
	_buffer.WriteString("\n            <input type=\"text\" class=\"form-control\"\n                   style=\"width:30px;height: 29px;padding: 2px 2px;display:inline;text-align:center;position: relative;top:1px;\"\n                   name=\"")
	_buffer.WriteString(gorazor.HTMLEscape(kmgPage.PageKeyName))
	_buffer.WriteString("\" value=\"")
	_buffer.WriteString(gorazor.HTMLEscape(kmgPage.CurrentPage))
	_buffer.WriteString("\"/>\n            <input type=\"submit\" class=\"btn btn-primary\" style=\"padding: 5px 12px\" value=\"跳转至该页\"/>\n        </form>\n    </li>\n</ul>")

	return _buffer.String()
}
