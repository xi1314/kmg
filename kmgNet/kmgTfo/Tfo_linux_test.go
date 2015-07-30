package kmgTfo

import (
	"net"
	"testing"

	"github.com/bronze1man/kmg/kmgNet/netTester"
)

func TestTfo(ot *testing.T) {
	netTester.RunTcpListenerDialerTest(func() (net.Listener, error) {
		return TfoListen("tcp", "127.0.0.1:50000")
	}, func() (net.Conn, error) {
		return TfoLazyDial("tcp", "127.0.0.1:50000")
	}, false)
}
