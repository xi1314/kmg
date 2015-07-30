package kmgSys

import (
	"bytes"
	"github.com/bronze1man/kmg/kmgCmd"
	"net"
	"runtime"
	"strconv"
)

type SetP2PIpRequest struct {
	IfaceName string // 必填
	SrcIp     net.IP // 必填
	DstIp     net.IP // 必填
	Mtu       int
	Mask      net.IPMask
}

//set tun p2p ip and up this device
// mtu default to 1500
func SetP2PIpAndUp(req SetP2PIpRequest) error {
	switch runtime.GOOS {
	case "darwin":
		cmdSlice := []string{"ifconfig", req.IfaceName, req.SrcIp.String(), req.DstIp.String()}
		if req.Mask != nil {
			cmdSlice = append(cmdSlice, "netmask", netmaskDotListString(req.Mask))
		}
		if req.Mtu > 0 {
			cmdSlice = append(cmdSlice, "mtu", strconv.Itoa(req.Mtu))
		}
		cmdSlice = append(cmdSlice, "up")
		return kmgCmd.StdioSliceRun(cmdSlice)
	case "linux":
		cmdSlice := []string{"ifconfig", req.IfaceName, req.SrcIp.String(), "pointopoint", req.DstIp.String()}
		if req.Mask != nil {
			cmdSlice = append(cmdSlice, "netmask", netmaskDotListString(req.Mask))
		}
		if req.Mtu > 0 {
			cmdSlice = append(cmdSlice, "mtu", strconv.Itoa(req.Mtu))
		}
		cmdSlice = append(cmdSlice, "up")
		return kmgCmd.StdioSliceRun(cmdSlice)
	default:
		return ErrPlatformNotSupport
	}
}

func netmaskDotListString(mask net.IPMask) string {
	buf := &bytes.Buffer{}
	for i, b := range mask {
		buf.WriteString(strconv.Itoa(int(b)))
		if i != len(mask)-1 {
			buf.WriteString(".")
		}
	}
	return buf.String()
}
