package kmgNet

import (
	"fmt"
	"golang.org/x/sys/unix"
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
	tcpAddr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	ip := tcpAddr.IP.To4()
	if ip == nil {
		return nil, fmt.Errorf("[TProxyListen] only support tcp4 right now.")
	}
	var ipA [4]byte
	copy(ipA[:], ip[:4])
	sockAddr := &unix.SockaddrInet4{
		Port: tcpAddr.Port,
		Addr: ipA,
	}
	err = unix.Bind(s, sockAddr)
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
