package kmgProxy

import (
	"net"
	"strconv"
)

//没有解析过的地址,可能有ip,有域名,也可能有端口号
type UnsolvedAddr struct {
	Ip     net.IP
	Domain string
	Port   uint16
}

func (uaddr UnsolvedAddr) String() string {
	if uaddr.Ip != nil {
		return net.JoinHostPort(uaddr.Ip.String(), strconv.Itoa(int(uaddr.Port)))
	} else {
		return net.JoinHostPort(uaddr.Domain, strconv.Itoa(int(uaddr.Port)))
	}
}

func ParseUnsolvedAddr(addr string) (uaddr UnsolvedAddr, err error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return
	}
	iport, err := strconv.Atoi(port)
	if err != nil {
		return
	}
	uaddr.Port = uint16(iport)
	uaddr.Ip = net.ParseIP(host)
	if uaddr.Ip == nil {
		uaddr.Domain = host
		return uaddr, nil
	}
	return uaddr, nil
}
