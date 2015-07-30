package kmgChart

import (
	"github.com/bronze1man/kmg/encoding/kmgJson"
	"github.com/bronze1man/kmg/kmgRand"
)

type Vector2 struct {
	X interface{}
	Y interface{}
}

type dataZoom struct {
	Show     bool        `json:"show"`
	RealTime bool        `json:"realtime"`
	Start    interface{} `json:"40"`
	End      interface{} `json:"end"`
}

type axis struct {
	Type string        `json:"type"` //category/value/time/log
	Data []interface{} `json:"data"`
	Min  interface{}   `json:"min,omitempty"`
	Max  interface{}   `json:"max,omitempty"`
}

type markPointData struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"` // max | min | average
	Value interface{} `json:"value"`
	X     interface{} `json:"x"`
	Y     interface{} `json:"y"`
}

type markPoint struct {
	Data  []*markPointData `json:"data"`
	Large bool             `json:"large"`
}

type series struct {
	Type          string        `json:"type"` // 支持： 'line'（折线图） | 'bar'（柱状图）尚未支持： 'scatter'（散点图） | 'k'（K线图） | 'pie'（饼图） | 'radar'（雷达图） | 'chord'（和弦图） | 'force'（力导向布局图） | 'map'（地图）
	ShowAllSymbol bool          `json:"showAllSymbol"`
	Data          []interface{} `json:"data"`
	MarkPoint     *markPoint    `json:"markPoint"`
}

type Title struct {
	Text    string `json:"text"`
	Subtext string `json:"subtext"`
}

type axisPointer struct {
	Show bool   `json:"show"`
	Type string `json:"type"` //辅助线一类的东西，可有可无 'line' | 'cross' | 'shadow' | 'none'
}

type tooltip struct {
	Trigger     string       `json:"trigger"` // 'item' | 'axis'
	AxisPointer *axisPointer `json:"axisPointer"`
}

type option struct {
	XAxis     *axis     `json:"xAxis"`
	YAxis     *axis     `json:"yAxis"`
	Series    []series  `json:"series"`
	Title     *Title    `json:"title"`
	Animation bool      `json:"animation"`
	Tooltip   *tooltip  `json:"tooltip"`
	DataZoom  *dataZoom `json:"dataZoom"`
}

type Chart struct {
	domId  string
	Width  string  //可选
	Height string  //可选
	Option *option `json:"options"`
	JS     string  //在 option 赋值后，渲染数据前调用
}

func newChartBaseConfig() *Chart {
	return &Chart{
		domId:  "chart-" + kmgRand.MustCryptoRandFromByteList(6, "abcdefghijklmnopqrstuvwxyz"),
		Width:  "100%",
		Height: "500px",
		Option: &option{
			XAxis: &axis{},
			YAxis: &axis{},
			Title: &Title{},
			Tooltip: &tooltip{
				Trigger: "axis",
				AxisPointer: &axisPointer{
					Show: true,
					Type: "none",
				},
			},
		},
	}
}

func (l Chart) GetOptionString() string {
	return kmgJson.MustMarshalIndentToString(l.Option)
}

func (l Chart) HtmlRender() string {
	return tplChart(l)
}

////X 是任意东西，均匀分布，类似枚举；Y 是数值,离散非均匀的
//func NewLineXIsAnyYIsNumber(data []Vector2) *Chart {
//	l := newChartBaseConfig()
//	l.Option.XAxis.Type = "category"
//	l.Option.XAxis.Data = []interface{}{}
//	l.Option.YAxis.Type = "value"
//	l.Option.Series = []series{
//		series{
//			Type: "line",
//			Data: []interface{}{},
//		},
//	}
//	for _, v := range data {
//		l.Option.XAxis.Data = append(l.Option.XAxis.Data, v.X)
//		l.Option.Series[0].Data = append(l.Option.Series[0].Data, v.Y)
//	}
//	return l
//}
