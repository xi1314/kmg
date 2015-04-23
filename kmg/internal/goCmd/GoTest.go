package goCommand

import (
	"flag"
	"fmt"
	//"github.com/bronze1man/kmg/console"
	"github.com/bronze1man/kmg/kmgFile"
	"go/build"
	"os"
	"path/filepath"

	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgConsole"
)

/*
递归目录的go test
 支持.kmg.yml目录结构提示文件(该文件必须存在)
 -v 更详细的描述
 -m 一个模块名,从这个模块名开始递归目录测试
 -d 一个目录名,从这个目录开始递归目录测试
 -bench benchmarks参数,直接传递到go test
 -onePackage 不递归目录测试,仅测试一个package
*/
func init() {
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "GoTest",
		Desc:   "递归目录的go test",
		Runner: runGoTest,
	})
}

func runGoTest() {
	command := GoTest{}
	flag.BoolVar(&command.v, "v", false, "show output of test")
	flag.StringVar(&command.dir, "d", "", "dir path to test")
	flag.StringVar(&command.moduleName, "m", "", "module name to test")
	flag.StringVar(&command.bench, "bench", "", "bench parameter pass to go test")
	flag.BoolVar(&command.onePackage, "onePackage", false, "only test one package")
	flag.StringVar(&command.runArg, "run", "", "Run only those tests and examples matching the regular expression.")
	flag.BoolVar(&command.onlyBuild, "onlyBuild", false, "only build all package(not test)")
	flag.Parse()

	kmgc, err := kmgConfig.LoadEnvFromWd()
	if err == nil {
		command.gopath = kmgc.GOPATH[0]
	} else {
		if kmgConfig.IsNotFound(err) {
			command.gopath = os.Getenv("GOPATH")
		} else {
			kmgConsole.ExitOnErr(err)
		}
	}
	//find root path
	root, err := command.findRootPath()
	kmgConsole.ExitOnErr(err)
	command.buildContext = &build.Context{
		GOPATH:   command.gopath,
		Compiler: build.Default.Compiler,
	}
	if command.onePackage {
		err = command.handlePath(root)
		kmgConsole.ExitOnErr(err)
	}
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}
		if kmgFile.IsDotFile(path) {
			return filepath.SkipDir
		}
		return command.handlePath(path)
	})
	kmgConsole.ExitOnErr(err)
}

/*
递归目录的 go test
 TODO 处理多个GOPATH的问题,从GOPATH里面找到这个模块
 支持.kmg.yml目录结构提示文件(该文件必须存在)
 -v 更详细的描述
 -m 一个模块名,从这个模块名开始递归目录测试
 -d 一个目录名,从这个目录开始递归目录测试
*/
type GoTest struct {
	gopath string
	//context      *console.Context
	v            bool
	dir          string
	moduleName   string
	bench        string
	onePackage   bool
	runArg       string
	buildContext *build.Context
	onlyBuild    bool
}

func (command *GoTest) findRootPath() (root string, err error) {
	if flag.NArg() == 1 {
		command.moduleName = flag.Arg(0)
	}
	if command.dir != "" {
		root = command.dir
		exist, err := kmgFile.FileExist(root)
		if err != nil {
			return "", err
		}
		if !exist {
			return "", fmt.Errorf("[GoTest] dir path:[%s] not found", root)
		}
		return root, nil
	}
	if command.moduleName != "" {
		//TODO 处理多个GOPATH的问题,从GOPATH里面找到这个模块
		root = filepath.Join(command.gopath, "src", command.moduleName)
		exist, err := kmgFile.FileExist(root)
		if err != nil {
			return "", err
		}
		if !exist {
			return "", fmt.Errorf("[GoTest] module name:[%s] not found", command.moduleName)
		}
		return root, nil
	}
	if root == "" {
		root, err = os.Getwd()
		if err != nil {
			return
		}
	}
	return
}

func (command *GoTest) handlePath(path string) error {
	pkg, err := command.buildContext.ImportDir(path, build.ImportMode(0))
	if err != nil {
		//仅忽略 不是golang的目录的错误
		_, ok := err.(*build.NoGoError)
		if ok {
			return nil
		}
		return err
	}
	if pkg.IsCommand() {
		return nil
	}
	if command.onlyBuild || len(pkg.TestGoFiles) == 0 {
		//如果没有测试文件,还会尝试build一下这个目录
		return command.gobuild(path)
		//return nil
	}
	return command.gotest(path)
}

func (command *GoTest) gotest(path string) error {
	fmt.Printf("[gotest] path[%s]\n", path)
	args := []string{"test"}
	if command.v {
		args = append(args, "-v")
	}
	if command.bench != "" {
		args = append(args, "-bench", command.bench)
	}
	if command.runArg != "" {
		args = append(args, "-run", command.runArg)
	}
	return kmgCmd.CmdSlice(append([]string{"go"}, args...)).
		MustSetEnv("GOPATH", command.gopath).
		SetDir(path).
		StdioRun()
}

func (command *GoTest) gobuild(path string) error {
	fmt.Printf("[gobuild] path[%s]\n", path)
	err := kmgCmd.CmdSlice([]string{"go", "build"}).
		MustSetEnv("GOPATH", command.gopath).
		SetDir(path).
		StdioRun()
	return err
}
