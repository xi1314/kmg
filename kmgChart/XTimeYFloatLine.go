package kmgChart

import (
	"sort"
	"time"

	"github.com/bronze1man/kmg/kmgTime"
)

type TimeFloatPair struct {
	X time.Time
	Y float64
}

func CreateLineFromTimeFloatPair(inputList []TimeFloatPair) *Chart {
	line := newChartBaseConfig()
	line.Option.XAxis.Type = "time"
	line.Option.YAxis.Type = "value"
	line.Option.Series = []series{
		series{
			Type:          "line",
			ShowAllSymbol: true,
			Data:          []interface{}{},
			MarkPoint: &markPoint{
				Data: []*markPointData{},
			},
		},
	}
	for _, v := range inputList {
		line.Option.Series[0].Data = append(line.Option.Series[0].Data, []interface{}{v.X, v.Y})
	}
	line.JS = `
		function formatDate(data) { //date JavaScript Date 对象
			return moment(data).format("YYYY-MM-DD HH:mm:ss:SSS")
		}
		//坐标显示值域通过 DataZoom 自适应
		delete option.xAxis.min
		delete option.xAxis.max
		delete option.yAxis.min;
		delete option.yAxis.max;
        option.xAxis.axisLabel = {
            formatter : function (value) {
                return formatDate(value)
            }
        };
        option.tooltip.formatter = function (params) {
            return "X: " + formatDate(params.value[0])
                    + "<br /> Y: " + params.value[1]
        };
        //把 X 轴输入的 MySQL 时间字符串转成 JavaScript Date 对象
        var list = option.series[0].data;
        var len = list.length;
        for (var i=0;i < len;i++ ) {
            list[i][0] = new Date(list[i][0])
        }
        option.series[0].data = list;
	`
	return line
}

type TimeFloatPairSortByDESC []TimeFloatPair

func (l TimeFloatPairSortByDESC) Len() int {
	return len(l)
}
func (l TimeFloatPairSortByDESC) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
func (l TimeFloatPairSortByDESC) Less(i, j int) bool {
	return l[i].X.After(l[j].X)
}

type LODLevel struct {
	Range   time.Duration
	Density time.Duration //每个 Density 时间间隔，取该时间间隔内的平均值
}

func TimeContainInRange(startTime time.Time, timeRange time.Duration, t time.Time) bool {
	endTime := startTime.Add(-timeRange)
	var period kmgTime.Period
	if endTime.After(startTime) {
		period = kmgTime.MustNewPeriod(startTime, endTime)
	} else {
		period = kmgTime.MustNewPeriod(endTime, startTime)
	}
	return period.IsIn(t)
}

func LODForTimeFloatLine(input []TimeFloatPair, levelList []LODLevel) []TimeFloatPair {
	sort.Sort(TimeFloatPairSortByDESC(input))
	output := []TimeFloatPair{}
	startTime := input[0].X
	lastTime := startTime
	density := time.Duration(0)
	total := float64(0)
	num := 0
	for _, v := range input {
		num++
		total += v.Y
		inLevel := false
		for _, level := range levelList {
			if TimeContainInRange(startTime, level.Range, v.X) {
				density = level.Density
				inLevel = true
				break
			}
		}
		if !inLevel {
			break
		}
		if TimeContainInRange(lastTime, density, v.X) {
			continue
		}
		output = append(output, TimeFloatPair{
			X: v.X,
			Y: total / float64(num),
		})
		lastTime = v.X
		num = 0
		total = float64(0)
	}
	return output
}

// 平均时间分析
func AvgTimeFloatPair(input []TimeFloatPair, Density time.Duration) []TimeFloatPair {
	if len(input) == 0 {
		return nil
	}
	output := []TimeFloatPair{}
	lastTime := input[0].X
	thisTotal := float64(0)
	thisNum := 0
	for _, v := range input {
		thisNum++
		thisTotal += v.Y
		if v.X.Sub(lastTime) < Density {
			continue
		}
		output = append(output, TimeFloatPair{
			X: v.X,
			Y: thisTotal / float64(thisNum),
		})
		thisNum = 0
		thisTotal = 0
		lastTime = v.X
	}
	// 直接忽略掉最后几条数据
	return output
}

// 累计时间分析
func AccTimePerSecondFloatPair(input []TimeFloatPair, Density time.Duration) []TimeFloatPair {
	if len(input) == 0 {
		return nil
	}
	output := []TimeFloatPair{}
	lastTime := input[0].X
	thisTotal := float64(0)
	for _, v := range input {
		thisTotal += v.Y
		if v.X.Sub(lastTime) < Density {
			continue
		}
		output = append(output, TimeFloatPair{
			X: v.X,
			Y: thisTotal / (float64(v.X.Sub(lastTime)) / 1e9),
		})
		thisTotal = 0
		lastTime = v.X
	}
	// 直接忽略掉最后几条数据
	return output
}
