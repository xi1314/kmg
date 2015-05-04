package kmgNet

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgCmd"
	"net"
	"regexp"
)

//一个网络设备上面的地址
type DeviceAddr struct {
	IP        net.IP     //地址上面的ip
	IPNet     *net.IPNet //地址上面的子网
	DevString string     //设备名称 eth0 什么的
}

//目前仅支持linux
func (a DeviceAddr) IpAddrDel() (err error) {
	one, _ := a.IPNet.Mask.Size()
	return kmgCmd.CmdString(fmt.Sprintf("ip addr del %s/%d dev %s", a.IP.String(), one, a.DevString)).Run()
}

//目前仅支持linux
func GetCurrentDeviceAddr() (ipnets []DeviceAddr, err error) {
	out, err := kmgCmd.CmdString("ip addr").RunAndReturnOutput()
	if err != nil {
		return
	}
	return getCurrentDeviceAddrFromIPAddr(out)
}

func getCurrentDeviceAddrFromIPAddr(cmdReturn []byte) (ipnets []DeviceAddr, err error) {
	reg := regexp.MustCompile(`inet ([^ ]+).* ([^\s]+)`)
	out := reg.FindAllSubmatch(cmdReturn, -1)
	ipnets = make([]DeviceAddr, len(out))
	for i := range out {
		ip, ipnet, err := net.ParseCIDR(string(out[i][1]))
		if err != nil {
			return nil, err
		}
		ipnets[i] = DeviceAddr{
			IP:        ip,
			IPNet:     ipnet,
			DevString: string(out[i][2]),
		}
	}
	return
}
