package kmgNet

import (
	"github.com/bronze1man/kmg/kmgTest"
	"net"
	"testing"
)

func TestGetCurrentDeviceAddrFromIPAddr(ot *testing.T) {
	addrs, err := getCurrentDeviceAddrFromIPAddr([]byte(`1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP group default qlen 1000
    link/ether 00:16:3e:00:05:7a brd ff:ff:ff:ff:ff:ff
    inet 192.168.1.2/24 brd 192.168.1.255 scope global eth0
       valid_lft forever preferred_lft forever
    inet 172.20.8.1/32 scope global eth0
       valid_lft forever preferred_lft forever
307: ppp0: <POINTOPOINT,MULTICAST,NOARP,UP,LOWER_UP> mtu 1496 qdisc pfifo_fast state UNKNOWN group default qlen 3
    link/ppp
    inet 172.20.0.1 peer 172.20.0.2/32 scope global ppp0
       valid_lft forever preferred_lft forever
`))
	kmgTest.Equal(err, nil)
	kmgTest.Equal(len(addrs), 4)
	kmgTest.Ok(addrs[0].IP.Equal(net.IPv4(127, 0, 0, 1)))
	//fmt.Println(addrs[0].IPNet.Mask.Size())
	one, _ := addrs[0].IPNet.Mask.Size()
	kmgTest.Equal(one, 8)
	kmgTest.Equal(addrs[0].DevString, "lo")

	kmgTest.Ok(addrs[1].IP.Equal(net.IPv4(192, 168, 1, 2)))
	kmgTest.Ok(addrs[1].IPNet.IP.Equal(net.IPv4(192, 168, 1, 0)))
	one, _ = addrs[1].IPNet.Mask.Size()
	kmgTest.Equal(one, 24)
	kmgTest.Equal(addrs[1].DevString, "eth0")

	kmgTest.Ok(addrs[2].IP.Equal(net.IPv4(172, 20, 8, 1)))
	one, _ = addrs[2].IPNet.Mask.Size()
	kmgTest.Equal(one, 32)
	kmgTest.Equal(addrs[2].DevString, "eth0")

	kmgTest.Ok(addrs[3].IP.Equal(net.IPv4(172, 20, 0, 1)))
	kmgTest.Equal(addrs[3].DevString, "ppp0")
}
