package kmgConfig

import (
	"database/sql"
	"fmt"
	"github.com/bronze1man/kmg/encoding/kmgYaml"
	"github.com/bronze1man/kmg/kmgSql"
	"path/filepath"
	"sync"
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

	dbOnce sync.Once
	db     *kmgSql.Db
}

func (p *Parameter) GetDbConfig() *kmgSql.DbConfig {
	return &kmgSql.DbConfig{
		Username: p.DatabaseUsername,
		Password: p.DatabasePassword,
		Host:     p.DatabaseHost,
		DbName:   p.DatabaseDbName,
	}
}

//放错地方了?
func (p *Parameter) GetDb() (db *kmgSql.Db) {
	p.dbOnce.Do(func() {
		odb, err := sql.Open("mysql", p.GetDbConfig().GetDsn())
		if err != nil {
			panic(err)
		}
		p.db = &kmgSql.Db{
			DB: odb,
		}
	})
	return p.db
}
