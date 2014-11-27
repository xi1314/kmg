package kmgConfig

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/bronze1man/kmg/encoding/kmgYaml"
)

//if you init it like &Context{xxx},please call Init()
//和目录相关的环境配置? .kmg.yml
type Env struct {
	GOPATH             []string
	CrossCompileTarget []CompileTarget

	//default to $ProjectPath/app
	AppPath string
	//default to $ProjectPath/config
	ConfigPath string
	//default to $AppPath/data
	DataPath string
	//default to $AppPath/tmp
	TmpPath string
	//default to ./app/log
	LogPath string

	//should come from environment
	GOROOT string
	//should come from dir of ".kmg.yml"
	ProjectPath string
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
	if context.AppPath == "" {
		context.AppPath = filepath.Join(context.ProjectPath, "app")
	}
	if context.DataPath == "" {
		context.DataPath = filepath.Join(context.AppPath, "data")
	}
	if context.TmpPath == "" {
		context.TmpPath = filepath.Join(context.AppPath, "tmp")
	}
	if context.ConfigPath == "" {
		context.ConfigPath = filepath.Join(context.AppPath, "config")
	}
	if context.LogPath == "" {
		context.LogPath = filepath.Join(context.AppPath, "log")
	}
	if len(context.GOPATH) == 0 {
		context.GOPATH = []string{context.ProjectPath}
	}
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
