package kmgCmd

import "bytes"

//是否是root,此处只返回是否,其他错误抛panic
func MustIsRoot() bool {
	return bytes.Equal(MustRunAndReturnOutput("whoami"), []byte("root\n"))
}
