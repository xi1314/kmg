package kmgNet

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

//从一个net.Listener里面读取需要Dial的地址(测试用的比较多)
func MustGetLocalAddrFromListener(listener net.Listener) string {
	return MustGetLocalAddrFromAddr(listener.Addr())
}

//从一个net.Listener里面读取需要Dial的地址(测试用的比较多)
func MustGetLocalAddrFromAddr(addr net.Addr) string {
	tcpAddr, err := net.ResolveTCPAddr(addr.Network(), addr.String())
	if err != nil {
		panic(err)
	}
	return "127.0.0.1:" + strconv.Itoa(tcpAddr.Port)
}

func MustTcpRandomListen() net.Listener {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	return l
}

func SpeedString(byteNum int, dur time.Duration) string {
	bytePerSecond := float64(byteNum) / (float64(dur) / float64(time.Second))
	if bytePerSecond > 1e9 {
		return fmt.Sprintf("%.2fGB/s", bytePerSecond/(1024*1024*1024))
	}
	if bytePerSecond > 1e6 {
		return fmt.Sprintf("%.2fMB/s", bytePerSecond/(1024*1024))
	}
	if bytePerSecond > 1e3 {
		return fmt.Sprintf("%.2fKB/s", bytePerSecond/1024)
	}
	return fmt.Sprintf("%.2fB/s", bytePerSecond)
}

func SizeString(byteNum int64) string {
	if byteNum > 1e9 || byteNum < -1e9{
		return fmt.Sprintf("%.2fGB", float64(byteNum)/(1024*1024*1024))
	}
	if byteNum > 1e6 || byteNum < -1e6{
		return fmt.Sprintf("%.2fMB", float64(byteNum)/(1024*1024))
	}
	if byteNum > 1e3 || byteNum < -1e3{
		return fmt.Sprintf("%.2fKB", float64(byteNum)/(1024))
	}
	return fmt.Sprintf("%dB", byteNum)
}
