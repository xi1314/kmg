package kmgSql

import (
	. "github.com/bronze1man/kmg/kmgTest"
	"github.com/bronze1man/kmg/kmgTime"
	"testing"
	"time"
)

func TestOrmFromRow(t *testing.T) {
	obj, err := OrmFromRow(&KmgSqlTestType{}, map[string]string{
		"Id":       "1",
		"Name":     "abc",
		"TimeDate": "2001-01-01",
		"Time":     "2001-01-01 01:01:01",
	})
	Equal(err, nil)
	objTest := obj.(*KmgSqlTestType)
	Equal(err, nil)
	Equal(objTest.Id, 1)
	Equal(objTest.Name, "abc")
	Equal(objTest.TimeDate, kmgTime.MustFromMysqlFormatDefaultTZ("2001-01-01 00:00:00"))
	Equal(objTest.Time, kmgTime.MustFromMysqlFormatDefaultTZ("2001-01-01 01:01:01"))
}
func TestOrmToRow(t *testing.T) {
	m, err := OrmToRow(&KmgSqlTestType{
		Id:       2,
		Name:     "edf",
		TimeDate: kmgTime.MustFromMysqlFormatDefaultTZ("2001-01-01 00:00:00"),
		Time:     kmgTime.MustFromMysqlFormatDefaultTZ("2001-01-01 01:01:01"),
	})
	Equal(err, nil)
	Equal(m["Id"], "2")
	Equal(m["Name"], "edf")
	Equal(m["TimeDate"], "2001-01-01 00:00:00")
	Equal(m["Time"], "2001-01-01 01:01:01")
}

func TestOrmPersistUpdate(t *testing.T) {
	setTestOrmTable()
	obj, err := OrmFromRow(&KmgSqlTestType{}, map[string]string{
		"Id":   "1",
		"Name": "Jack",
	})
	Equal(err, nil)
	id, err := OrmPersist(obj)
	Equal(err, nil)
	Equal(id, 1)
	one, err := GetOneWhere(obj.GetTableName(), obj.GetIdFieldName(), "1")
	Equal(err, nil)
	obj, err = OrmFromRow(&KmgSqlTestType{}, one)
	Equal(err, nil)
	objTest := obj.(*KmgSqlTestType)
	Equal(objTest.Name, "Jack")
	Equal(objTest.Id, 1)
}

func TestOrmPersistInsert(t *testing.T) {
	setTestOrmTable()
	obj, err := OrmFromRow(&KmgSqlTestType{}, map[string]string{
		"Name": "Lucy",
	})
	Equal(err, nil)
	id, err := OrmPersist(obj)
	Equal(err, nil)
	Equal(id, 2)
	one, err := GetOneWhere(obj.GetTableName(), "Name", "Lucy")
	Equal(err, nil)
	obj, err = OrmFromRow(&KmgSqlTestType{}, one)
	Equal(err, nil)
	objTest := obj.(*KmgSqlTestType)
	Equal(objTest.Name, "Lucy")
	Equal(objTest.Id, 2)
}

type KmgSqlTestType struct {
	Id       int
	Name     string
	TimeDate time.Time
	Time     time.Time
}

func (t *KmgSqlTestType) GetIdFieldName() string {
	return "Id"
}
func (t *KmgSqlTestType) GetTableName() string {
	return "kmgSql_test_table"
}

func setTestOrmTable() {
	_, err := Exec("DROP TABLE IF EXISTS `kmgSql_test_table`")
	Equal(err, nil)
	_, err = Exec("CREATE TABLE `kmgSql_test_table` ( `Id` int(11) NOT NULL AUTO_INCREMENT, `Name` varchar(255) COLLATE utf8_bin DEFAULT NULL, `Time` datetime NOT NULL,`TimeDate` datetime NOT NULL,PRIMARY KEY (`Id`) ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin")
	Equal(err, nil)
	err = GetDb().SetTablesDataYaml(`
kmgSql_test_table:
  - Id: 1
    Name: Tom
    Time: "2015-01-12 09:23:59"
`)
	Equal(err, nil)
}
