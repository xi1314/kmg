package kmgNet_test

import (
	"github.com/bronze1man/kmg/kmgNet"
	"github.com/bronze1man/kmg/kmgNet/netTester"
	"net"
	"testing"
)

func TestTfo(ot *testing.T) {
	netTester.RunTcpListenerDialerTest(func() (net.Listener, error) {
		return kmgNet.TfoListen("tcp", "127.0.0.1:50000")
	}, func() (net.Conn, error) {
		return kmgNet.TfoLazyDial("tcp", "127.0.0.1:50000")
	}, false)
}
