package kmgSql

import (
	. "github.com/bronze1man/kmg/kmgTest"
	"testing"
)

func TestConnectToDb(t *testing.T) {
	db := GetDb()
	err := db.Ping()
	Equal(err, nil)
}

func TestExec(t *testing.T) {
	setTestSqlTable()
	_, err := Exec("DELETE FROM `kmgSql_test_table` WHERE Id=? AND Info=?", "2", "World")
	Equal(err, nil)
}

func TestQuery(t *testing.T) {
	setTestSqlTable()
	row, err := Query("select * from kmgSql_test_table")
	Equal(err, nil)
	Equal(len(row), 2)
	setTestSqlTable()
	rowA, err := Query("select * from kmgSql_test_table limit 1")
	rowB, err := QueryOne("select * from kmgSql_test_table")
	Equal(rowA[0]["Id"], rowB["Id"])
	Equal(rowA[0]["Info"], rowB["Info"])

	_, err = Query("SELECT * from kmgSql_test_table WHERE (Id=?) LIMIT 0,10", "2")
	Equal(err, nil)
}

func TestInsert(t *testing.T) {
	setTestSqlTable()
	id, err := Insert("kmgSql_test_table", map[string]string{
		"Id":   "3",
		"Info": "Tom",
	})
	Equal(err, nil)
	Equal(id, 3)
	one, err := QueryOne("select * from kmgSql_test_table where Id=?", "3")
	Equal(one["Info"], "Tom")
	Equal(err, nil)
}

func TestUpdateById(t *testing.T) {
	setTestSqlTable()
	err := UpdateById("kmgSql_test_table", "Id", map[string]string{
		"Id":   "1",
		"Info": "Ok",
	})
	Equal(err, nil)
	one, err := QueryOne("select * from kmgSql_test_table where Id=?", "1")
	Equal(one["Info"], "Ok")
	Equal(err, nil)
}

func TestReplaceById(t *testing.T) {
	setTestSqlTable()
	row := map[string]string{
		"Id":   "1",
		"Info": "Sky",
	}
	id, err := ReplaceById("kmgSql_test_table", "Id", row)
	Equal(err, nil)
	Equal(id, 1)
	setTestSqlTable()
	row = map[string]string{
		"Info": "Sky",
	}
	id, err = ReplaceById("kmgSql_test_table", "Id", row)
	Equal(err, nil)
	Equal(id, 3)
	one, err := GetOneWhere("kmgSql_test_table", "Id", "3")
	Equal(one["Info"], "Sky")
	Equal(err, nil)
	setTestSqlTable()
	row = map[string]string{
		"Id":   "10",
		"Info": "Sky",
	}
	id, err = ReplaceById("kmgSql_test_table", "Id", row)
	Equal(err, nil)
	Equal(id, 10)
	one, err = GetOneWhere("kmgSql_test_table", "Id", "10")
	Equal(one["Info"], "Sky")
	Equal(err, nil)
}

func TestGetOneWhere(t *testing.T) {
	setTestSqlTable()
	one, err := GetOneWhere("kmgSql_test_table", "Id", "1")
	Equal(err, nil)
	Equal(one["Info"], "Hello")
}

func TestDeleteById(t *testing.T) {
	setTestSqlTable()
	err := DeleteById("kmgSql_test_table", "Id", "1")
	Equal(err, nil)
	one, err := GetOneWhere("kmgSql_test_table", "Id", "1")
	Equal(one, nil)
	Equal(err, nil)
}

func TestGetAllInTable(t *testing.T) {
	setTestSqlTable()
	row, err := GetAllInTable("kmgSql_test_table")
	Equal(err, nil)
	Equal(len(row), 2)
}

func setTestSqlTable() {
	_, err := Exec("DROP TABLE IF EXISTS `kmgSql_test_table`")
	Equal(err, nil)
	_, err = Exec("CREATE TABLE `kmgSql_test_table` ( `Id` int(11) NOT NULL AUTO_INCREMENT, `Info` varchar(255) COLLATE utf8_bin DEFAULT NULL, PRIMARY KEY (`Id`) ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin")
	Equal(err, nil)
	err = GetDb().SetTablesDataYaml(`
kmgSql_test_table:
  - Id: 1
    Info: Hello
  - Id: 2
    Info: World
`)
	Equal(err, nil)
}
