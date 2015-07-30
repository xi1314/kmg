package kmgTapTun

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Interface is a TUN/TAP interface.
type iInterface struct {
	deviceType DeviceType
	file       *os.File
	name       string
}

// Create a new TUN interface whose name is ifName.
// you need install tuntaposx on your mac os,
// name should be something like tap0 .. tap15
// If ifName is empty, a default name (tap0, tap1, ... ) will be assigned.
func NewTap(ifName string) (ifce Interface, err error) {
	return newDevice(ifName, DeviceTypeTap)
}

// Create a new TUN interface whose name is ifName.
// you need install tuntaposx on your mac os,
// name should be something like tun0 .. tun15
// If ifName is empty, a default name (tun0, tun1, ... ) will be assigned.
// If your put two empty name to this, you will get the same device.
func NewTun(ifName string) (ifce Interface, err error) {
	return newDevice(ifName, DeviceTypeTun)
}

func newDevice(ifName string, deviceType DeviceType) (ifce Interface, err error) {
	iifce := &iInterface{deviceType: deviceType, name: ifName}

	devTypeString := deviceType.String()
	if ifName != "" {
		if !strings.HasPrefix(ifName, devTypeString) {
			return nil, fmt.Errorf("name should look like %s0 .", devTypeString)
		}
		iifce.file, err = os.OpenFile("/dev/"+ifName, os.O_RDWR, 0)
		if err != nil {
			return nil, err
		}
		return iifce, nil
	} else {
		//一个一个尝试
		for i := 0; i <= 15; i++ {
			iifce.name = devTypeString + strconv.Itoa(i)
			fmt.Println("Open tun ", iifce.name, "start")
			iifce.file, err = os.OpenFile("/dev/"+iifce.name, os.O_RDWR, 0) // 这个调用有时候会卡死,换一个名字往往能好.
			fmt.Println("Open tun ", iifce.name, "finish")

			if err == nil {
				return iifce, nil
			}
			if err != nil {
				if strings.Contains(err.Error(), "resource busy") {
					continue
				}
				return nil, err
			}
		}
		return nil, ErrAllDeviceBusy
	}
}

func (ifce *iInterface) GetDeviceType() DeviceType {
	return ifce.deviceType
}

// Returns the interface name of ifce, e.g. tun0, tap1, etc..
func (ifce *iInterface) Name() string {
	return ifce.name
}

// Implement io.Writer interface.
func (ifce *iInterface) Write(p []byte) (n int, err error) {
	return ifce.file.Write(p)
}

// Implement io.Reader interface.
func (ifce *iInterface) Read(p []byte) (n int, err error) {
	return ifce.file.Read(p)
}
func (ifce *iInterface) Close() (err error) {
	return ifce.file.Close()
}
