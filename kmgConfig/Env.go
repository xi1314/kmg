package kmgConfig

import (
	"os"
	"path/filepath"
	"strings"

	"fmt"
	"github.com/bronze1man/kmg/encoding/kmgYaml"
	"sync"
)

//if you init it like &Context{xxx},please call Init()
//和目录相关的环境配置? .kmg.yml
type Env struct {
	GOPATH             []string
	CrossCompileTarget []CompileTarget
	//default to $ProjectPath/config
	ConfigPath string
	//default to $AppPath/data
	DataPath string
	//default to $AppPath/tmp
	TmpPath string
	//default to $AppPath/log
	LogPath string
	//should come from environment
	GOROOT string
	//should come from dir of ".kmg.yml"
	ProjectPath string
	//make command,使用kmg make可以运行这个命令
	Make string
	//当前是否是测试
	IsTest bool
	//Http 请求最大内存占用 默认100M
	HttpRequestMaxMemory int64
}

func (context *Env) GOPATHToString() string {
	if len(context.GOPATH) == 0 {
		return ""
	}
	return strings.Join(context.GOPATH, ":")
}
func (context *Env) Init() {
	for i, p := range context.GOPATH {
		if filepath.IsAbs(p) {
			continue
		}
		context.GOPATH[i] = filepath.Join(context.ProjectPath, p)
	}
	if context.GOROOT == "" {
		context.GOROOT = os.Getenv("GOROOT")
	}
	if context.DataPath == "" {
		context.DataPath = filepath.Join(context.ProjectPath, "data")
	}
	if context.TmpPath == "" {
		context.TmpPath = filepath.Join(context.ProjectPath, "tmp")
	}
	if context.ConfigPath == "" {
		context.ConfigPath = filepath.Join(context.ProjectPath, "config")
	}
	if context.LogPath == "" {
		context.LogPath = filepath.Join(context.ProjectPath, "log")
	}
	if len(context.GOPATH) == 0 {
		context.GOPATH = []string{context.ProjectPath}
	}
	if context.HttpRequestMaxMemory == 0 {
		context.HttpRequestMaxMemory = 100 << 20
	}
}
func (context *Env) PathInProject(relPath string) string {
	return filepath.Join(context.ProjectPath, relPath)
}
func FindFromPath(p string) (context *Env, err error) {
	p, err = filepath.Abs(p)
	if err != nil {
		return
	}
	var kmgFilePath string
	for {
		kmgFilePath = filepath.Join(p, ".kmg.yml")
		_, err = os.Stat(kmgFilePath)
		if err == nil {
			//found it
			break
		}
		if !os.IsNotExist(err) {
			return
		}
		thisP := filepath.Dir(p)
		if p == thisP {
			err = NotFoundError{}
			return
		}
		p = thisP
	}
	context = &Env{}
	err = kmgYaml.ReadFile(kmgFilePath, context)
	if err != nil {
		return
	}
	context.ProjectPath, err = filepath.Abs(filepath.Dir(kmgFilePath))
	if err != nil {
		return
	}
	context.Init()
	return
}

func LoadEnvFromWd() (context *Env, err error) {
	p, err := os.Getwd()
	if err != nil {
		return
	}
	return FindFromPath(p)
}

type NotFoundError struct {
}

func (e NotFoundError) Error() string {
	return "not found .kmg.yml in the project dir"
}
func IsNotFound(err error) (ok bool) {
	_, ok = err.(NotFoundError)
	return
}

var envOnce sync.Once
var env *Env

func DefaultEnv() *Env {
	envOnce.Do(func() {
		var err error
		env, err = LoadEnvFromWd()
		if err != nil {
			panic(fmt.Errorf("can not getEnv,do you forget create a .kmg.yml at project root? err: %s", err))
		}
	})
	return env
}
