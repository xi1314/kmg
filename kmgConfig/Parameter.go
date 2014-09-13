package kmgConfig

import (
	"fmt"
	"github.com/bronze1man/kmg/encoding/kmgYaml"
	"github.com/bronze1man/kmg/kmgSql"
	"path/filepath"
)

var DefParameter *Parameter

func init() {
	DefParameter = &Parameter{}
	path := filepath.Join(DefEnv.ConfigPath, "parameters.yml")
	err := kmgYaml.ReadFile(path, DefParameter)
	if err != nil {
		panic(fmt.Errorf("can not get Parameters config,do you forget write a config file at %s ? err: %s", path, err))
	}
}

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

func (p *Parameter) GetDbConfig() *kmgSql.DbConfig {
	return &kmgSql.DbConfig{
		Username: p.DatabaseUsername,
		Password: p.DatabasePassword,
		Host:     p.DatabaseHost,
		DbName:   p.DatabaseDbName,
	}
}
