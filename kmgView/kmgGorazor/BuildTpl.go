package kmgGorazor

import (
	"fmt"

	"github.com/bronze1man/kmg/encoding/kmgBase64"
	"github.com/bronze1man/kmg/kmgCache"
	"github.com/sipin/gorazor/gorazor"
)

func MustBuildTplWithPath(path string) {
	kmgCache.MustMd5FileChangeCache("tpl_"+kmgBase64.Base64EncodeStringToString(path), []string{path}, func() {
		fmt.Println("build tpl at " + path)
		err := gorazor.GenFolder(path, path, gorazor.Option{"NameNotChange": true})
		if err != nil {
			panic(err)
		}
	})
}
