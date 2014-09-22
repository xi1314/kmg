package beegoBundle

import (
	"github.com/astaxie/beego/orm"
	"github.com/bronze1man/kmg/kmgContext"
	"github.com/bronze1man/kmg/kmgSql"

	"github.com/bronze1man/kmg/kmgConfig/defaultParameter"
	"time"
)

type tBeegoOrmKey struct{}

var beegoOrmKey tBeegoOrmKey = tBeegoOrmKey{}

func InitOrm() {
	orm.RegisterDataBase("default", "mysql", kmgSql.GetDbConfigFromConfig(defaultParameter.Parameter()).GetDsn())
	orm.SetDataBaseTZ("default", time.UTC)
}

func ContextGetOrm(c kmgContext.Context) orm.Ormer {
	o, ok := c.Value(beegoOrmKey).(orm.Ormer)
	if !ok {
		o = orm.NewOrm()
		c.SetValue(beegoOrmKey, o)
	}
	return o
}
