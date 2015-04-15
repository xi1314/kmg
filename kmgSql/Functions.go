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

//func UpdateById(tableName string, row map[string]string, primaryKeyName string) {
//}

//func Replace(query string, args ...string) (sql.Result, error) {
//}

func argsStringToInterface(args ...string) []interface{} {
	_args := []interface{}{}
	for _, value := range args {
		_args = append(_args, value)
	}
	return _args
}
