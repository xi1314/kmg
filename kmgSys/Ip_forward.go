package kmgSys

import (
	"bytes"
	"fmt"

	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgPlatform"
)

func IsIpForwardOn() bool {
	if !kmgPlatform.IsLinux(){
		panic("[IsIpForwardOn] only support linux now")
	}
	b := kmgFile.MustReadFile("/proc/sys/net/ipv4/ip_forward")
	if bytes.Equal(b, []byte{'0'}) {
		return false
	}
	if bytes.Equal(b, []byte{'1'}) {
		return true
	}
	panic(fmt.Errorf("[IsIpForwardOn] unable to understand info in /proc/sys/net/ipv4/ip_forward %#v", b))
}

// 证实可用
func SetIpForwardOn() {
	if !kmgPlatform.IsLinux(){
		panic("[SetIpForwardOn] only support linux now")
	}
	kmgFile.MustWriteFile("/proc/sys/net/ipv4/ip_forward", []byte("1"))
	// 已经证实,多次写入不会出现任何问题.
	// TODO 正确解析/etc/sysctl.conf 如果后面又加一条 = 0 估计就挂了.
	if !bytes.Contains(kmgFile.MustReadFile("/etc/sysctl.conf"), []byte("\nnet.ipv4.ip_forward = 1")) {
		kmgFile.MustAppendFile("/etc/sysctl.conf", []byte("\nnet.ipv4.ip_forward = 1"))
	}
}
