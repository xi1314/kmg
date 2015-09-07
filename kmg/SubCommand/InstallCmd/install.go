package InstallCmd

import (
	"github.com/bronze1man/kmg/kmgConsole"

	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgSys"
	//"strings"
	"fmt"

	"github.com/bronze1man/kmg/kmgCompress"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	"github.com/bronze1man/kmg/kmgPlatform"
	"github.com/bronze1man/kmg/kmgTime"
	"path"
	"time"
)

func AddCommandList() {

	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "install",
		Desc:   "install tool",
		Runner: installCmd,
	})
}

func installCmd() {
	cg := kmgConsole.NewCommandGroup()
	cg.AddCommandWithName("golang", installGolang)
	cg.AddCommandWithName("golang1.5", installGolang15)
	cg.Main()
}

/*
本次实现遇到下列问题:
	* 不同操作系统判断
	* root权限判断
	* 如果同名 wget会修改下载组件的名称,使用临时文件夹处理
	* 如果已经安装 cp -rf go /usr/local/go 会再创建一个/usr/local/go/go 目录,而不是更新它
	* 如果已经在/usr/local/bin/go处放了一个执行,又在/bin/go处放了一个执行文件, /bin/go的版本不会被使用,使用函数专门判断这种情况,并且把多余的去除掉
	* 如果上一种情况发生,当前bash不能执行go version,因为当前bash有路径查询缓存(暂时无解了..)
	* TODO 国外服务器因为网络太卡而不能安装
*/
func installGolang() {
	installGolangWithUrlMap(map[string]string{
		"windows_amd64": "http://kmgtools.qiniudn.com/v1/go1.4.2.windows-amd64.zip",
		"linux_amd64":   "http://kmgtools.qiniudn.com/v1/go1.4.2.linux-amd64.tar.gz",
		"darwin_amd64":  "http://kmgtools.qiniudn.com/v1/go1.4.2.darwin-amd64-osx10.8.tar.gz",
	})
}

func installGolang15() {
	installGolangWithUrlMap(map[string]string{
		"windows_amd64": "https://storage.googleapis.com/golang/go1.5.windows-amd64.zip",
		"linux_amd64":   "https://storage.googleapis.com/golang/go1.5.linux-amd64.tar.gz",
		"darwin_amd64":  "https://storage.googleapis.com/golang/go1.5.darwin-amd64.tar.gz",
	})
}

func installGolangWithUrlMap(urlMap map[string]string) {
	p := kmgPlatform.GetCompiledPlatform()
	if p.Compatible(kmgPlatform.WindowsAmd64) {
		contentB, err := kmgHttp.UrlGetContent(urlMap["windows_amd64"])
		kmgConsole.ExitOnErr(err)
		err = kmgCompress.ZipUncompressFromBytesToDir(contentB, `c:\go`, "go")
		kmgConsole.ExitOnErr(err)
		err = kmgFile.CopyFile(`c:\go\bin\go.exe`, `c:\windows\system32\go.exe`)
		kmgConsole.ExitOnErr(err)
		err = kmgFile.CopyFile(`c:\go\bin\godoc.exe`, `c:\windows\system32\godoc.exe`)
		kmgConsole.ExitOnErr(err)
		err = kmgFile.CopyFile(`c:\go\bin\gofmt.exe`, `c:\windows\system32\gofmt.exe`)
		kmgConsole.ExitOnErr(err)
		return
	}
	tmpPath := kmgFile.MustChangeToTmpPath()
	defer kmgFile.MustDelete(tmpPath)
	if !kmgSys.MustIsRoot() {
		fmt.Println("you need to be root to install golang")
		return
	}

	url, ok := urlMap[p.String()]
	if !ok {
		kmgConsole.ExitOnErr(fmt.Errorf("not support platform [%s]", p))
	}
	packageName := path.Base(url)
	contentB := kmgHttp.MustUrlGetContentProcess(url)

	kmgFile.MustWriteFile(packageName, contentB)
	kmgCmd.ProxyRun("tar -xf " + packageName)
	if kmgFile.MustFileExist("/usr/local/go") {
		kmgCmd.ProxyRun("mv /usr/local/go /usr/local/go.bak." + time.Now().Format(kmgTime.FormatFileNameV2))
	}
	kmgCmd.ProxyRun("cp -rf go /usr/local")
	kmgFile.MustDeleteFile("/bin/go")
	kmgCmd.ProxyRun("ln -s /usr/local/go/bin/go /bin/go")
	kmgFile.MustDeleteFile("/bin/godoc")
	kmgCmd.ProxyRun("ln -s /usr/local/go/bin/godoc /bin/godoc")
	kmgFile.MustDeleteFile("/bin/gofmt")
	kmgCmd.ProxyRun("ln -s /usr/local/go/bin/gofmt /bin/gofmt")
	kmgFile.MustEnsureBinPath("/bin/go")
	kmgFile.MustEnsureBinPath("/bin/godoc")
	kmgFile.MustEnsureBinPath("/bin/gofmt")
}
