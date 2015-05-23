package kmgSql

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/go-sql-driver/mysql"
)

// @deprecated
func GetDbConfigFromConfig(p *kmgConfig.Parameter) *DbConfig {
	return &DbConfig{
		Username: p.DatabaseUsername,
		Password: p.DatabasePassword,
		Host:     p.DatabaseHost,
		DbName:   p.DatabaseDbName,
	}
}

// @deprecated
type TransactionableDb interface {
	Begin() error
	Commit() error
	Rollback() error
}

// @deprecated
//transaction callback on beego.orm,but not depend on it
// 有可能会多次运行
func TransactionCallback(db TransactionableDb, f func() error) (err error) {
	for i := 0; i < 3; i++ {
		err = runTransaction(db, f)
		mysqlErr, ok := err.(*mysql.MySQLError)
		if ok && mysqlErr.Number == 1213 {
			//1213 错误可以重试解决
			continue
		}
		return err
	}
	return err
}

// @deprecated
func runTransaction(db TransactionableDb, f func() error) error {
	hasFinish := false
	defer func() { //panic的时候处理
		if !hasFinish {
			db.Rollback()
			//不用recover,让异常继续向上传递
		}
	}()
	err := db.Begin()
	if err != nil {
		return err
	}
	err = f()
	if err != nil {
		errR := db.Rollback()
		hasFinish = true
		if errR != nil {
			return fmt.Errorf("rollback fail:%s,origin fail:%s", errR.Error(), err.Error())
		}
		return err
	}
	err = db.Commit()
	if err != nil {
		return err
	}
	hasFinish = true
	return nil
}
