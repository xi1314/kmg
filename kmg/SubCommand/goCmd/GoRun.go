package goCmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bronze1man/kmg/encoding/kmgGob"
	"github.com/bronze1man/kmg/encoding/kmgJson"
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgConsole"
	"github.com/bronze1man/kmg/kmgCrypto"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgGoSource"
	"github.com/bronze1man/kmg/kmgPlatform"
	//"github.com/bronze1man/kmg/kmgDebug"
)

// go install bug
func GoRunCmd() {
	kmgc, err := kmgConfig.LoadEnvFromWd()
	kmgConsole.ExitOnErr(err)
	goPath := kmgc.GOPATHToString()

	//假设第一个是文件或者package名称,后面是传给命令行的参数
	if len(os.Args) < 2 {
		kmgConsole.ExitOnErr(fmt.Errorf("you need pass in running filename"))
		return
	}
	pathOrPkg := os.Args[1]
	_, err = os.Stat(pathOrPkg)
	switch {
	case os.IsNotExist(err): //package名称
		goRunPackageName(goPath, pathOrPkg)
		return
	case err != nil: //其他错误
		kmgConsole.ExitOnErr(err)
		return

	default: //文件或目录
		wd, err := os.Getwd()
		kmgConsole.ExitOnErr(err)
		if wd == filepath.Join(goPath, "src") {
			//用户在src下
			goRunPackageName(goPath, pathOrPkg)
			return
		}

		//  找出指向的这个文件的所有import的包,全部install一遍,再go run
		//靠谱实现这个东西的复杂度太高,目前已有的方案不能达到目标,暂时先使用go run
		// 如果有需要使用请把这个文件放到package里面,或者运行前删除pkg目录.
		// TODO 速度比较慢.
		//已经证实不行的方案:
		// 1.在临时目录建一个package,并且使用GOPATH指向那个临时目录,缓存会出现问题,并且效果和 go build -i 没有区别
		// 2.使用go build -i 效果和直接go run没有区别(缓存还是会出现问题)

		//找出这个文件所有的 import然后install 一遍
		importPathList, err := kmgGoSource.GetImportPathListFromFile(pathOrPkg)
		kmgConsole.ExitOnErr(err)
		for _, pkgPath := range importPathList {
			runCmdSliceWithGoPath(goPath, []string{"go", "install", pkgPath})
		}
		runCmdSliceWithGoPath(goPath, append([]string{"go", "run"}, os.Args[1:]...))
		return
	}
	kmgConsole.ExitOnErr(fmt.Errorf("unexpected run path"))
}

//不回显命令
func runCmdSliceWithGoPath(gopath string, cmdSlice []string) {
	err := kmgCmd.CmdSlice(cmdSlice).
		MustSetEnv("GOPATH", gopath).StdioRun()
	kmgConsole.ExitOnErr(err)
}

func goRunPackageName(goPath string, pathOrPkg string) {
	goRunInstall(goPath, pathOrPkg)
	//run
	outPath := filepath.Join(goPath, "bin", filepath.Base(pathOrPkg))
	runCmdSliceWithGoPath(goPath, append([]string{outPath}, os.Args[2:]...))
}

type gorunCacheInfo struct {
	PkgMap map[string]*gorunCachePkgInfo // key pkgname
}

type gorunCachePkgInfo struct {
	GoFileMap map[string]string
	PkgMd5    string
	IsMain    bool
	Name      string
}

func (info *gorunCachePkgInfo) getPkgBinPath(gopath string, platform string) string {
	if info.IsMain {
		return filepath.Join("bin", filepath.Base(info.Name))
	} else {
		return filepath.Join("pkg", platform, info.Name+".a")
	}
}

