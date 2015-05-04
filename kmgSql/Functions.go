package kmgSql

import (
	"database/sql"
	"fmt"
	"github.com/bronze1man/kmg/kmgSql/MysqlAst"
	"strconv"
	"strings"
)

func MustQuery(query string, args ...string) []map[string]string {
	out, err := Query(query, args...)
	if err != nil {
		panic(err)
	}
	return out
}

func Query(query string, args ...string) (output []map[string]string, err error) {
	rows, err := GetDb().Query(query, argsStringToInterface(args...)...)
	if err != nil {
		return nil, fmt.Errorf("[Query] sql: [%s] data: [%s] err:[%s]", query, strings.Join(args, ","), err.Error())
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	lenColumn := len(columns)
	for rows.Next() {
		rowArray := make([]interface{}, lenColumn)
		//box value with *RawByte
		for k1 := range rowArray {
			var s sql.RawBytes
			rowArray[k1] = &s
		}
		if err := rows.Scan(rowArray...); err != nil {
			return nil, err
		}
		rowMap := make(map[string]string)
		for rowIndex, rowName := range columns {
			//unbox value with *string
			rowMap[rowName] = string(*(rowArray[rowIndex].(*sql.RawBytes)))
		}
		output = append(output, rowMap)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return
}

func MustQueryOne(query string, args ...string) map[string]string {
	out, err := QueryOne(query, args...)
	if err != nil {
		panic(err)
	}
	return out
}

func QueryOne(query string, args ...string) (output map[string]string, err error) {
	list, err := Query(query, args...)
	if err != nil {
		return nil, err
	}
	if len(list) <= 0 {
		return nil, err
	}
	output = list[0]
	return output, err
}

func MustExec(query string, args ...string) {
	_, err := Exec(query, args...)
	if err != nil {
		panic(err)
	}
}

func Exec(query string, args ...string) (sql.Result, error) {
	ret, err := GetDb().Exec(query, argsStringToInterface(args...)...)
	if err != nil {
		return nil, fmt.Errorf("[Exec] sql: [%s] data: [%s] err:[%s]", query, strings.Join(args, ","), err.Error())
	}
	return ret, err
}

func MustInsert(tableName string, row map[string]string) (lastInsertId int) {
	lastInsertId, err := Insert(tableName, row)
	if err != nil {
		panic(err)
	}
	return lastInsertId
}

func Insert(tableName string, row map[string]string) (lastInsertId int, err error) {
	keyList := []string{}
	valueList := []string{}
	for key, value := range row {
		keyList = append(keyList, key)
		valueList = append(valueList, value)
	}
	keyStr := "`" + strings.Join(keyList, "`,`") + "`"
	valueStr := strings.Repeat("?,", (len(row)-1)) + "?"
	sql := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)", tableName, keyStr, valueStr)
	result, err := Exec(sql, valueList...)
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	lastInsertId = int(id)
	return lastInsertId, err
}

func MustUpdateById(tableName string, primaryKeyName string, row map[string]string) {
	err := UpdateById(tableName, primaryKeyName, row)
	if err != nil {
		panic(err)
	}
}

func UpdateById(tableName string, primaryKeyName string, row map[string]string) error {
	keyList := []string{}
	valueList := []string{}
	var idValue string
	for key, value := range row {
		if primaryKeyName == key {
			idValue = value
			continue
		}
		keyList = append(keyList, "`"+key+"`=?")
		valueList = append(valueList, value)
	}
	if idValue == "" {
		return fmt.Errorf("%s no set", primaryKeyName)
	}
	valueList = append(valueList, idValue)
	updateStr := strings.Join(keyList, ",")
	//sql例子 UPDATE AdminUser SET username=?,password=? where id = 1;
	sql := fmt.Sprintf("UPDATE `%s` SET %s where `%s` = ?", tableName, updateStr, primaryKeyName)
	_, err := Exec(sql, valueList...)
	if err != nil {
		return err
	}
	return nil
}

func MustReplaceById(tableName string, primaryKeyName string, row map[string]string) (lastInsertId int) {
	lastInsertId, err := ReplaceById(tableName, primaryKeyName, row)
	if err != nil {
		panic(err)
	}
	return lastInsertId
}

func ReplaceById(tableName string, primaryKeyName string, row map[string]string) (lastInsertId int, err error) {
	var one map[string]string
	if idValue, ok := row[primaryKeyName]; ok {
		one, _ = GetOneWhere(tableName, primaryKeyName, idValue)
	}
	if one == nil {
		return Insert(tableName, row)
	}
	err = UpdateById(tableName, primaryKeyName, row)
	lastInsertId, err = strconv.Atoi(one[primaryKeyName])
	if err != nil {
		lastInsertId = 0
	}
	return lastInsertId, err
}

func MustGetOneWhere(tableName string, fieldName string, value string) (output map[string]string) {
	output, err := GetOneWhere(tableName, fieldName, value)
	if err != nil {
		panic(err)
	}
	return output
}

func GetOneWhere(tableName string, fieldName string, value string) (output map[string]string, err error) {
	sql := fmt.Sprintf("SELECT * FROM `%s` WHERE `%s`=? LIMIT 1", tableName, fieldName)
	return QueryOne(sql, value)
}

func MustDeleteById(tableName string, fieldName string, value string) {
	err := DeleteById(tableName, fieldName, value)
	if err != nil {
		panic(err)
	}
}

func DeleteById(tableName string, fieldName string, value string) error {
	sql := fmt.Sprintf("DELETE FROM `%s` WHERE `%s`=?", tableName, fieldName)
	_, err := Exec(sql, value)
	return err
}

func MustGetAllInTable(tableName string) (output []map[string]string) {
	output, err := GetAllInTable(tableName)
	if err != nil {
		panic(err)
	}
	return output
}

func GetAllInTable(tableName string) (output []map[string]string, err error) {
	output, err = Query("SELECT * FROM `" + tableName + "`")
	return output, err
}

func MustRunSelectCommand(selectCommand *MysqlAst.SelectCommand) (mapValue []map[string]string) {
	output, paramList := selectCommand.GetPrepareParameter()
	list, error := Query(output, paramList...)
	if error != nil {
		panic(error)
	}
	return list
}

func Ping() (err error) {
	return GetDb().Ping()
}

func argsStringToInterface(args ...string) []interface{} {
	_args := []interface{}{}
	for _, value := range args {
		_args = append(_args, value)
	}
	return _args
}
