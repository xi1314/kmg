package kmgSql

import (
	"fmt"

	"github.com/bronze1man/kmg/kmgConfig"
)

type DbConfig struct {
	Username string
	Password string
	Host     string
	DbName   string
}

func (config *DbConfig) GetDsn() string {
	return fmt.Sprintf("%s:%s@%s/%s?charset=utf8&timeout=5s",
		config.Username,
		config.Password,
		config.Host,
		config.DbName)
}

func GetDbConfigFromConfig(p *kmgConfig.Parameter) *DbConfig {
	return &DbConfig{
		Username: p.DatabaseUsername,
		Password: p.DatabasePassword,
		Host:     p.DatabaseHost,
		DbName:   p.DatabaseDbName,
	}
}
