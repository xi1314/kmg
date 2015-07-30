package kmgSql

import (
	"database/sql"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

//a wrap of database/sql.Db
//表示一个数据库
// 这个可以建立很多连接
type DB struct {
	DbQueryer
	*sql.DB
}

func (q DB) Query(query string, args ...string) (output []map[string]string, err error) {
	return q.DbQueryer.Query(query, args...)
}

func (q DB) Exec(query string, args ...string) (sql.Result, error) {
	return q.DbQueryer.Exec(query, args...)
}

func NewDb(db *sql.DB) DB {
	return DB{
		DbQueryer: DbQueryer{
			SqlQueryer: db,
		},
		DB: db,
	}
}

//表示一个数据库Queryer
// 可以拿去查数据库,
type DbQueryer struct {
	SqlQueryer
}

// 表示一个数据事务
// 这个只有一条连接
type Tx struct {
	DbQueryer
	*sql.Tx
}

func (q Tx) Query(query string, args ...string) (output []map[string]string, err error) {
	return q.DbQueryer.Query(query, args...)
}

func (q Tx) Exec(query string, args ...string) (sql.Result, error) {
	return q.DbQueryer.Exec(query, args...)
}

func NewTx(tx *sql.Tx) Tx {
	return Tx{
		DbQueryer: DbQueryer{
			SqlQueryer: tx,
		},
		Tx: tx,
	}
}

type SqlQueryer interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type SqlTxer interface {
	Commit() error
	Rollback() error
}

var dbLock sync.Mutex
var db DB

func GetDb() DB {
	dbLock.Lock()
	defer dbLock.Unlock()
	if db.DB == nil {
		if defaultDbConfig == nil {
			panic("you need use SetDefaultDbConfig to set the database config")
		}
		odb, err := sql.Open("mysql", defaultDbConfig.GetDsn())
		if err != nil {
			panic(err)
		}
		db = NewDb(odb)
	}
	return db
}
