package kmgSql

import "database/sql"

//a wrap of database/sql.Db
type Db struct {
	*sql.DB
}

/*
func GetDbFromConfig(p *kmgConfig.Parameter) (db *Db) {
	odb, err := sql.Open("mysql", GetDbConfigFromConfig(p).GetDsn())
	if err != nil {
		panic(err)
	}
	return Db{DB:odb}
}
*/
