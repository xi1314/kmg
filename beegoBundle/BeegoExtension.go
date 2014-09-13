package beegoBundle

import (
	"github.com/astaxie/beego/orm"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgContext"
	"time"
)

type tBeegoOrmKey struct{}

var beegoOrmKey tBeegoOrmKey = tBeegoOrmKey{}

func init() {
	orm.RegisterDataBase("default", "mysql", kmgConfig.DefParameter.GetDbConfig().GetDsn())
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
