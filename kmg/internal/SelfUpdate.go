package internal

import (
	"errors"
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgConsole"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	"github.com/bronze1man/kmg/kmgPlatform"
	"github.com/bronze1man/kmg/kmgRand"
	"io/ioutil"
	"net/http"
	"os"
)

func init() {
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "SelfUpdate",
		Desc:   "update kmg tool from our server",
		Runner: selfUpdate,
	})
}

func selfUpdate() {
	response, err := http.Get("http://kmgtools.qiniudn.com/v1/installKmg.bash?v=" + kmgRand.MustCryptoRandToAlphaNum(5))
	kmgConsole.ExitOnErr(err)
	defer response.Body.Close()
	if kmgPlatform.GetCompiledPlatform().Compatible(kmgPlatform.WindowsAmd64) {
		if !kmgPlatform.GetCompiledPlatform().Compatible(kmgPlatform.WindowsAmd64) {
			kmgConsole.ExitOnErr(errors.New("this file should only run on windows_amd64"))
		}
		prefixB, err := kmgHttp.UrlGetContent("http://kmgtools.qiniudn.com/v1/kmgUrlPrefix.txt")
		kmgConsole.ExitOnErr(err)

		exeContent, err := kmgHttp.UrlGetContent(string(prefixB) + "_windows_amd64.exe")
		kmgConsole.ExitOnErr(err)

		//这个也还是不行,还需要继续研究方案.
		err = os.Rename(`C:\windows\system32\kmg.exe`, `C:\windows\system32\kmg.old.exe`)
		kmgConsole.ExitOnErr(err)

		kmgFile.MustWriteFile(`C:\windows\system32\kmg.exe`, exeContent)
	} else {
		baseFileContent, err := ioutil.ReadAll(response.Body)
		kmgConsole.ExitOnErr(err)

		baseFilePath := "/tmp/installKmg.bash"
		kmgFile.MustDeleteFile(baseFilePath)
		kmgFile.MustAppendFile(baseFilePath, baseFileContent)
		kmgCmd.MustRunInBash(baseFilePath)
	}
}
