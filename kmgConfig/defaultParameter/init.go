package defaultParameter

import (
	"fmt"
	"github.com/bronze1man/kmg/encoding/kmgYaml"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgConfig/defaultEnv"
	"path/filepath"
	"sync"
)

var parameterOnce sync.Once
var parameter *kmgConfig.Parameter

func Parameter() *kmgConfig.Parameter {
	parameterOnce.Do(func() {
		parameter = &kmgConfig.Parameter{}
		path := filepath.Join(defaultEnv.Env().ConfigPath, "parameters.yml")
		err := kmgYaml.ReadFile(path, parameter)
		if err != nil {
			panic(fmt.Errorf("can not get Parameters config,do you forget write a config file at %s ? err: %s", path, err))
		}
	})
	return parameter
}
