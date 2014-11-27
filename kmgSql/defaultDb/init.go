package defaultDb

import (
	"database/sql"
	"sync"

	"github.com/bronze1man/kmg/kmgConfig/defaultParameter"
	"github.com/bronze1man/kmg/kmgSql"
)

var dbonce sync.Once
var db *kmgSql.Db

func GetDb() *kmgSql.Db {
	dbonce.Do(func() {
		odb, err := sql.Open("mysql", kmgSql.GetDbConfigFromConfig(defaultParameter.Parameter()).GetDsn())
		if err != nil {
			panic(err)
		}
		db = &kmgSql.Db{DB: odb}
	})
	return db
}
