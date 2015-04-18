package kmgSql

import (
	"database/sql"
	"fmt"
	"strings"
)

func Query(query string, args ...string) (output []map[string]string, error error) {
	rows, err := GetDb().Query(query, argsStringToInterface(args...)...)
	if err != nil {
		return nil, err
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

func QueryOne(query string, args ...string) (output map[string]string, error error) {
	list, error := Query(query, args...)
	if error != nil {
		return nil, error
	}
	if len(list) <= 0 {
		return nil, error
	}
	output = list[0]
	return output, error
}

func Exec(query string, args ...string) (sql.Result, error) {
	return GetDb().Exec(query, argsStringToInterface(args...)...)
}

func Insert(tableName string, row map[string]string) (lastInsertId int, error error) {
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

func UpdateById(tableName string, row map[string]string, primaryKeyName string) error {
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

func GetOneWhere(tableName string, fieldName string, value string) (output map[string]string, error error) {
	sql := fmt.Sprintf("SELECT * FROM `%s` WHERE `%s`=? LIMIT 1", tableName, fieldName)
	output, error = QueryOne(sql, value)
	if error != nil {
		return nil, error
	}
	return output, error
}

func DeleteById(tableName string, fieldName string, value string) error {
	sql := fmt.Sprintf("DELETE FROM `%s` WHERE `%s`=?", tableName, fieldName)
	_, err := Exec(sql, value)
	return err
}

//func Replace(query string, args ...string) (sql.Result, error) {
//}

func argsStringToInterface(args ...string) []interface{} {
	_args := []interface{}{}
	for _, value := range args {
		_args = append(_args, value)
	}
	return _args
}
