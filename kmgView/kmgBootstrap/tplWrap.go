package kmgBootstrap

import (
	"bytes"

	"github.com/bronze1man/kmg/kmgXss"
)

func tplWrap(w Wrap) string {
	var _buf bytes.Buffer
	_buf.WriteString(`<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>`)
	_buf.WriteString(kmgXss.H(w.Title))
	_buf.WriteString(`</title>
    `)
	_buf.WriteString(w.Head.HtmlRender())
	_buf.WriteString(`</head>
<body style="padding: 20px;">`)
	_buf.WriteString(w.Body.HtmlRender())
	_buf.WriteString(`</body>
</html>`)
	return _buf.String()
}
