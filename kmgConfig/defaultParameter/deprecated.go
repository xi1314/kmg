package defaultParameter

import (
	"github.com/bronze1man/kmg/kmgConfig"
)

// @deprecated
func Parameter() *kmgConfig.Parameter {
	return kmgConfig.DefaultParameter()
}
