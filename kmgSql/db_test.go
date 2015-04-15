package kmgSql

import (
	"github.com/bronze1man/kmg/kmgDebug"
	. "github.com/bronze1man/kmg/kmgTest"
	"testing"
)

func TestConnectToDb(t *testing.T) {
	db := GetDb()
	err := db.Ping()
	Equal(err, nil)
}

func TestQuery(t *testing.T) {
	//	rows, err := Query("SELECT * FROM AdminUser WHERE Id=?", 1)
	//	handleError(err)
	//	for key, value := range rows {
	//		kmgDebug.Println(key, value)
	//	}
	//	handleError(err)
	//	one, err := QueryOne("SELECT * FROM AdminUser WHERE Id=?", "1")
	//	handleError(err)
	//	kmgDebug.Println(one)
	//	_,err = Exec("DELETE FROM ArticleList where Id=?","1")
	//	handleError(err)
	row := map[string]string{
		"Title":  "Degas",
		"Source": "123",
	}
	id, err := Insert("ArticleList", row)
	handleError(err)
	kmgDebug.Println(id)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
