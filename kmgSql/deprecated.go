package kmgSql

import "github.com/bronze1man/kmg/kmgConfig"

// @deprecated
func GetDbConfigFromConfig(p *kmgConfig.Parameter) *DbConfig {
	return &DbConfig{
		Username: p.DatabaseUsername,
		Password: p.DatabasePassword,
		Host:     p.DatabaseHost,
		DbName:   p.DatabaseDbName,
	}
}
