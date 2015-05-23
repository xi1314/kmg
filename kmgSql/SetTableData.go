package kmgSql

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/bronze1man/kmg/encoding/kmgYaml"
)

//设置表数据
// 注意会删除数据
func MustSetTableDataYaml(yaml string) {
	err := GetDb().SetTablesDataYaml(yaml)
	if err != nil {
		panic(err)
	}
}

// @deprecated
func (db DB) MustSetTablesDataYaml(yaml string) {
	err := db.SetTablesDataYaml(yaml)
	if err != nil {
		panic(err)
	}
}

// @deprecated
func (db DB) SetTablesDataYaml(yaml string) (err error) {
	data := make(map[string][]map[string]string)
	err = kmgYaml.Unmarshal([]byte(yaml), &data)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return fmt.Errorf("[SetTablesDataYaml] try to set tables with no data,wrong format?")
	}
	return db.SetTablesData(data)
}

// @deprecated
// Set some tables data in this database.
// mostly for test
// not guarantee next increment id will be!!
//设置表数据
// 注意:
//   * 会删除数据
//   * 保证 auto_increase 的值是数据里面的最大值+1
func (db DB) SetTablesData(data map[string][]map[string]string) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	err = setTablesDataTransaction(data, tx)
	if err != nil {
		errRoll := tx.Rollback()
		if errRoll != nil {
			return fmt.Errorf("error [transaction] %s,[rollback] %s", err, errRoll)
		}
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
func setTablesDataTransaction(data map[string][]map[string]string, tx *sql.Tx) error {
	for tableName, tableData := range data {
		sql := fmt.Sprintf("truncate `%s`", tableName)
		_, err := tx.Exec(sql)
		if err != nil {
			return err
		}
		for _, row := range tableData {
			colNameList := []string{}
			placeHolderNum := len(row)
			valueList := []interface{}{}
			for name, value := range row {
				colNameList = append(colNameList, name)
				valueList = append(valueList, value)
			}
			sqlColNamePart := "`" + strings.Join(colNameList, "`, `") + "`"
			sqlValuePart := strings.Repeat("?, ", placeHolderNum-1) + "?"
			sql = fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)", tableName, sqlColNamePart, sqlValuePart)
			_, err := tx.Exec(sql, valueList...)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
