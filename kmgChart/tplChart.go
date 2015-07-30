package kmgChart

import (
	"bytes"

	"github.com/bronze1man/kmg/kmgXss"
)

func tplChart(chartData Chart) string {
	var _buf bytes.Buffer
	_buf.WriteString(`<div id="`)
	_buf.WriteString(kmgXss.H(chartData.domId))
	_buf.WriteString(`" style="width: `)
	_buf.WriteString(kmgXss.H(chartData.Width))
	_buf.WriteString(`;height: `)
	_buf.WriteString(kmgXss.H(chartData.Height))
	_buf.WriteString(`;">
</div>
<script>
    (function () {
        var chart = echarts.init(document.getElementById(`)
	_buf.WriteString(kmgXss.Jsonv(chartData.domId))
	_buf.WriteString(`));
        var option = `)
	_buf.WriteString(chartData.GetOptionString())
	_buf.WriteString(`;
        `)
	_buf.WriteString(chartData.JS)
	_buf.WriteString(`        chart.setOption(option);
    })();
</script>`)
	return _buf.String()
}
