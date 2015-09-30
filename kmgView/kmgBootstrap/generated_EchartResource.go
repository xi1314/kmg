package kmgBootstrap

import (
	"github.com/bronze1man/kmg/kmgView/kmgViewResource"
	"sync"
)

var EchartsOnce sync.Once
var Echartsgenerated *kmgViewResource.Generated

func getEchartsViewResource() *kmgViewResource.Generated {
	EchartsOnce.Do(func() {
		Echartsgenerated = &kmgViewResource.Generated{
			Name:                 "Echarts",
			GeneratedJsFileName:  "df41c1d5c015f3fc194e7bc28381e89b.js",
			GeneratedCssFileName: "1694220fc74c1fd1b6ccf55ecef64aeb.css",
			GeneratedUrlPrefix:   "http://kmgtools.qiniudn.com/kmgBootstrap",
			RequestImportList:    []string{"github.com/bronze1man/kmg/kmgView/kmgWeb/bootstrap", "github.com/bronze1man/kmg/kmgView/kmgWeb/font-awesome", "github.com/bronze1man/kmg/kmgView/kmgWeb/jquery", "github.com/bronze1man/kmg/kmgView/kmgWeb/moment", "github.com/bronze1man/kmg/kmgView/kmgWeb/echart"},
		}
		Echartsgenerated.Init()
	})
	return Echartsgenerated
}
