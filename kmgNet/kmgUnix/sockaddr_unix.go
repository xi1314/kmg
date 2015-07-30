// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package kmgUnix

import (
	"net"

	"golang.org/x/sys/unix"
)

func IPv4TcpAddrToUnixSocksAddr(addr string) (sa unix.Sockaddr, err error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	return ipToSocksAddr(unix.AF_INET, tcpAddr.IP, tcpAddr.Port, tcpAddr.Zone)
}

//copy from http://golang.org/src/net/ipsock_posix.go?s=4615:4700#L130
func ipToSocksAddr(family int, ip net.IP, port int, zone string) (unix.Sockaddr, error) {
	switch family {
	case unix.AF_INET:
		if len(ip) == 0 {
			ip = net.IPv4zero
		}
		if ip = ip.To4(); ip == nil {
			return nil, net.InvalidAddrError("non-IPv4 address")
		}
		sa := new(unix.SockaddrInet4)
		for i := 0; i < net.IPv4len; i++ {
			sa.Addr[i] = ip[i]
		}
		sa.Port = port
		return sa, nil
	case unix.AF_INET6:
		if len(ip) == 0 {
			ip = net.IPv6zero
		}
		// IPv4 callers use 0.0.0.0 to mean "announce on any available address".
		// In IPv6 mode, Linux treats that as meaning "announce on 0.0.0.0",
		// which it refuses to do.  Rewrite to the IPv6 unspecified address.
		if ip.Equal(net.IPv4zero) {
			ip = net.IPv6zero
		}
		if ip = ip.To16(); ip == nil {
			return nil, net.InvalidAddrError("non-IPv6 address")
		}
		sa := new(unix.SockaddrInet6)
		for i := 0; i < net.IPv6len; i++ {
			sa.Addr[i] = ip[i]
		}
		sa.Port = port
		sa.ZoneId = uint32(zoneToInt(zone))
		return sa, nil
	}
	return nil, net.InvalidAddrError("unexpected socket family")
}

func zoneToInt(zone string) int {
	if zone == "" {
		return 0
	}
	if ifi, err := net.InterfaceByName(zone); err == nil {
		return ifi.Index
	}
	n, _, _ := dtoi(zone, 0)
	return n
}

const big = 0xFFFFFF

func dtoi(s string, i0 int) (n int, i int, ok bool) {
	n = 0
	for i = i0; i < len(s) && '0' <= s[i] && s[i] <= '9'; i++ {
		n = n*10 + int(s[i]-'0')
		if n >= big {
			return 0, i, false
		}
	}
	if i == i0 {
		return 0, i, false
	}
	return n, i, true
}
