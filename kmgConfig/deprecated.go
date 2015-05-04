package kmgConfig

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/bronze1man/kmg/encoding/kmgYaml"
)

// @deprecated
// 你应该在你自己的package里面管理你的数据库配置,并且把这个数据库配置注册到 kmgSql.DefaultDbConfig
type Parameter struct {
	DatabaseUsername   string
	DatabasePassword   string
	DatabaseHost       string
	DatabaseDbName     string
	DatabaseTestDbName string

	MemcacheHostList []string

	SessionPrefix     string
	SessionExpiration string
}

var parameterOnce sync.Once
var parameter *Parameter

// @deprecated
func DefaultParameter() *Parameter {
	parameterOnce.Do(func() {
		parameter = &Parameter{}
		path := filepath.Join(DefaultEnv().ConfigPath, "Parameters.yml")
		err := kmgYaml.ReadFile(path, parameter)
		if err != nil {
			panic(fmt.Errorf("can not get Parameters config,do you forget write a config file at %s ? err: %s", path, err))
		}
	})
	return parameter
}
