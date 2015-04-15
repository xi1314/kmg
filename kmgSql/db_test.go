package kmgSql

import (
	"testing"
	//	."github.com/bronze1man/kmg/kmgTest"
	"fmt"
	"github.com/bronze1man/kmg/kmgConfig/defaultParameter"
	"github.com/bronze1man/kmg/kmgDebug"
)

func TestConnectToDb(t *testing.T) {
	db := GetDb()
	config := GetDbConfigFromConfig(defaultParameter.Parameter()).GetDsn()
	kmgDebug.Println(config)
	kmgDebug.Println(db)
	_, err := db.GetTableData("AdminUser")
	//fmt.Printf("%#v",defaultParameter.Parameter())
	//	kmgDebug.Println(defaultParameter.Parameter())
	//	for key, value := range row {
	//		kmgDebug.Println(key, value)
	//	}
	if err != nil {
		fmt.Errorf("%v", err)
		kmgDebug.Println(err)
	}
	//	fmt.Printf("%v", len(row))
}
