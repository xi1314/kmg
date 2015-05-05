package kmgConfig

import (
	"os"
	"path/filepath"
	"strings"

	"fmt"
	"github.com/bronze1man/kmg/encoding/kmgYaml"
	"github.com/bronze1man/kmg/kmgFile"
	"sync"
)

//if you init it like &Context{xxx},please call Init()
//和目录相关的环境配置? .kmg.yml
type Env struct {
	GOPATH             []string
	CrossCompileTarget []CompileTarget
	//default to $ProjectPath/config
	ConfigPath string
	//default to $ProjectPath/data
	DataPath string
	//default to $ProjectPath/tmp
	TmpPath string
	//default to $ProjectPath/log
	LogPath string
	//should come from environment 此参数可以配置
	GOROOT string
	//the dir of ".kmg.yml" 此参数不能配置
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
		context.GOPATH[i] = kmgFile.FullPathOnPath(context.ProjectPath, p)
	}
	if context.GOROOT == "" {
		context.GOROOT = os.Getenv("GOROOT")
	}
	if context.DataPath == "" {
		context.DataPath = filepath.Join(context.ProjectPath, "data")
	}
	context.DataPath = kmgFile.FullPathOnPath(context.ProjectPath, context.DataPath)
	if context.TmpPath == "" {
		context.TmpPath = filepath.Join(context.ProjectPath, "tmp")
	}
	context.TmpPath = kmgFile.FullPathOnPath(context.ProjectPath, context.TmpPath)
	if context.ConfigPath == "" {
		context.ConfigPath = filepath.Join(context.ProjectPath, "config")
	}
	context.ConfigPath = kmgFile.FullPathOnPath(context.ProjectPath, context.ConfigPath)
	if context.LogPath == "" {
		context.LogPath = filepath.Join(context.ProjectPath, "log")
	}
	context.LogPath = kmgFile.FullPathOnPath(context.ProjectPath, context.LogPath)
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
func (context *Env) PathInConfig(relPath string) string {
	return filepath.Join(context.ConfigPath, relPath)
}
func (context *Env) PathInTmp(relPath string) string {
	return filepath.Join(context.TmpPath, relPath)
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
