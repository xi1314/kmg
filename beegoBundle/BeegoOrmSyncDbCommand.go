package beegoBundle

import (
	"flag"
	"github.com/astaxie/beego/orm"
	"github.com/bronze1man/kmg/console"
	"os"
)

type BeegoOrmSyncDbCommand struct {
	env string
}

func (command *BeegoOrmSyncDbCommand) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{
		Name:  "BeegoOrmSyncDb",
		Short: "beego orm command",
	}
}

func (command *BeegoOrmSyncDbCommand) ConfigFlagSet(flag *flag.FlagSet) {
	flag.StringVar(&command.env, "env", "dev", "database env(dev,test)")
}
func (command *BeegoOrmSyncDbCommand) Execute(context *console.Context) error {
	InitOrm()
	//TODO register database config stuff.
	os.Args = []string{
		os.Args[0], "orm", "syncdb",
	}
	orm.RunCommand()
	return nil
}
