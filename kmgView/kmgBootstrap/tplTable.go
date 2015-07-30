package kmgBootstrap

import (
	"bytes"
)

func tplTable(table Table) string {
	var _buffer bytes.Buffer
	_buffer.WriteString("\n<table class=\"table table-hover\">\n    ")
	if table.Caption != nil {

		_buffer.WriteString("<caption>")
		_buffer.WriteString((table.Caption.HtmlRender()))
		_buffer.WriteString("</caption>")

	}
	_buffer.WriteString("\n    ")
	if table.TitleList != nil {

		_buffer.WriteString("<thead>\n        <tr>\n            ")
		for _, title := range table.TitleList {

			_buffer.WriteString("<th>")
			_buffer.WriteString((title.HtmlRender()))
			_buffer.WriteString("</th>")

		}
		_buffer.WriteString("\n        </tr>\n        </thead>")

	}
	_buffer.WriteString("\n    <tbody>\n    ")
	for _, row := range table.DataList {

		_buffer.WriteString("<tr>\n            ")
		for _, cell := range row {

			_buffer.WriteString("<td>")
			_buffer.WriteString((cell.HtmlRender()))
			_buffer.WriteString("</td>")

		}
		_buffer.WriteString("\n        </tr>")

	}
	_buffer.WriteString("\n    </tbody>\n</table>")

	return _buffer.String()
}
