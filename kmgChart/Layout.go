package kmgChart

import (
	"bytes"
	"github.com/bronze1man/kmg/kmgView"
)

func Layout(r ...kmgView.HtmlRenderer) string {
	var _buf bytes.Buffer
	_buf.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Kmg Chart</title>
    <script src="http://7xjyrb.com2.z0.glb.qiniucdn.com/echarts.js"></script>
</head>
<body>`)
	for _, v := range r {
		_buf.WriteString(`    `)
		_buf.WriteString(v.HtmlRender())
		_buf.WriteString(``)
	}
	_buf.WriteString(`</body>
</html>`)
	return _buf.String()
}
