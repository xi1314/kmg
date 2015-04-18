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
	if err != nil {
		return nil, err
	}
	return obj, err
}

func OrmToRow(obj OrmObject) (row map[string]string, err error) {
	row = map[string]string{}
	err = typeTransform.Transform(obj, &row)
	if err != nil {
		return nil, err
	}
	return row, err
}

func OrmPersist(obj OrmObject) (lastInsertId int, err error) {
	row, err := OrmToRow(obj)
	if err != nil {
		return 0, err
	}
	return ReplaceById(obj.GetTableName(), row, obj.GetIdFieldName())
}
