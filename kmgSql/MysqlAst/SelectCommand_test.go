package MysqlAst

import (
	"testing"

	. "github.com/bronze1man/kmg/kmgTest"
)

func TestSelectCommand1(t *testing.T) {
	text, parameterList := NewSelectCommand().From("tbf_test").Where("real_name=1").GetPrepareParameter()
	Equal(text, "SELECT * FROM tbf_test WHERE real_name=1")
	Equal(len(parameterList), 0)
}

func TestSelectCommand2(t *testing.T) {
	s1 := NewSelectCommand().From("Table1").Where("a=1 AND b=?", "c").Limit("1,2")
	text, parameterList := s1.GetPrepareParameter()
	Equal(text, "SELECT * FROM Table1 WHERE a=1 AND b=? LIMIT 1,2")
	Equal(parameterList, []string{"c"})

	s2 := s1.Copy()
	s2.Limit("2,2")
	text, _ = s1.GetPrepareParameter()
	Equal(text, "SELECT * FROM Table1 WHERE a=1 AND b=? LIMIT 1,2")
	text, _ = s2.GetPrepareParameter()
	Equal(text, "SELECT * FROM Table1 WHERE a=1 AND b=? LIMIT 2,2")
}

func TestAndWhereConditionAddPrepare(t *testing.T) {
	and := NewAndWhereCondition().AddPrepare("a=1").AddPrepare("d=?", "c")
	s1 := NewSelectCommand().From("Table1").WhereObj(and)
	text, parameterList := s1.GetPrepareParameter()
	Equal(text, "SELECT * FROM Table1 WHERE (a=1) AND (d=?)")
	Equal(parameterList, []string{"c"})

	s2 := s1.Copy()
	text, parameterList = s2.GetPrepareParameter()
	Equal(text, "SELECT * FROM Table1 WHERE (a=1) AND (d=?)")
	Equal(parameterList, []string{"c"})
}

func TestOrWhereConditionAddPrepare(t *testing.T) {
	and := NewOrWhereCondition().AddPrepare("a=1").AddPrepare("d=?", "c")
	s1 := NewSelectCommand().From("Table1").WhereObj(and)
	text, parameterList := s1.GetPrepareParameter()
	Equal(text, "SELECT * FROM Table1 WHERE (a=1) OR (d=?)")
	Equal(parameterList, []string{"c"})
}
