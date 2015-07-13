package kmgProxy

import (
	"encoding/binary"
	"fmt"
	"github.com/bronze1man/kmg/kmgNet"
	"io"
	"net"
	"time"
)

func Socks4aDialTimeout(proxyAddr string, targetAddr string, timeout time.Duration) (conn net.Conn, err error) {
	conn, err = net.DialTimeout("tcp", proxyAddr, timeout)
	if err != nil {
		return
	}
	err = socks4aConnect(conn, targetAddr)
	if err != nil {
		conn.Close()
		return
	}
	return
}
func Socks4aDial(proxyAddr string, targetAddr string) (conn net.Conn, err error) {
	return Socks4aDialTimeout(proxyAddr, targetAddr, 0)
}

func NewSocks4aDialer(proxyAddr string) kmgNet.Dialer {
	return func(network, address string) (net.Conn, error) {
		return Socks4aDialTimeout(proxyAddr, address, 0)
	}
}

//按照socks4a接口,连接到那个地址去,后面可以当成这个连接已经连接到了对方主机,本实现没有用户密码功能
func socks4aConnect(conn net.Conn, addr string) (err error) {
	uaddr, err := ParseUnsolvedAddr(addr)
	if err != nil {
		return
	}
	toWriteHeader := []byte{4, 1, 0, 0, 0, 0, 0, 1, 0}
	binary.BigEndian.PutUint16(toWriteHeader[2:4], uint16(uaddr.Port))
	if uaddr.Ip != nil {
		ip := uaddr.Ip.To4()
		if ip == nil {
			return fmt.Errorf("you can not connect to a ipv6 addr with socks4a")
		}
		copy(toWriteHeader[4:8], []byte(ip[:4]))
	} else {
		toWriteHeader = append(toWriteHeader, []byte(uaddr.Domain)...)
		toWriteHeader = append(toWriteHeader, 0)
	}
	_, err = conn.Write(toWriteHeader)
	if err != nil {
		return
	}
	_, err = io.ReadFull(conn, toWriteHeader[:8])
	if err != nil {
		return
	}
	switch toWriteHeader[1] {
	case 0x5a:
	//ok
	case 0x5b:
		return fmt.Errorf("socks4a: request rejected or failed")
	case 0x5c:
		return fmt.Errorf("socks4a: request failed because client is not running identd (or not reachable from the server)")
	case 0x5d:
		return fmt.Errorf("socks4a: request failed because client is not running identd (or not reachable from the server)")
	default:
		return fmt.Errorf("socks4a: protoal error 1 %d", toWriteHeader[1])
	}
	return
}
