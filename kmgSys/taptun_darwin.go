package kmgSys

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
	"github.com/bronze1man/kmg/kmgNet"
)

// Interface is a TUN/TAP interface.
type iInterface struct {
	deviceType DeviceType
	file       io.ReadWriteCloser
	name       string
	isUtun     bool
}

// Create a new TUN interface whose name is ifName.
// you need install tuntaposx on your mac os,
// name should be something like tap0 .. tap15
// If ifName is empty, a default name (tap0, tap1, ... ) will be assigned.
func NewTap(ifName string) (ifce TunTapInterface, err error) {
	return newDeviceGeneric(ifName, DeviceTypeTap)
}

// Create a new TUN interface whose name is ifName.
// you need install tuntaposx on your mac os,
// name should be something like tun0 .. tun15
// If ifName is empty, a default name (tun0, tun1, ... ) will be assigned.
// If your put two empty name to this, you will get the same device.
func NewTun(ifName string) (ifce TunTapInterface, err error) {
	return newDeviceGeneric(ifName, DeviceTypeTun)
}

func NewTunNoName() (ifce TunTapInterface, err error) {
	for i := 0; i < 255; i++ {
		fd, err := utunOpenHelper(i)
		if err != nil {
			if kmgNet.IsResourceBusy(err){
				// 换一个tun的id号试试?
				continue
			}
			return nil, err
		}
		utunname := [20]byte{}
		utunname_len := uint32(20)
		err = SyscallGetSockopt(fd, SYSPROTO_CONTROL, UTUN_OPT_IFNAME, uintptr(unsafe.Pointer(&utunname[0])), &utunname_len)
		if err != nil {
			return nil, fmt.Errorf("Opening utun Error retrieving utun interface name %s", err.Error())
		}
		name := string(utunname[:int(utunname_len)-1])
		file := os.NewFile(uintptr(fd), name)
		return &iInterface{
			deviceType: DeviceTypeTun,
			//file: kmgIo.NewDebugRwc(file,"tun"),
			file:   file,
			name:   name,
			isUtun: true,
		}, nil
	}
	return nil, ErrAllDeviceBusy
}

///usr/include//sys/sys_domain.h
const SYSPROTO_CONTROL = 2
const AF_SYS_CONTROL = 2

// /usr/include/sys/kern_control.h
const CTLIOCGINFO = 0xc0644e03 //从宏计算出来的.
const MAX_KCTL_NAME = 96

type ctl_info struct {
	ctl_id   uint32
	ctl_name [MAX_KCTL_NAME]byte
}
type sockaddr_ctl struct {
	sc_len      byte
	sc_family   byte
	ss_sysaddr  uint16
	sc_id       uint32
	sc_unit     uint32
	sc_reserved [5]uint32
}

const UTUN_CONTROL_NAME = "com.apple.net.utun_control\x00" // /usr/include/net/if_utun.h
const UTUN_OPT_IFNAME = 2

func utunOpenHelper(num int) (fd int, err error) {
	fd, err = syscall.Socket(syscall.AF_SYSTEM, syscall.SOCK_DGRAM, SYSPROTO_CONTROL)
	if err != nil {
		return 0, fmt.Errorf("Opening utun socket(SYSPROTO_CONTROL) %s ", err.Error())
	}
	ctlInfo := ctl_info{}
	copy(ctlInfo.ctl_name[:], UTUN_CONTROL_NAME)
	err = SyscallIoctl(fd, CTLIOCGINFO, uintptr(unsafe.Pointer(&ctlInfo)))
	if err != nil {
		syscall.Close(fd)
		return 0, fmt.Errorf("Opening utun ioctl(CTLIOCGINFO) %s ", err.Error())
	}
	sa := sockaddr_ctl{
		sc_id:      ctlInfo.ctl_id,
		sc_family:  syscall.AF_SYSTEM,
		ss_sysaddr: AF_SYS_CONTROL,
		sc_unit:    uint32(num + 1),
	}
	sa.sc_len = byte(unsafe.Sizeof(sa))

	err = SyscallConnect(fd, uintptr(unsafe.Pointer(&sa)), uintptr(sa.sc_len))
	if err != nil {
		syscall.Close(fd)
		return 0, fmt.Errorf("Opening utun connect(AF_SYS_CONTROL) %s ", err.Error())
	}
	err = syscall.SetNonblock(fd, false)
	if err != nil {
		syscall.Close(fd)
		return 0, fmt.Errorf("Opening utun SetNonblock %s ", err.Error())
	}
	_, err = SyscallFcntl(fd, syscall.F_SETFD, syscall.FD_CLOEXEC)
	if err != nil {
		syscall.Close(fd)
		return 0, fmt.Errorf("Opening utun Fcntl %s ", err.Error())
	}
	return fd, nil
}

func newDeviceGeneric(ifName string, deviceType DeviceType) (ifce TunTapInterface, err error) {
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
func (ifce *iInterface) Write(buf []byte) (n int, err error) {
	inLen := len(buf)
	if ifce.isUtun {
		if buf[0]&0xf0 == 0x40 {
			buf = append([]byte{0, 0, 0, 2}, buf...)
		} else if buf[0]&0xf0 == 0x60 {
			buf = append([]byte{0, 0, 0, 0x1e}, buf...)
		} else {
			panic(fmt.Errorf("unexpect first byte %d", buf[0]))
		}
	}
	_, err = ifce.file.Write(buf)
	if err != nil {
		return 0, err
	}
	return inLen, nil
}

// Implement io.Reader interface.
func (ifce *iInterface) Read(p []byte) (nr int, err error) {
	nr, err = ifce.file.Read(p)
	if err != nil {
		return
	}
	if ifce.isUtun {
		copy(p[:nr-4], p[4:nr])
		return nr - 4, nil
	} else {
		return nr, nil
	}
}
func (ifce *iInterface) Close() (err error) {
	return ifce.file.Close()
}
