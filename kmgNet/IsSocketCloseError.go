package kmgNet

import (
	"io"
	"strings"
)

//是否是因为socket没有关闭,或者socket根本没有连接而导致的错误,或者被reset (表示这个连接应该被马上关闭)
func IsSocketCloseError(err error) bool {
	return err != nil && (strings.Contains(err.Error(), "use of closed network connection") || //来自 conn.Read
		//strings.Contains(err.Error(), "connection reset by peer") || //来自 conn.Read
		err == io.ErrClosedPipe) //来自 conn.Read
	//strings.Contains(err.Error(), "Stream closed") || //来自 muxado
	//strings.Contains(err.Error(), "Session closed")) //来自 muxado
}
