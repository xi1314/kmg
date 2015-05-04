package beegoBundle

import (
	"database/sql"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/bronze1man/kmg/kmgConfig/defaultParameter"
	"github.com/bronze1man/kmg/kmgConsole"
	"github.com/bronze1man/kmg/kmgSql"
	"os"
)

func AddCommandList() {
	kmgConsole.AddCommandWithName("BeegoOrmCreateDb", createDbCmd)
	kmgConsole.AddCommandWithName("BeegoOrmSyncDb", syncDbCmd)
}

func createDbCmd() {
	InitOrm()
	//work around for container bug
	DbConfig := kmgSql.GetDbConfigFromConfig(defaultParameter.Parameter())

	dsn := fmt.Sprintf("%s:%s@%s/?charset=utf8&timeout=5s",
		DbConfig.Username,
		DbConfig.Password,
		DbConfig.Host)
	db, err := sql.Open("mysql", dsn)
	kmgConsole.ExitOnErr(err)
	_, err = db.Exec(fmt.Sprintf("create database %s", DbConfig.DbName))
	kmgConsole.ExitOnErr(err)
}

func syncDbCmd() {
	InitOrm()
	//TODO register database config stuff.
	os.Args = []string{
		os.Args[0], "orm", "syncdb",
	}
	orm.RunCommand()
}
