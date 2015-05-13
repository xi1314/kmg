package SubCommand

import (
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgConsole"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"

	"github.com/bronze1man/kmg/kmgRand"
)

func selfUpdate() {
	prefixB, err := kmgHttp.UrlGetContent("http://kmgtools.qiniudn.com/v1/kmgUrlPrefix.txt?v=" + kmgRand.MustCryptoRandToAlphaNum(16))
	kmgConsole.ExitOnErr(err)

	exeContent, err := kmgHttp.UrlGetContent(string(prefixB) + "_windows_amd64.exe")
	kmgConsole.ExitOnErr(err)

	//cmd 这个东西有超级神力,直接os.Rename不行 但是360会报警
	// 已经试过下列方案:
	// 1.os.Rename 后面的write会没有权限,原因不明
	// 2.move windows上面没有这个命令
	kmgCmd.ProxyRun(`cmd /c move C:\windows\system32\kmg.exe C:\windows\system32\kmg-old.exe`)

	kmgFile.MustWriteFile(`C:\windows\system32\kmg.exe`, exeContent)
}
