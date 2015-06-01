package kmgPage_test

import (
	"github.com/bronze1man/kmg/kmgPage"
	"github.com/bronze1man/kmg/kmgSql"
	"github.com/bronze1man/kmg/kmgSql/MysqlAst"
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
)

func test_init() []map[string]string {
	test_map := map[string]string{"a": "a1", "b": "b1"}
	arr := []map[string]string{test_map, test_map, test_map}
	arr[0] = test_map
	arr[1] = test_map
	arr[2] = test_map
	return arr
}

func TestPagination(ot *testing.T) {
	page := &kmgPage.KmgPage{}
	page.TotalItem = 10
	page.ItemPerPage = 2
	kmgTest.Equal(page.GetTotalPage(), 5)
}

func TestKmgPageWithDb(ot *testing.T) {
	setupDb()
	kmgSql.MustExec("DROP TABLE IF EXISTS `kmgSql_test_table`")
	kmgSql.MustExec("CREATE TABLE `kmgSql_test_table` ( `Id` int(11) NOT NULL AUTO_INCREMENT, `Info` varchar(255) COLLATE utf8_bin DEFAULT NULL, PRIMARY KEY (`Id`) ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin")
	kmgSql.MustSetTableDataYaml(`
kmgSql_test_table:
   - Id: 1
   - Id: 2
   - Id: 3
   - Id: 4
`)
	pager := kmgPage.CreateFromSelectCommand(kmgPage.CreateFromSelectCommandRequest{
		Select:      MysqlAst.NewSelectCommand().From("kmgSql_test_table"),
		Url:         "/?n=a",
		ItemPerPage: 2,
		CurrentPage: 1,
	})
	kmgTest.Equal(len(pager.Data), 2)

	pager = kmgPage.CreateFromSelectCommand(kmgPage.CreateFromSelectCommandRequest{
		Select: MysqlAst.NewSelectCommand().From("kmgSql_test_table"),
		Url:    "/?n=a",
	})
	kmgTest.Equal(len(pager.Data), 4)
}

func setupDb() {
	kmgSql.MustLoadTestConfig()
}
