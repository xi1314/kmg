package kmgSys

import (
	"fmt"
	"os"
	"testing"

	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgTest"
)

func TestTun(ot *testing.T) {
	fmt.Println(1)
	tun, err := NewTun("")
	fmt.Println(2)
	if os.IsPermission(err) {
		ot.Skip("you need root permission to run this test.")
		return
	}
	kmgTest.Equal(err, nil)
	defer tun.Close()
	kmgTest.Equal(tun.GetDeviceType(), DeviceTypeTun)
	fmt.Println(3)

	err = SetP2PIpAndUp(tun, "10.209.34.1", "10.209.34.2")
	kmgTest.Equal(err, nil)

	err = SetMtu(tun, 1420)
	kmgTest.Equal(err, nil)
	fmt.Println(4)
	cmd := kmgCmd.CmdString("ping 10.209.34.2").GetExecCmd()
	err = cmd.Start()
	kmgTest.Equal(err, nil)
	defer cmd.Process.Kill()

	buf := make([]byte, 4096)
	n, err := tun.Read(buf)
	kmgTest.Equal(err, nil)
	kmgTest.Ok(n > 0)
	/*
		tun2, err := NewTun("")
		t.Equal(err, nil)
		defer tun2.Close()
	*/
}

func TestTap(ot *testing.T) {
	tap, err := NewTap("")
	if os.IsPermission(err) {
		ot.Skip("you need root permission to run this test.")
		return
	}
	kmgTest.Equal(err, nil)
	defer tap.Close()
	kmgTest.Equal(tap.GetDeviceType(), DeviceTypeTap)

	err = kmgCmd.CmdString("ifconfig " + tap.Name() + " 10.209.34.1 up").GetExecCmd().Run()
	kmgTest.Equal(err, nil)
	/*
		cmd := kmgCmd.NewOsStdioCmdString("ping 10.0.0.2")
		err = cmd.Start()
		t.Equal(err, nil)
		defer cmd.Process.Kill()

		buf := make([]byte, 4096)
		n, err := tun.Read(buf)
		t.Equal(err, nil)
		t.Ok(n > 0)

		tun2, err := NewTap("")
		t.Equal(err, nil)
		defer tun2.Close()
	*/
}
