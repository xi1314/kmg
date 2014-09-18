package kmgDefault

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgYaml"
	"path/filepath"
)

var Env *kmgConfig.Env

func init() {
	var err error
	Env, err = kmgConfig.LoadEnvFromWd()
	if err != nil {
		panic(fmt.Errorf("can not getEnv,do you forget create a .kmg.yml at project root? err: %s", err))
	}
}

var Parameter *kmgConfig.Parameter

func init() {
	Parameter = &kmgConfig.Parameter{}
	path := filepath.Join(Env.ConfigPath, "parameters.yml")
	err := kmgYaml.ReadFile(path, Parameter)
	if err != nil {
		panic(fmt.Errorf("can not get Parameters config,do you forget write a config file at %s ? err: %s", path, err))
	}
}

/*
var dbOnce sync.Once
var db     *kmgSql.Db
*/
