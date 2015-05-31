package kmgSql

import (
	"github.com/bronze1man/kmg/typeTransform"
)

type OrmObject interface {
	GetIdFieldName() string
	GetTableName() string
}

func OrmFromRow(obj OrmObject, row map[string]string) (OrmObject, error) {
	err := typeTransform.Transform(row, &obj)
	return obj, err
}

func MustOrmFromRow(obj OrmObject, row map[string]string) {
	err := typeTransform.Transform(row, &obj)
	if err != nil {
		panic(err)
	}
}

func OrmToRow(obj OrmObject) (row map[string]string, err error) {
	row = map[string]string{}
	err = typeTransform.Transform(obj, &row)
	return row, err
}

func OrmPersist(obj OrmObject) (lastInsertId int, err error) {
	row, err := OrmToRow(obj)
	if err != nil {
		return 0, err
	}
	return ReplaceById(obj.GetTableName(), obj.GetIdFieldName(), row)
}

func GetOrmById(obj OrmObject, id string) {
	row := MustGetOneWhere(obj.GetTableName(), obj.GetIdFieldName(), id)
	if row == nil {
		obj = nil
		return
	}
	obj, err := OrmFromRow(obj, row)
	if err != nil {
		panic(err)
	}
}
