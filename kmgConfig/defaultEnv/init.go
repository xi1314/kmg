package defaultEnv

import (
	"fmt"
	"sync"

	"github.com/bronze1man/kmg/kmgConfig"
)

var envOnce sync.Once
var env *kmgConfig.Env

func Env() *kmgConfig.Env {
	envOnce.Do(func() {
		var err error
		env, err = kmgConfig.LoadEnvFromWd()
		if err != nil {
			panic(fmt.Errorf("can not getEnv,do you forget create a .kmg.yml at project root? err: %s", err))
		}
	})
	return env
}
