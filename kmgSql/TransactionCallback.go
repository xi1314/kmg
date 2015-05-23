package kmgSql

import (
	"database/sql"
)

//使用默认数据库开启事务回调
func MustTransactionCallback(f func(tx Tx)) {
	db := GetDb()
	var err error
	var tx *sql.Tx
	hasFinish := false
	defer func() { //panic的时候处理
		if !hasFinish {
			tx.Rollback()
			//不用recover,让异常继续向上传递
		}
	}()
	tx, err = db.Begin()
	if err != nil {
		panic(err)
	}
	f(NewTx(tx))
	err = tx.Commit() //TODO commit失败怎么搞?
	if err != nil {
		panic(err)
	}
	hasFinish = true
	return
}
