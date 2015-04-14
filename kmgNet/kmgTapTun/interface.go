package kmgTapTun

import (
	"errors"
	"github.com/bronze1man/kmg/kmgCmd"
	"io"
	"runtime"
	"strconv"
)

type DeviceType string

func (s DeviceType) String() string {
	return string(s)
}

var DeviceTypeTap DeviceType = "tap"
var DeviceTypeTun DeviceType = "tun"

var ErrPlatformNotSupport = errors.New("tun/tap: platform is not support")
var ErrAllDeviceBusy = errors.New("tun/tap: all dev is busy.")

// Interface is a TUN/TAP interface.
type Interface interface {
	io.ReadWriteCloser
	GetDeviceType() DeviceType
	Name() string
}

//set tun p2p ip and up this device
func SetP2PIpAndUp(ifac Interface, srcIp string, destIp string) error {
	switch runtime.GOOS {
	case "darwin":
		return kmgCmd.StdioSliceRun("ifconfig", ifac.Name(), srcIp, destIp, "up")
	case "linux":
		return kmgCmd.StdioSliceRun("ifconfig", ifac.Name(), srcIp, "pointopoint", destIp, "up")
	default:
		return ErrPlatformNotSupport
	}
}

//set mtu on a device
func SetMtu(ifac Interface, mtu int) error {
	return kmgCmd.StdioSliceRun("ifconfig", ifac.Name(), "mtu", strconv.Itoa(mtu))
}
