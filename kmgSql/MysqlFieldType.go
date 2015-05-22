package kmgSql

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgStrconv"
	"strconv"
	"strings"
)

type MysqlField struct {
	Name string
	Type MysqlFieldType
}

type MysqlFieldType struct {
	DataType         MysqlDataType
	IsUnsigned       bool
	IsAutoIncrement  bool
	CharacterSetName string //utf8
	CollationName    string //utf8_bin
	Default          string
	StringLength     int
}

func (t1 MysqlFieldType) Equal(t2 MysqlFieldType) bool {
	return t1 == t2
}

func (t1 MysqlFieldType) String() string {
	out := string(t1.DataType)
	if t1.StringLength != 0 {
		out += "(" + strconv.Itoa(t1.StringLength) + ")"
	}
	if t1.IsUnsigned {
		out += " unsigned"
	}
	if t1.IsAutoIncrement {
		out += " auto_increment"
	}
	if t1.CollationName != "" {
		out += " COLLATE " + t1.CollationName
	}
	switch t1.DataType {
	case MysqlDataTypeInt32, MysqlDataTypeInt8, MysqlDataTypeFloat:
		out += " DEFAULT " + strconv.Itoa(kmgStrconv.AtoIDefault0(t1.Default))
	case MysqlDataTypeVarchar, MysqlDataTypeDateTime:
		out += " DEFAULT " + fmt.Sprintf("%#v", t1.Default) //TODO 正确的序列化方式
	}
	return out
}

type MysqlDataType string

const (
	MysqlDataTypeVarchar  MysqlDataType = `varchar`
	MysqlDataTypeInt32    MysqlDataType = `int`
	MysqlDataTypeLongText MysqlDataType = `longtext`
	MysqlDataTypeFloat    MysqlDataType = `float`
	MysqlDataTypeDateTime MysqlDataType = `datetime`
	MysqlDataTypeInt8     MysqlDataType = `tinyint`
	MysqlDataTypeLongBlob MysqlDataType = `longblob`
)

func mustMysqlGetTableFieldTypeList(TableName string) (out []MysqlField) {
	fieldRowList := MustQuery(`SELECT * FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?`, GetDefaultDbConfig().DbName, TableName)
	for _, row := range fieldRowList {
		field := MysqlFieldType{}
		field.DataType = MysqlDataType(row["DATA_TYPE"])
		field.Default = row["COLUMN_DEFAULT"]
		switch field.DataType {
		case MysqlDataTypeVarchar:
			field.CharacterSetName = row["CHARACTER_SET_NAME"]
			field.CollationName = row["COLLATION_NAME"]
			field.StringLength = kmgStrconv.AtoIDefault0(row["CHARACTER_MAXIMUM_LENGTH"])
		case MysqlDataTypeLongText:
			field.CharacterSetName = row["CHARACTER_SET_NAME"]
			field.CollationName = row["COLLATION_NAME"]
		case MysqlDataTypeInt32, MysqlDataTypeInt8:
			field.IsUnsigned = strings.Contains(row["COLUMN_TYPE"], "unsigned")
			field.IsAutoIncrement = strings.Contains(row["EXTRA"], "auto_increment")
		case MysqlDataTypeDateTime, MysqlDataTypeFloat, MysqlDataTypeLongBlob:
		default:
			panic(fmt.Errorf("TODO implement MysqlDataType %s", field.DataType))
		}
		out = append(out, MysqlField{
			Name: row["COLUMN_NAME"],
			Type: field,
		})
	}
	return out
}
