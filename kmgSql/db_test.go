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
	err := DeleteById("AdminUser","Id","3")
	kmgDebug.Println(err)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
