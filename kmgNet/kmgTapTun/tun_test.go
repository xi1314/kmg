package kmgTapTun

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgTest"
	"os"
	"testing"
)

func TestTun(ot *testing.T) {
	t := kmgTest.NewTestTools(ot)
	fmt.Println(1)
	tun, err := NewTun("")
	fmt.Println(2)
	if os.IsPermission(err) {
		ot.Skip("you need root permission to run this test.")
		return
	}
	t.Equal(err, nil)
	defer tun.Close()
	t.Equal(tun.GetDeviceType(), DeviceTypeTun)
	fmt.Println(3)

	err = SetP2PIpAndUp(tun, "10.209.34.1", "10.209.34.2")
	t.Equal(err, nil)

	err = SetMtu(tun, 1420)
	t.Equal(err, nil)
	fmt.Println(4)
	cmd := kmgCmd.CmdString("ping 10.209.34.2").GetExecCmd()
	err = cmd.Start()
	t.Equal(err, nil)
	defer cmd.Process.Kill()

	buf := make([]byte, 4096)
	n, err := tun.Read(buf)
	t.Equal(err, nil)
	t.Ok(n > 0)
	/*
		tun2, err := NewTun("")
		t.Equal(err, nil)
		defer tun2.Close()
	*/
}

func TestTap(ot *testing.T) {
	t := kmgTest.NewTestTools(ot)
	tap, err := NewTap("")
	if os.IsPermission(err) {
		ot.Skip("you need root permission to run this test.")
		return
	}
	t.Equal(err, nil)
	defer tap.Close()
	t.Equal(tap.GetDeviceType(), DeviceTypeTap)

	err = kmgCmd.CmdString("ifconfig " + tap.Name() + " 10.209.34.1 up").GetExecCmd().Run()
	t.Equal(err, nil)
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
