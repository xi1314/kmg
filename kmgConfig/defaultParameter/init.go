package defaultParameter

import (
	"fmt"
	"github.com/bronze1man/kmg/encoding/kmgYaml"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgConfig/defaultEnv"
	"path/filepath"
)

var Parameter *kmgConfig.Parameter

func init() {
	Parameter = &kmgConfig.Parameter{}
	path := filepath.Join(defaultEnv.Env.ConfigPath, "parameters.yml")
	err := kmgYaml.ReadFile(path, Parameter)
	if err != nil {
		panic(fmt.Errorf("can not get Parameters config,do you forget write a config file at %s ? err: %s", path, err))
	}
}
