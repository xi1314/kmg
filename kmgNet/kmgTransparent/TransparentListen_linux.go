package kmgTransparent

import (
	"net"
	"os"

	"github.com/bronze1man/kmg/kmgNet/kmgUnix"
	"golang.org/x/sys/unix"
)

//linux transparent listen
//use for iptables TProxy
func TransparentListen(listenAddr string) (listener net.Listener, err error) {
	s, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		return nil, err
	}
	defer unix.Close(s)
	err = unix.SetsockoptInt(s, unix.SOL_IP, unix.IP_TRANSPARENT, 1)
	if err != nil {
		return nil, err
	}
	sa, err := kmgUnix.IPv4TcpAddrToUnixSocksAddr(listenAddr)
	if err != nil {
		return nil, err
	}
	err = unix.Bind(s, sa)
	if err != nil {
		return nil, err
	}
	err = unix.Listen(s, 10)
	if err != nil {
		return nil, err
	}
	f := os.NewFile(uintptr(s), "TProxy")
	defer f.Close()
	return net.FileListener(f)
}
