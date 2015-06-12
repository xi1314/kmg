package kmgNet

import (
	"fmt"
	"github.com/bronze1man/kmg/errors"
	"github.com/bronze1man/kmg/kmgCmd"
	"net"
	"regexp"
	"strconv"
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

//返回当前机器上面的所有ip列表.没有ip会报错
func MustGetCurrentIpList() (ipList []net.IP) {
	deviceAddrList, err := GetCurrentDeviceAddr()
	if err != nil {
		panic(err)
	}
	if len(deviceAddrList) == 0 {
		panic(errors.New("[MustGetCurrentIpList] do not find any ip address."))
	}
	ipList = make([]net.IP, len(deviceAddrList))
	for i, addr := range deviceAddrList {
		ipList[i] = addr.IP
	}
	return ipList
}

func MustGetCurrentIpWithPortList(port uint16) (sList []string) {
	ipList := MustGetCurrentIpList()
	sList = make([]string, len(ipList))
	sPort := strconv.Itoa(int(port))
	for i, ip := range ipList {
		sList[i] = ip.String() + ":" + sPort
	}
	return sList
}

func getCurrentDeviceAddrFromIPAddr(cmdReturn []byte) (ipnets []DeviceAddr, err error) {
	//可能性1 inet 127.0.0.1/8 scope host lo
	//可能性2 inet 10.169.224.99/21 brd 10.169.231.255 scope global eth0
	//可能性3 inet 172.20.0.1 peer 172.20.0.2/32 scope global ppp0
	reg := regexp.MustCompile(`inet ([^ ]+).* ([^\s]+)`)
	out := reg.FindAllSubmatch(cmdReturn, -1)
	ipnets = make([]DeviceAddr, len(out))
	for i := range out {
		ip, ipnet, err := net.ParseCIDR(string(out[i][1]))
		if err != nil {
			_, ok := err.(*net.ParseError)
			if !ok {
				return nil, err
			}
			ip = net.ParseIP(string(out[i][1]))
			if ip == nil {
				return nil, fmt.Errorf("[getCurrentDeviceAddrFromIPAddr] can not parse CIDR or IP [%s]", out[i][0])
			}
		}
		ipnets[i] = DeviceAddr{
			IP:        ip,
			IPNet:     ipnet,
			DevString: string(out[i][2]),
		}
	}
	return
}
