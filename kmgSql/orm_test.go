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
	return "tbf_test"
}
