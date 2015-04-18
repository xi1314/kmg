package kmgPage

import (
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
)

type Data struct {
	ItemPerPage int                 // 每页数据量
	ShowPageNum int                 // 显示分页的页面的时候，显示出来的页数 如 10，11，12，13，14，15 就是显示了6个页
	PageKeyName string              // 页面参数的名字
	CurrentPage int                 // 当前页数，默认为1
	TotalItem   int                 // 总数据量
	BaseUrl     string              // 不包含页码参数的url,具体渲染url的时候会把页面参数加上
	Data        []map[string]string // 本次分页查询到的数据
	TotalPage   int                 // 共有多少个页面
}

func test_init() []map[string]string {

	test_map := map[string]string{"a": "a1", "b": "b1"}
	arr := []map[string]string{test_map, test_map, test_map}
	arr[0] = test_map
	arr[1] = test_map
	arr[2] = test_map
	return arr
}

func TestPagination(ot *testing.T) {
	var page kmgPage
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
