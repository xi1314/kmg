package kmgNet

import (
	"errors"
	"net"
	"strconv"
)

var ProtocolNotSupportPortError = errors.New("Protocol Not Support Port")

func PortFromNetAddr(addr net.Addr) (int, error) {
	switch saddr := addr.(type) {
	case *net.TCPAddr:
		return saddr.Port, nil
	case *net.UDPAddr:
		return saddr.Port, nil
	case *net.IPAddr:
		return 0, ProtocolNotSupportPortError
	case *net.UnixAddr:
		return -1, ProtocolNotSupportPortError
	}
	return PortFromAddrString(addr.String())
}

func PortFromAddrString(addr string) (int, error) {
	_, portS, err := net.SplitHostPort(addr)
	if err != nil {
		return 0, err
	}
	portI, err := strconv.Atoi(portS)
	if err != nil {
		return 0, err
	}
	return portI, nil
}

func JoinHostPortInt(host string, port int) string {
	return net.JoinHostPort(host, strconv.Itoa(port))
}

func MustGetHostFromAddr(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		panic(err)
	}
	return host
}
