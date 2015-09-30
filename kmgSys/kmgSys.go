package kmgSys

import (
	"errors"
	"github.com/bronze1man/kmg/kmgErr"
	"os/user"
)

var ErrPlatformNotSupport = errors.New("Platform Not Support")

func GetCurrentUserHomeDir() string {
	u, err := user.Current()
	kmgErr.PanicIfError(err)
	return u.HomeDir
}
