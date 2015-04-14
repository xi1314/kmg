package kmgTfo

import (
	"golang.org/x/sys/unix"

	"github.com/bronze1man/kmg/kmgNet/kmgUnix"
	"net"
	"os"
)

//network is useless ,it will always use tcp4
func TfoListen(network string, listenAddr string) (listener net.Listener, err error) {
	s, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM|unix.SOCK_NONBLOCK|unix.SOCK_CLOEXEC, 0)
	if err != nil {
		return nil, err
	}
	defer unix.Close(s)
	err = unix.SetsockoptInt(s, unix.SOL_TCP, 23, 10)
	if err != nil {
		return nil, err
	}
	err = unix.SetsockoptInt(s, unix.SOL_SOCKET, unix.SO_REUSEADDR, 1)
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
	f := os.NewFile(uintptr(s), "TFOListen")
	defer f.Close()
	return net.FileListener(f)
}
