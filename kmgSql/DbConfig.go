package kmgSql

import (
	"fmt"

	"github.com/bronze1man/kmg/encoding/kmgYaml"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgFile"
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

func (config *DbConfig) GetDsnWithoutDbName() string {
	return fmt.Sprintf("%s:%s@tcp(%s:3306)/?charset=utf8&timeout=5s",
		config.Username,
		config.Password,
		config.Host)
}

var defaultDbConfig *DbConfig

func SetDefaultDbConfig(conf *DbConfig) {
	dbLock.Lock()
	defer dbLock.Unlock()
	if db.DB != nil {
		db.DB.Close()
		db = DB{}
	}
	defaultDbConfig = conf
}

func GetDefaultDbConfig() *DbConfig {
	dbLock.Lock()
	defer dbLock.Unlock()
	if defaultDbConfig == nil {
		panic("you need use SetDefaultDbConfig to set the config")
	}
	return defaultDbConfig
}

type TestDbConf struct {
	Db *DbConfig
}

func HasTestConfig() bool {
	return kmgFile.MustFileExist(kmgConfig.DefaultEnv().PathInConfig("Test.yml"))
}

func HasProdConfig() bool {
	return kmgFile.MustFileExist(kmgConfig.DefaultEnv().PathInConfig("Prod.yml"))
}
func MustLoadTestConfig() {
	mustLoadConfigByFilename("Test.yml")
}

func MustLoadProdConfig() {
	mustLoadConfigByFilename("Prod.yml")
}

func LoadConfigWithDbName(dbname string) {
	SetDefaultDbConfig(&DbConfig{
		Username: "root",
		Password: "",
		Host:     "127.0.0.1",
		DbName:   dbname,
	})
}

func mustLoadConfigByFilename(filename string) {
	conf := TestDbConf{}
	err := kmgYaml.ReadFile(kmgConfig.DefaultEnv().PathInConfig(filename), &conf)
	if err != nil {
		panic(err)
	}
	SetDefaultDbConfig(conf.Db)
}
