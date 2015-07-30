package kmgSql

import (
	"testing"

	"github.com/bronze1man/kmg/kmgTest"
)

func TestMustTransactionCallback(ot *testing.T) {
	setTestSqlTable()
	MustTransactionCallback(func(tx Tx) {
		tx.Insert("kmgSql_test_table", map[string]string{
			"Info": "Hello123",
		})
	})
	all := MustGetAllInTable("kmgSql_test_table")
	kmgTest.Equal(len(all), 3)

	out := kmgTest.AssertPanic(func() {
		MustTransactionCallback(func(tx Tx) {
			tx.Insert("kmgSql_test_table", map[string]string{
				"Info": "Hello123",
			})
			panic("abc")
		})
	})
	kmgTest.Equal(out, "abc")
	all = MustGetAllInTable("kmgSql_test_table")
	kmgTest.Equal(len(all), 3)

	kmgTest.AssertPanic(func() {
		MustTransactionCallback(func(tx Tx) {
			_, err := tx.Insert("kmgSql_test_table", map[string]string{
				"Info": "Hello123",
			})
			if err != nil {
				panic(err)
			}
			_, err = tx.Insert("kmgSql_test_table", map[string]string{
				"Info":                 "Hello123",
				"InfoWhateverNotExist": "abc",
			})
			if err != nil {
				panic(err)
			}
		})
	})
	all = MustGetAllInTable("kmgSql_test_table")
	kmgTest.Equal(len(all), 3)
}
