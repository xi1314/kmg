package kmgSql

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgStrings"
	"strings"
	//"github.com/bronze1man/kmg/kmgDebug"
)

//读取数据库的表字段,不区分大小写(某些系统的mysql不区分大小写)
//写入数据库的表字段,区分大小写
type Table struct {
	Name       string
	FieldList  map[string]DbType
	PrimaryKey string
	UniqueKey  [][]string
	Null       []string
}

type DbType string

const (
	DbTypeInt              DbType = `int(11) DEFAULT 0`
	DbTypeIntAutoIncrement DbType = `int(11) unsigned AUTO_INCREMENT`
	DbTypeString           DbType = `varchar(255) COLLATE utf8_bin DEFAULT ""`
	DbTypeLongString       DbType = `longtext COLLATE utf8_bin DEFAULT ""`
	DbTypeFloat            DbType = `float default 0`
	DbTypeDatetime         DbType = `datetime DEFAULT "0000-00-00 00:00:00"`
	DbTypeBool             DbType = `tinyint(4) DEFAULT 0`
	DbTypeLongBlob         DbType = `LONGBLOB`
)

func (t DbType) GetMysqlFieldType() MysqlFieldType {
	switch t {
	case DbTypeInt:
		return MysqlFieldType{
			DataType: MysqlDataTypeInt32,
			Default:  "0",
		}
	case DbTypeIntAutoIncrement:
		return MysqlFieldType{
			DataType:        MysqlDataTypeInt32,
			IsUnsigned:      true,
			IsAutoIncrement: true,
		}
	case DbTypeString:
		return MysqlFieldType{
			DataType:         MysqlDataTypeVarchar,
			Default:          "",
			CharacterSetName: "utf8",
			CollationName:    "utf8_bin",
			StringLength:     255,
		}
	case DbTypeLongString:
		return MysqlFieldType{
			DataType:         MysqlDataTypeLongText,
			Default:          "",
			CharacterSetName: "utf8",
			CollationName:    "utf8_bin",
		}
	case DbTypeFloat:
		return MysqlFieldType{
			DataType: MysqlDataTypeFloat,
			Default:  "0",
		}
	case DbTypeDatetime:
		return MysqlFieldType{
			DataType: MysqlDataTypeDateTime,
			Default:  "0000-00-00 00:00:00",
		}
	case DbTypeBool:
		return MysqlFieldType{
			DataType: MysqlDataTypeInt8,
			Default:  "0",
		}
	case DbTypeLongBlob:
		return MysqlFieldType{
			DataType: MysqlDataTypeLongBlob,
		}
	default:
		panic(fmt.Errorf("Unsupport DbType %s", t))
	}
}

func MustSyncTable(tableConf Table) {
	MustVerifyTableConfig(tableConf)
	if MustIsTableExist(tableConf.Name) {
		MustModifyTable(tableConf)
	} else {
		MustCreateTable(tableConf)
	}
}

func MustForceSyncTable(tableConf Table) {
	MustVerifyTableConfig(tableConf)
	if MustIsTableExist(tableConf.Name) {
		MustForceModifyTable(tableConf)
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

func MustVerifyTableConfig(tableConf Table) {
	fieldNameFieldMap := map[string]bool{}
	for name := range tableConf.FieldList {
		name := strings.ToLower(name)
		if fieldNameFieldMap[name] {
			panic(fmt.Errorf("[MustVerifyTableConfig] Table[%s] Field[%s] 两个字段名只有大小写不一致",
				tableConf.Name, name))
		}
		fieldNameFieldMap[name] = true
	}
}

func MustCreateTable(tableConf Table) {
	sql := "CREATE TABLE IF NOT EXISTS `" + tableConf.Name + "` \n("
	sqlItemList := []string{}
	hasPrimaryKey := false
	for fieldName, fieldType := range tableConf.FieldList {
		if tableConf.PrimaryKey == fieldName {
			hasPrimaryKey = true
			//continue
		}
		sqlField := "`" + fieldName + "` " + string(fieldType)
		if !kmgStrings.IsInSlice(tableConf.Null, fieldName) {
			sqlField += " NOT NULL"
		}
		sqlItemList = append(sqlItemList, sqlField)
	}
	if tableConf.PrimaryKey != "" {
		if !hasPrimaryKey {
			panic(fmt.Sprintf(`tableConf.PrimaryKey[%s], 但是这个主键不在字段列表里面`, tableConf.PrimaryKey))
		}
		//sqlItemList = append(sqlItemList, "`"+tableConf.PrimaryKey+"` int(11) unsigned AUTO_INCREMENT")
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
	MysqlFieldTypeList := mustMysqlGetTableFieldTypeList(tableConf.Name)
	dbFieldNameList := []string{}
	for _, row := range MysqlFieldTypeList {
		dbFieldNameList = append(dbFieldNameList, strings.ToLower(row.Name))
	}
	for _, f1 := range dbFieldNameList {
		found := false
		for f2 := range tableConf.FieldList {
			if strings.EqualFold(f2, f1) {
				found = true
				break
			}
		}
		if !found {
			fmt.Printf("[kmgSql.SyncTable] 数据库中包含多余字段 Table[%s] Field[%s]\n", tableConf.Name, f1)
		}
	}
	for fieldName, fieldType := range tableConf.FieldList {
		if !kmgStrings.IsInSlice(dbFieldNameList, strings.ToLower(fieldName)) {
			MustAddNewField(tableConf, fieldName)
			continue
		}
		for _, row := range MysqlFieldTypeList {
			if row.Name == fieldName {
				if !fieldType.GetMysqlFieldType().Equal(row.Type) {
					fmt.Printf("[kmgSql.SyncTable] Table[%s] Field[%s] OldType[%s] NewType[%s] 数据库字段类型不一致\n",
						tableConf.Name, fieldName, row.Type.String(), fieldType.GetMysqlFieldType().String())
				}
				break
			}
			if strings.EqualFold(row.Name, fieldName) {
				fmt.Printf("[kmgSql.SyncTable] Table[%s] OldField[%s] NewField[%s] 数据库字段大小写不一致\n",
					tableConf.Name, fieldName, row.Name)
				break
			}
		}
	}
}

func MustForceModifyTable(tableConf Table) {
	MysqlFieldTypeList := mustMysqlGetTableFieldTypeList(tableConf.Name)
	dbFieldNameList := []string{}
	for _, row := range MysqlFieldTypeList {
		dbFieldNameList = append(dbFieldNameList, row.Name)
	}
	for _, f1 := range dbFieldNameList {
		found := false
		for f2 := range tableConf.FieldList {
			if f2 == f1 {
				found = true
				break
			}
		}
		if !found {
			MustExec(fmt.Sprintf("ALTER TABLE `%s` DROP COLUMN `%s`", tableConf.Name, f1))
		}
	}
	for fieldName, fieldType := range tableConf.FieldList {
		if kmgStrings.IsInSlice(dbFieldNameList, fieldName) {
			for _, row := range MysqlFieldTypeList {
				if row.Name == fieldName {
					if !fieldType.GetMysqlFieldType().Equal(row.Type) {
						MustExec(fmt.Sprintf("ALTER TABLE `%s` CHANGE COLUMN `%s` `%s` %s NOT NULL", tableConf.Name, fieldName, fieldName, fieldType))
					}
					break
				}
			}
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
