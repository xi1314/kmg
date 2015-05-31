package InstallCmd

import (
	"github.com/bronze1man/kmg/kmgConsole"

	"github.com/bronze1man/kmg/kmgCmd"
	//"strings"
	"fmt"
	"github.com/bronze1man/kmg/kmgCompress"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	"github.com/bronze1man/kmg/kmgPlatform"
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
	p := kmgPlatform.GetCompiledPlatform()
	if p.Compatible(kmgPlatform.WindowsAmd64) {
		contentB, err := kmgHttp.UrlGetContent("http://kmgtools.qiniudn.com/v1/go1.4.2.windows-amd64.zip")
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
	kmgFile.MustChangeToTmpPath()
	if !kmgCmd.MustIsRoot() {
		fmt.Println("you need to be root to install golang")
		return
	}

	packageName := ""

	switch {
	case p.Compatible(kmgPlatform.LinuxAmd64):
		packageName = "go1.4.2.linux-amd64.tar.gz"
		kmgCmd.ProxyRun("apt-get install -y gcc")
	case p.Compatible(kmgPlatform.DarwinAmd64):
		packageName = "go1.4.2.darwin-amd64-osx10.8.tar.gz"
	default:
		kmgConsole.ExitOnErr(fmt.Errorf("not support platform [%s]", p))
	}
	contentB := kmgHttp.MustUrlGetContentProcess("http://kmgtools.qiniudn.com/v1/" + packageName)

	kmgFile.MustWriteFile(packageName, contentB)
	kmgCmd.ProxyRun("tar -xf " + packageName)
	kmgCmd.ProxyRun("cp -rf go /usr/local")
	kmgFile.MustDeleteFile("/bin/go")
	kmgCmd.ProxyRun("ln -s /usr/local/go/bin/go /bin/go")
	kmgFile.MustDeleteFile("/bin/godoc")
	kmgCmd.ProxyRun("ln -s /usr/local/go/bin/godoc /bin/godoc")
	kmgFile.MustDeleteFile("/bin/gofmt")
	kmgCmd.ProxyRun("ln -s /usr/local/go/bin/gofmt /bin/gofmt")
	kmgCmd.MustEnsureBinPath("/bin/go")
	kmgCmd.MustEnsureBinPath("/bin/godoc")
	kmgCmd.MustEnsureBinPath("/bin/gofmt")

}
