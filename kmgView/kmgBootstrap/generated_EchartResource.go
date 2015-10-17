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
			GeneratedJsFileName:  "8198e8a7f0938b622f51b0f434e3d50e.js",
			GeneratedCssFileName: "5590f633af977f14cf224771c6e2f57e.css",
			GeneratedUrlPrefix:   "http://kmgtools.qiniudn.com/kmgBootstrap",
			RequestImportList:    []string{"github.com/bronze1man/kmg/kmgView/kmgBootstrap/WebResource", "github.com/bronze1man/kmg/kmgView/kmgWeb/echart"},
		}
		Echartsgenerated.Init()
	})
	return Echartsgenerated
}
