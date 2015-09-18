package kmgSys

import (
	"bytes"
	"fmt"
	"github.com/bronze1man/kmg/kmgCmd"
	"os"
)

//是否是root,此处只返回是否,其他错误抛panic
// TODO 名字比较费解.
func MustIsRoot() bool {
	return bytes.Equal(kmgCmd.MustCombinedOutput("whoami"), []byte("root\n"))
}

func MustIsRootOnCmd() {
	if !MustIsRoot() {
		fmt.Println("need root to run this command.")
		os.Exit(1)
	}
}
