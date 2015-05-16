package kmgSql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"sync"
)

//a wrap of database/sql.Db
type Db struct {
	*sql.DB
}

var dbLock sync.Mutex
var db *Db

func GetDb() *Db {
	dbLock.Lock()
	defer dbLock.Unlock()
	if db == nil {
		if defaultDbConfig == nil {
			panic("you need use SetDefaultDbConfig to set the database config")
		}
		odb, err := sql.Open("mysql", defaultDbConfig.GetDsn())
		if err != nil {
			panic(err)
		}
		db = &Db{DB: odb}
	}
	return db
}
