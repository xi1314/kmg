package beegoBundle

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/bronze1man/kmg/console"
	_ "github.com/go-sql-driver/mysql"
	"github.com/bronze1man/kmg/kmgConfig/defaultParameter"
	"github.com/bronze1man/kmg/kmgSql"

)

type BeegoOrmCreateDbCommand struct {
	env string
}

func (command *BeegoOrmCreateDbCommand) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{
		Name:  "BeegoOrmCreateDb",
		Short: "beego orm create db",
	}
}
func (command *BeegoOrmCreateDbCommand) ConfigFlagSet(flag *flag.FlagSet) {
	flag.StringVar(&command.env, "env", "dev", "database env(dev,test)")
}
func (command *BeegoOrmCreateDbCommand) Execute(context *console.Context) (err error) {
	//work around for container bug
	DbConfig := kmgSql.GetDbConfigFromConfig(defaultParameter.Parameter)

	dsn := fmt.Sprintf("%s:%s@%s/?charset=utf8&timeout=5s",
		DbConfig.Username,
		DbConfig.Password,
		DbConfig.Host)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return
	}
	_, err = db.Exec(fmt.Sprintf("create database %s", DbConfig.DbName))
	if err != nil {
		return
	}
	return
}
