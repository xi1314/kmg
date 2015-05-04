package kmgSql

import (
	"fmt"
)

type DbConfig struct {
	Username string // example: root
	Password string // example: password
	Host     string // example: 127.0.0.1
	DbName   string // example: kmg_test
}

func (config *DbConfig) GetDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8&timeout=5s",
		config.Username,
		config.Password,
		config.Host,
		config.DbName)
}

func (config *DbConfig) GetDsnWithoutDbname() string {
	return fmt.Sprintf("%s:%s@tcp(%s:3306)/?charset=utf8&timeout=5s",
		config.Username,
		config.Password,
		config.Host)
}

var defaultDbConfig *DbConfig

func SetDefaultDbConfig(conf *DbConfig) {
	dbLock.Lock()
	defer dbLock.Unlock()
	if db != nil {
		db.DB.Close()
		db = nil
	}
	defaultDbConfig = conf
}

func GetDefaultDbConfig() *DbConfig {
	dbLock.Lock()
	defer dbLock.Unlock()
	return defaultDbConfig
}
