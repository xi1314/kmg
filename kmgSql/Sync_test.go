package kmgSql_test

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgSql"
	. "github.com/bronze1man/kmg/kmgTest"
	"strings"
	"testing"
)

var testTableConfig = kmgSql.Table{
	Name: "testTable",
	FieldList: map[string]kmgSql.DbType{
		"Id":   kmgSql.DbTypeInt,
		"Name": kmgSql.DbTypeString,
	},
	PrimaryKey: "Id",
	UniqueKey: [][]string{
		[]string{"Id", "Name"},
	},
	Null: []string{"Name"},
}

var testBadTableConfig = kmgSql.Table{
	Name: "testTable",
	FieldList: map[string]kmgSql.DbType{
		"Id":   kmgSql.DbTypeInt,
		"Name": kmgSql.DbTypeString,
	},
}

func TestIsTableExist(t *testing.T) {
	setTest()
	Ok(!kmgSql.MustIsTableExist("testTable"))
}

func TestCreateTable(t *testing.T) {
	setTest()
	kmgSql.MustCreateTable(testTableConfig)
	ret := kmgSql.MustQueryOne("SHOW CREATE TABLE testTable")
	Ok(strings.Contains(fmt.Sprint(ret), "Id"))
	Ok(strings.Contains(fmt.Sprint(ret), "Name"))
	Ok(kmgSql.MustIsTableExist("testTable"))
}

func TestModifyTable(t *testing.T) {
	setTest()
	kmgSql.MustCreateTable(testTableConfig)
	newTestTableConfig := kmgSql.Table{
		Name: "testTable",
		FieldList: map[string]kmgSql.DbType{
			"Id":   kmgSql.DbTypeInt,
			"Name": kmgSql.DbTypeString,
			"Age":  kmgSql.DbTypeInt,
		},
		PrimaryKey: "Id",
		UniqueKey: [][]string{
			[]string{"Id", "Name"},
		},
		Null: []string{"Name"},
	}
	kmgSql.MustModifyTable(newTestTableConfig)
	ret := kmgSql.MustQueryOne("SHOW CREATE TABLE testTable")
	Ok(strings.Contains(fmt.Sprint(ret), "Id"))
	Ok(strings.Contains(fmt.Sprint(ret), "Name"))
	Ok(strings.Contains(fmt.Sprint(ret), "Age"))
	Ok(kmgSql.MustIsTableExist("testTable"))
}

func TestSyncTable(t *testing.T) {
	setTest()
	kmgSql.MustSyncTable(testTableConfig)
	Ok(kmgSql.MustIsTableExist("testTable"))
}

func TestSyncTableBad(t *testing.T) {
	setTest()
	kmgSql.MustSyncTable(testBadTableConfig)
	Ok(kmgSql.MustIsTableExist("testTable"))
}

func setTest() {
	kmgSql.SetDefaultDbConfig(kmgSql.MustGetTestConfig().Db)
	kmgSql.Exec("DROP TABLE IF EXISTS `testTable`")
}
