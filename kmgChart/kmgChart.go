package kmgChart

import (
	"github.com/bronze1man/kmg/encoding/kmgJson"
	"github.com/bronze1man/kmg/kmgRand"
)

type Vector2 struct {
	X interface{}
	Y interface{}
}

type axis struct {
	Type string        `json:"type"` //category/value/time/log
	Data []interface{} `json:"data"`
}

type series struct {
	Type string        `json:"type"` // 支持： 'line'（折线图） | 'bar'（柱状图）尚未支持： 'scatter'（散点图） | 'k'（K线图） | 'pie'（饼图） | 'radar'（雷达图） | 'chord'（和弦图） | 'force'（力导向布局图） | 'map'（地图）
	Data []interface{} `json:"data"`
}

type title struct {
	Text    string `json:"text"`
	Subtext string `json:"subtext"`
}

type option struct {
	XAxis   *axis    `json:"xAxis"`
	YAxis   *axis    `json:"yAxis"`
	Series  []series `json:"series"`
	Title   *title   `json:"title"`
	Toolbox string   `json:"toolbox"`
}

type Line struct {
	domId  string
	Width  string //可选
	Height string //可选
	Option *option
}

func newLineBase() *Line {
	return &Line{
		domId:  "line-" + kmgRand.MustCryptoRandFromByteList(6, "abcdefghijklmnopqrstuvwxyz"),
		Width:  "100%",
		Height: "500px",
		Option: &option{
			XAxis: &axis{},
			YAxis: &axis{},
			Title: &title{},
		},
	}
}

//XY都是数值
func NewLineXYIsNumber(data []Vector2) *Line {
	l := newLineBase()
	l.Option.XAxis.Type = "value"
	l.Option.YAxis.Type = "value"
	l.Option.Series = []series{
		series{
			Type: "line",
			Data: []interface{}{},
		},
	}
	for _, v := range data {
		l.Option.Series[0].Data = append(l.Option.Series[0].Data, []interface{}{v.X, v.Y})
	}
	return l
}

//X 是任意东西，均匀的；Y 是数值,离散非均匀的
func NewLineXIsAnyYIsNumber(data []Vector2) *Line {
	l := newLineBase()
	l.Option.XAxis.Type = "category"
	l.Option.XAxis.Data = []interface{}{}
	l.Option.YAxis.Type = "value"
	l.Option.Series = []series{
		series{
			Type: "line",
			Data: []interface{}{},
		},
	}
	for _, v := range data {
		l.Option.XAxis.Data = append(l.Option.XAxis.Data, v.X)
		l.Option.Series[0].Data = append(l.Option.Series[0].Data, v.Y)
	}
	return l
}

func (l Line) GetOptionString() string {
	return kmgJson.MustMarshalToString(l.Option)
}

func (l Line) HtmlRender() string {
	return viewLine(l)
}
