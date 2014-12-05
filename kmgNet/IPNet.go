package kmgNet

import (
	"net"
)

func MustParseCIDR(s string) *net.IPNet {
	_, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		panic(err)
	}
	return ipnet
}

var localAddrMaskList = []*net.IPNet{
	MustParseCIDR("0.0.0.0/8"),      // broadcast messages
	MustParseCIDR("10.0.0.0/8"),     // private network
	MustParseCIDR("127.0.0.0/8"),    // loopback addresses
	MustParseCIDR("172.16.0.0/12"),  // private network
	MustParseCIDR("169.254.0.0/16"), // link-local addresses
	MustParseCIDR("192.168.0.0/16"), // private network
	MustParseCIDR("198.18.0.0/15"),  // private network
	MustParseCIDR("::1/128"),        // loopback addresses
	MustParseCIDR("fc00::/7"),
	MustParseCIDR("fe80::/10"),
}

func IsPrivateNetwork(ip net.IP) bool {
	for _, ipnet := range localAddrMaskList {
		if ipnet.Contains(ip) {
			return true
		}
	}
	return false
}
