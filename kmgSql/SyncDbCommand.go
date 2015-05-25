package kmgSql

import (
	"flag"
)

func SyncDbCommand() {
	var Force bool
	flag.BoolVar(&Force, "Force", false, "force change the database")
	flag.Parse()
	//切换到prod配置
	MustLoadProdConfig()
	if Force {
		MustForceSyncDefaultDbConfig()
	} else {
		MustSyncDefaultDbConfig()
	}

	//切换到test配置
	MustLoadTestConfig()
	if Force {
		MustForceSyncDefaultDbConfig()
	} else {
		MustSyncDefaultDbConfig()
	}
}
