package defaultEnv

import (
	"github.com/bronze1man/kmg/kmgConfig"
)

// @deprecated
func Env() *kmgConfig.Env {
	return kmgConfig.DefaultEnv()
}
