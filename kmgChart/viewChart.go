package kmgChart

import (
	"bytes"
	"github.com/bronze1man/kmg/kmgXss"
)

func viewLine(in Line) string {
	var _buf bytes.Buffer
	_buf.WriteString(`<div id="`)
	_buf.WriteString(kmgXss.H(in.domId))
	_buf.WriteString(`" style="width: `)
	_buf.WriteString(kmgXss.H(in.Width))
	_buf.WriteString(`;height: `)
	_buf.WriteString(kmgXss.H(in.Height))
	_buf.WriteString(`;">
</div>
<script>
    (function () {
        var line = echarts.init(document.getElementById(`)
	_buf.WriteString(kmgXss.Jsonv(in.domId))
	_buf.WriteString(`));
        var option = `)
	_buf.WriteString(in.GetOptionString())
	_buf.WriteString(`;
        line.setOption(option);
    })();
</script>`)
	return _buf.String()
}
