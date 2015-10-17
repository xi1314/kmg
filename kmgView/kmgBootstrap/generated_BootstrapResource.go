package kmgBootstrap

import (
	"github.com/bronze1man/kmg/kmgView/kmgViewResource"
	"sync"
)

var BootstrapOnce sync.Once
var Bootstrapgenerated *kmgViewResource.Generated

func getBootstrapViewResource() *kmgViewResource.Generated {
	BootstrapOnce.Do(func() {
		Bootstrapgenerated = &kmgViewResource.Generated{
			Name:                 "Bootstrap",
			GeneratedJsFileName:  "0265f0f0d0695359fcd14680b7c7bb7b.js",
			GeneratedCssFileName: "ba62a1e2c3f272532a297e1f2340d471.css",
			GeneratedUrlPrefix:   "http://kmgtools.qiniudn.com/kmgBootstrap",
			RequestImportList:    []string{"github.com/bronze1man/kmg/kmgView/kmgBootstrap/WebResource"},
		}
		Bootstrapgenerated.Init()
	})
	return Bootstrapgenerated
}
