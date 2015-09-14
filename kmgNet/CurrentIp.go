package kmgNet

import (
	"fmt"
	"net"
	"regexp"
	"strconv"

	"github.com/bronze1man/kmg/errors"
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgPlatform"
)

// 获取当前地址的列表的字符串
func GetCurrentAddrListString() string{
	addrList,err:=net.InterfaceAddrs()
	if err!=nil{
		return "[GetCurrentAddrListString] fail "+err.Error()
	}
	out:=""
	for _,addr:=range addrList{
		out+=addr.String()+";"
	}
	return out
}

//一个网络设备上面的地址
type DeviceAddr struct {
	IP        net.IP     //地址上面的ip
	IPNet     *net.IPNet //地址上面的子网
	DevString string     //设备名称 eth0 什么的
}

// 返回nil表示没找到这个设备,或者这个设备上面没有ip
func MustGetFirstIPByNetDeviceName(devname string) net.IP{
	if !kmgPlatform.IsLinux(){
		panic("[GetFirstIPByNetDeviceName] only support linux now")
	}
	deviceAddrList, err := GetCurrentDeviceAddr()
	if err != nil {
		panic(err)
	}
	for _,deviceAddr:=range deviceAddrList{
		if deviceAddr.DevString==devname{
			return deviceAddr.IP
		}
	}
	return nil
}

//目前仅支持linux
func (a DeviceAddr) IpAddrDel() (err error) {
	if !kmgPlatform.IsLinux(){
		panic("[DeviceAddr.IpAddrDel] only support linux now")
	}
	one, _ := a.IPNet.Mask.Size()
	return kmgCmd.CmdString(fmt.Sprintf("ip addr del %s/%d dev %s", a.IP.String(), one, a.DevString)).Run()
}

//目前仅支持linux
func GetCurrentDeviceAddr() (ipnets []DeviceAddr, err error) {
	if !kmgPlatform.IsLinux(){
		panic("[GetCurrentDeviceAddr] only support linux now")
	}
	out, err := kmgCmd.CmdString("ip addr").RunAndReturnOutput()
	if err != nil {
		return
	}
	return getCurrentDeviceAddrFromIPAddr(out)
}

//返回当前机器上面的所有ip列表.没有ip会报错
func MustGetCurrentIpList() (ipList []net.IP) {
	if !kmgPlatform.IsLinux(){
		panic("[MustGetCurrentIpList] only support linux now")
	}
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
	if !kmgPlatform.IsLinux(){
		panic("[MustGetCurrentIpWithPortList] only support linux now")
	}
	deviceAddrList, err := GetCurrentDeviceAddr()
	if err != nil {
		panic(err)
	}
	if len(deviceAddrList) == 0 {
		panic(errors.New("[MustGetCurrentIpList] do not find any ip address."))
	}
	sList = make([]string, 0, len(deviceAddrList))
	sPort := strconv.Itoa(int(port))
	for _, addr := range deviceAddrList {
		ones, size := addr.IPNet.Mask.Size()
		if ones == size { // 实践表明 不能监听这种子网只有一个ip的地址.
			continue
		}
		sList = append(sList, net.JoinHostPort(addr.IP.String(), sPort))
	}
	return sList
}

func getCurrentDeviceAddrFromIPAddr(cmdReturn []byte) (ipnets []DeviceAddr, err error) {
	//可能性1 本地回路     inet 127.0.0.1/8 scope host lo
	//可能性2 物理网卡     inet 10.169.224.99/21 brd 10.169.231.255 scope global eth0
	//可能性3 pptp虚拟网卡 inet 172.20.0.1 peer 172.20.0.2/32 scope global ppp0
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
			ip, ipnet, err = net.ParseCIDR(string(out[i][1]) + "/32")
			if err != nil {
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
