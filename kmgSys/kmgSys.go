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

// 经过多次尝试,发现这个值只能设置到 1048576
const MaxMaxFileNum = 1048576

func MustSetCurrentMaxFileNum(limit uint64){
	err := SetCurrentMaxFileNum(limit)
	if err!=nil{
		panic(err)
	}
}