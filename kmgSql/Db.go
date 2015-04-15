package kmgSql

import (
	"database/sql"
	"github.com/bronze1man/kmg/kmgConfig/defaultParameter"
	_ "github.com/go-sql-driver/mysql"
	"sync"
)

//a wrap of database/sql.Db
type Db struct {
	*sql.DB
}

var dbOnce sync.Once
var db *Db

func GetDb() *Db {
	dbOnce.Do(func() {
		odb, err := sql.Open("mysql", GetDbConfigFromConfig(defaultParameter.Parameter()).GetDsn())
		if err != nil {
			panic(err)
		}
		db = &Db{DB: odb}
	})
	return db
}
