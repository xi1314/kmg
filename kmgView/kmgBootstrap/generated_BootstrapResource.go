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
			GeneratedJsFileName:  "ce818a784e8497b6d760b34985d0a3e5.js",
			GeneratedCssFileName: "edbcee494dbe2bebd093ca823faaa72f.css",
			GeneratedUrlPrefix:   "http://kmgtools.qiniudn.com/kmgBootstrap",
			RequestImportList:    []string{"github.com/bronze1man/kmg/kmgView/kmgWeb/bootstrap", "github.com/bronze1man/kmg/kmgView/kmgWeb/font-awesome", "github.com/bronze1man/kmg/kmgView/kmgWeb/jquery", "github.com/bronze1man/kmg/kmgView/kmgWeb/moment"},
		}
		Bootstrapgenerated.Init()
	})
	return Bootstrapgenerated
}
