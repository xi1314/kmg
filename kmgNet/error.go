package kmgNet

import (
	"errors"
	"io"
	"strings"
)

var ErrClosing = errors.New("use of closed network connection")

//是否是因为socket没有关闭,或者socket根本没有连接而导致的错误,或者被reset (表示这个连接应该被马上关闭)
func IsSocketCloseError(err error) bool {
	return err != nil && (err == io.EOF || //conn.Read
		err == io.ErrClosedPipe || //conn.Read
		strings.Contains(err.Error(), "use of closed network connection") || //来自 conn.Read
		strings.Contains(err.Error(), "socket is not connected")) //conn.CloseRead shutdown tcp 127.0.0.1:30002: socket is not connected
	//strings.Contains(err.Error(), "connection reset by peer") || //来自 conn.Read
	//err == io.ErrClosedPipe) //来自 conn.Read
	//strings.Contains(err.Error(), "Stream closed") || //来自 muxado
	//strings.Contains(err.Error(), "Session closed")) //来自 muxado

}

func IsBadFileDescripter(err error) bool {
	return err != nil && strings.Contains(err.Error(), "bad file descriptor")
}

// icmp 报告服务器不存在,通常是服务器掉进程了.
// 出现位置 udp write
func IsConnectionRefused(err error) bool {
	return err != nil && strings.Contains(err.Error(), "connection refused")
}

// 当前设备没有可用的网络设备,在网络切换时会出现
// 出现位置 udp read, udp write,udp dial
func IsNetworkIsUnreachable(err error) bool {
	return err != nil && strings.Contains(err.Error(), "network is unreachable")
}
