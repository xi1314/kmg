package kmgBootstrap

import "github.com/bronze1man/kmg/kmgView"

type Wrap struct {
	//<html><head><title>
	Title string
	Head  kmgView.HtmlRendererList
	Body  kmgView.HtmlRendererList
}

func (w Wrap) HtmlRender() string {
	return tplWrap(w)
}

func NewWrap(title string, content ...kmgView.HtmlRenderer) *Wrap {
	body := kmgView.HtmlRendererList{}
	for _, r := range content {
		body = append(body, r)
	}
	wrap := &Wrap{
		Title: title,
		Head: kmgView.HtmlRendererList{
			GetJQueryCDN(),
			GetBootstrapCDN(),
			GetMomentJsCDN(),
		},
		Body: body,
	}
	return wrap
}

//带有 ECharts JS CDN 信息，ECharts JS 有 1MB 左右，不是图表的页面慎用
func NewWrapWithChart(title string, content ...kmgView.HtmlRenderer) *Wrap {
	w := NewWrap(title, content...)
	w.Head = append(w.Head, GetEChartsCDN())
	return w
}

func GetJQueryCDN() kmgView.HtmlRenderer {
	return kmgView.Html(`
	<script src="http://7xjyrb.com2.z0.glb.qiniucdn.com/jquery-1.11.3.min.js"></script>
	`)
}

func GetBootstrapCDN() kmgView.HtmlRenderer {
	return kmgView.Html(`
    <link rel="stylesheet" href="http://7xjyrb.com2.z0.glb.qiniucdn.com/bootstrap.min.css">
	`)
}

func GetEChartsCDN() kmgView.HtmlRenderer {
	return kmgView.Html(`
	<script src="http://7xjyrb.com2.z0.glb.qiniucdn.com/echarts-all.js"></script>
	`)
}

func GetMomentJsCDN() kmgView.HtmlRenderer {
	return kmgView.Html(`
	<script src="http://7xjyrb.com2.z0.glb.qiniucdn.com/moment.js"></script>
	`)
}
