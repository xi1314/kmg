package defaultEnv

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgConfig"
)

var Env *kmgConfig.Env

func init() {
	var err error
	Env, err = kmgConfig.LoadEnvFromWd()
	if err != nil {
		panic(fmt.Errorf("can not getEnv,do you forget create a .kmg.yml at project root? err: %s", err))
	}
}
