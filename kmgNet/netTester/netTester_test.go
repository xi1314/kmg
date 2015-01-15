package netTester

import (
	"net"
	"testing"
)

func TestRunTcpListenerDialerTest(ot *testing.T) {
	RunTcpListenerDialerTest(func() (net.Listener, error) {
		return net.Listen("tcp", "127.0.0.1:50000")
	}, func() (net.Conn, error) {
		return net.Dial("tcp", "127.0.0.1:50000")
	}, false)
}

func TestRunTcpTestWithNetDialAndNetListener(ot *testing.T) {
	RunTcpTestWithNetDialAndNetListener("127.0.0.1:50000", "127.0.0.1:50000", false)
}
