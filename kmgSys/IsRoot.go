package kmgSys

import (
	"bytes"
	"github.com/bronze1man/kmg/kmgCmd"
)

//是否是root,此处只返回是否,其他错误抛panic
func MustIsRoot() bool {
	return bytes.Equal(kmgCmd.MustCombinedOutput("whoami"), []byte("root\n"))
}
