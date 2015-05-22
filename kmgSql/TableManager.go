package kmgSql

import (
	"database/sql"
	"fmt"
)

var registerTableList = []Table{}

//向系统中注册表,不允许重复注册
// 不支持并发调用
func MustRegisterTable(table Table) {
	for i := range registerTableList {
		if registerTableList[i].Name == table.Name {
			panic(fmt.Errorf("[MustRegisterTable] table name %s repeat", table.Name))
		}
	}
	registerTableList = append(registerTableList, table)
}

// 不支持并发调用
func ClearReisterTable() {
	registerTableList = []Table{}
}

//同步注册进去的表(只会增加字段,保证不掉数据,会使用fmt显示有哪些字段存在问题.
// 不支持并发调用
func MustSyncRegisterTable() {
	for i := range registerTableList {
		MustSyncTable(registerTableList[i])
	}
}

//强制同步注册进去的表,可能会掉数据,保证字段达到配置的样子
// 不支持并发调用
func MustForceSyncRegisterTable() {
	for i := range registerTableList {
		MustForceSyncTable(registerTableList[i])
	}
}

//创建数据库
// 不支持并发调用
func MustCreateDb() {
	conf := GetDefaultDbConfig()
	db, err := sql.Open("mysql", conf.GetDsnWithoutDbName())
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", conf.DbName))
	db.Close()
	if err != nil {
		panic(err)
	}
}

//同步默认数据库配置
func MustSyncDefaultDbConfig() {
	MustCreateDb()
	MustSyncRegisterTable()
}

func MustForceSyncDefaultDbConfig() {
	MustCreateDb()
	MustForceSyncRegisterTable()
}
