package kmgSql

import (
	"database/sql"
	"github.com/bronze1man/kmg/kmgConfig"
)

//a wrap of database/sql.Db
type Db struct {
	*sql.DB
}

func GetDbFromConfig(p *kmgConfig.Parameter) (db *Db) {
	p.dbOnce.Do(func() {
		odb, err := sql.Open("mysql", p.GetDbConfig().GetDsn())
		if err != nil {
			panic(err)
		}
		p.db = &Db{
			DB: odb,
		}
	})
	return p.db
}
