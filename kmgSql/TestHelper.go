package kmgSql

import (
	"github.com/bronze1man/kmg/encoding/kmgYaml"
	"github.com/bronze1man/kmg/kmgConfig"
)

type TestDbConf struct {
	Db *DbConfig
}

func mustGetTestConfig() TestDbConf {
	conf := TestDbConf{}
	err := kmgYaml.ReadFile(kmgConfig.DefaultEnv().PathInConfig("Test.yml"), &conf)
	if err != nil {
		panic(err)
	}
	return conf
}

func MustLoadTestConfig() {
	SetDefaultDbConfig(mustGetTestConfig().Db)
}
