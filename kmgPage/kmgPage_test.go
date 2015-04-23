package kmgPage

import (
	"github.com/bronze1man/kmg/kmgSql"
	"github.com/bronze1man/kmg/kmgSql/MysqlAst"
	"github.com/bronze1man/kmg/kmgTest"
	"strings"
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
	var page KmgPage
	data_init := test_init()
	data := page.CreateFromData(data_init, "http://sig.com", 1)
	kmgTest.Equal(data.BaseUrl, "http://sig.com")
	kmgTest.Equal(data.ItemPerPage, 1)
	kmgTest.Equal(data.CurrentPage, 1)
	kmgTest.Equal(data.TotalItem, 3)
	kmgTest.Equal(page.GetTotalPage(), 3)
	kmgTest.Ok(page.IsAfterPageActive())
	kmgTest.Ok(!page.IsBeforePageActive())
}

func TestGetUrlWithPage(ot *testing.T) {
	var page KmgPage
	data_init := test_init()
	data := page.CreateFromData(data_init, "http://sig.com?home=eqw", 1)
	kmgTest.Ok(len(Splite(data.GetUrlWithPage(2), "?")) == 2)
	//	fmt.Println(len(Splite(data.GetUrlWithPage(2),"?")))
}

func TestKmgPageWithDb(ot *testing.T) {
	kmgSql.MustExec("DROP TABLE IF EXISTS `kmgSql_test_table`")
	kmgSql.MustExec("CREATE TABLE `kmgSql_test_table` ( `Id` int(11) NOT NULL AUTO_INCREMENT, `Info` varchar(255) COLLATE utf8_bin DEFAULT NULL, PRIMARY KEY (`Id`) ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin")
	kmgSql.MustSetTableDataYaml(`
kmgSql_test_table:
   - Id: 1
   - Id: 2
   - Id: 3
   - Id: 4
`)
	pager := CreateFromSelectCommand(CreateFromSelectCommandRequest{
		Select:      MysqlAst.NewSelectCommand().From("kmgSql_test_table"),
		Url:         "/?n=a",
		ItemPerPage: 2,
		CurrentPage: 1,
	})
	kmgTest.Equal(len(pager.Data), 2)

	pager = CreateFromSelectCommand(CreateFromSelectCommandRequest{
		Select: MysqlAst.NewSelectCommand().From("kmgSql_test_table"),
		Url:    "/?n=a",
	})
	kmgTest.Equal(len(pager.Data), 4)
}

func Splite(s, sep string) []string {
	return strings.Split(s, sep)
}