func goRunInstall(goPath string, pathOrPkg string) {
	//只能更新本GOPATH里面的pkg,不能更多多个GOPATH里面其他GOPATH的pkg缓存.
	// TODO 已知bug1 删除某个package里面的部分文件,然后由于引用到了旧的实现的代码,不会报错.删除pkg解决问题.
	// TODO 已知bug2 如果一个package先是main,然后build了一个东西,然后又改成了非main,再gorun会使用旧的缓存/bin/里面的缓存.
	ok := goRunInstallIsValidAndInvalidCache(goPath, pathOrPkg)
	if ok {
		//fmt.Println("use cache")
		return
	}
	//fmt.Println("not use cache")
	runCmdSliceWithGoPath(goPath, []string{"go", "install", pathOrPkg})
	// 填充缓存
	platform := kmgPlatform.GetCompiledPlatform().String()
	info := &gorunCacheInfo{
		PkgMap: map[string]*gorunCachePkgInfo{},
	}
	outputJson := kmgCmd.CmdSlice([]string{"go", "list", "-json", pathOrPkg}).
		MustSetEnv("GOPATH", goPath).MustCombinedOutput()
	listObj := &struct {
		Deps []string
		Name string
	}{}
	kmgJson.MustUnmarshal(outputJson, &listObj)
	if listObj.Name != "main" {
		fmt.Printf("run non main package %s\n", pathOrPkg)
		return
	}
	listObj.Deps = append(listObj.Deps, pathOrPkg)
	for _, pkgName := range listObj.Deps {
		srcpkgPath := filepath.Join(goPath, "src", pkgName)
		fileList, err := kmgFile.ReadDirFileOneLevel(srcpkgPath)
		if err != nil {
			if !os.IsNotExist(err) {
				panic(err)
			}
			// 没有找到pkg,可能是这个pkg在GOROOT出现过,此处暂时不管.
			continue
		}
		pkgInfo := &gorunCachePkgInfo{
			GoFileMap: map[string]string{},
		}
		for _, file := range fileList {
			ext := filepath.Ext(file)
			if ext == ".go" {
				pkgInfo.GoFileMap[file] = kmgCrypto.MustMd5File(filepath.Join(srcpkgPath, file))
			}
		}
		pkgInfo.IsMain = pkgName == pathOrPkg
		pkgInfo.Name = pkgName
		pkgBinPath := pkgInfo.getPkgBinPath(goPath, platform)
		pkgInfo.PkgMd5 = kmgCrypto.MustMd5File(pkgBinPath)
		info.PkgMap[pkgName] = pkgInfo
	}
	tmpPath := filepath.Join(goPath, "tmp", "gorun", kmgCrypto.Md5HexFromString(pathOrPkg+"_"+platform))
	kmgGob.MustWriteFile(tmpPath, info)
}

func goRunInstallIsValidAndInvalidCache(goPath string, pathOrPkg string) bool {
	info := &gorunCacheInfo{}
	// TODO 获取本次install的 平台名称
	platform := kmgPlatform.GetCompiledPlatform().String()
	tmpPath := filepath.Join(goPath, "tmp", "gorun", kmgCrypto.Md5HexFromString(pathOrPkg+"_"+platform))
	err := kmgGob.ReadFile(tmpPath, &info)
	if err != nil { //此处故意忽略错误 没有缓存文件 TODO 此处需要折腾其他东西吗?
		return false
	}
	//kmgDebug.Println(info)
	isValid := true
	for pkgName, pkgInfo := range info.PkgMap {
		pkgPath := pkgInfo.getPkgBinPath(goPath, platform)
		if !checkFileWithMd5(pkgPath, pkgInfo.PkgMd5) {
			kmgFile.MustDelete(pkgPath)
			isValid = false
			continue
		}
		srcpkgPath := filepath.Join(goPath, "src", pkgName)
		fileList, err := kmgFile.ReadDirFileOneLevel(srcpkgPath)
		if err != nil {
			if !os.IsNotExist(err) {
				panic(err)
			}
			kmgFile.MustDelete(pkgPath)
			isValid = false
			continue
		}
		isThisPkgValid := true
		for _, file := range fileList {
			ext := filepath.Ext(file)
			if ext == ".go" {
				// 多了一个文件
				if pkgInfo.GoFileMap[file] == "" {
					kmgFile.MustDelete(pkgPath)
					isValid = false
					isThisPkgValid = false
					break
				}
			}
		}
		if !isThisPkgValid {
			continue
		}
		for name, md5 := range pkgInfo.GoFileMap {
			goFilePath := filepath.Join(srcpkgPath, name)
			if !checkFileWithMd5(goFilePath, md5) {
				kmgFile.MustDelete(pkgPath)
				isValid = false
				break
			}
		}
	}
	return isValid
}

func checkFileWithMd5(path string, shouldMd5 string) (ok bool) {
	content, err := kmgFile.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}
		return false
	}
	return kmgCrypto.Md5Hex(content) == shouldMd5
}
