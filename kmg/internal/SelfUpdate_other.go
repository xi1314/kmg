// +build !windows

package internal

import (
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgConsole"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	"github.com/bronze1man/kmg/kmgRand"
)

func selfUpdate() {
	baseFileContent, err := kmgHttp.UrlGetContent("http://kmgtools.qiniudn.com/v1/installKmg.bash?v=" + kmgRand.MustCryptoRandToAlphaNum(16))
	kmgConsole.ExitOnErr(err)

	baseFilePath := "/tmp/installKmg.bash"
	kmgFile.MustDeleteFile(baseFilePath)
	kmgFile.MustAppendFile(baseFilePath, baseFileContent)
	kmgCmd.MustRunInBash(baseFilePath)
}
