package kmgSql

import (
	"github.com/bronze1man/kmg/kmgStrings"
	"strings"
)

type Table struct {
	Name       string
	FieldList  map[string]DbType
	PrimaryKey string
	UniqueKey  [][]string
	Null       []string
}

type DbType string

const (
	DbTypeInt        DbType = `int(11) DEFAULT 0`
	DbTypeString     DbType = `varchar(255) COLLATE utf8_bin DEFAULT ""`
	DbTypeLongString DbType = `longtext COLLATE utf8_bin DEFAULT ""`
	DbTypeFloat      DbType = `float default 0`
	DbTypeDatetime   DbType = `datetime DEFAULT "0000-00-00 00:00:00"`
	DbTypeBool       DbType = `tinyint(4) DEFAULT 0`
)

func MustSyncTable(tableConf Table) {
	if MustIsTableExist(tableConf.Name) {
		MustModifyTable(tableConf)
	} else {
		MustCreateTable(tableConf)
	}
}

func MustIsTableExist(tableName string) bool {
	ret := MustQueryOne("SHOW TABLE STATUS WHERE Name=?", tableName)
	if len(ret) <= 0 {
		return false
	} else {
		return true
	}
}

func MustCreateTable(tableConf Table) {
	sql := "CREATE TABLE IF NOT EXISTS `" + tableConf.Name + "` \n("
	sqlItemList := []string{}
	for fieldName, fieldType := range tableConf.FieldList {
		sqlField := "`" + fieldName + "` " + string(fieldType)
		if !kmgStrings.IsInSlice(tableConf.Null, fieldName) {
			sqlField += " NOT NULL"
		}
		if tableConf.PrimaryKey == fieldName {
			continue
		}
		sqlItemList = append(sqlItemList, sqlField)
	}
	if tableConf.PrimaryKey != "" {
		sqlItemList = append(sqlItemList, "`"+tableConf.PrimaryKey+"` int(11) unsigned AUTO_INCREMENT")
		sqlItemList = append(sqlItemList, "PRIMARY KEY (`"+tableConf.PrimaryKey+"`)")
	}
	for _, group := range tableConf.UniqueKey {
		uniqueSql := "UNIQUE INDEX ("
		uniqueKeyList := []string{}
		for _, key := range group {
			uniqueKeyList = append(uniqueKeyList, "`"+key+"`")
		}
		uniqueSql += strings.Join(uniqueKeyList, ",") + ")"
		sqlItemList = append(sqlItemList, uniqueSql)
	}
	sql += strings.Join(sqlItemList, ",\n")
	sql += "\n) engine=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin"
	MustExec(sql)
}

func MustModifyTable(tableConf Table) {
	fieldRow := MustQuery("SHOW COLUMNS FROM `" + tableConf.Name + "`")
	dbFieldNameList := []string{}
	for _, row := range fieldRow {
		dbFieldNameList = append(dbFieldNameList, row["Field"])
	}
	for fieldName, _ := range tableConf.FieldList {
		if kmgStrings.IsInSlice(dbFieldNameList, fieldName) {
			continue
		}
		MustAddNewField(tableConf, fieldName)
	}
}

func MustAddNewField(tableConf Table, newFieldName string) {
	newFieldType := tableConf.FieldList[newFieldName]
	sql := "ALTER TABLE `" + tableConf.Name + "` ADD `" + newFieldName + "` " + string(newFieldType)
	if !kmgStrings.IsInSlice(tableConf.Null, newFieldName) {
		sql += " NOT NULL"
	}
	MustExec(sql)
}
