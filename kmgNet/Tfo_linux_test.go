package kmgNet_test

import (
	"github.com/bronze1man/kmg/kmgNet"
	"github.com/bronze1man/kmg/kmgNet/netTester"
	"github.com/bronze1man/kmg/kmgTest"
	"net"
	"testing"
)

func TestTfo(ot *testing.T) {
	t := kmgTest.NewTestTools(ot)
	l, err := kmgNet.TfoListen("tcp", "127.0.0.1:50000")
	t.Equal(err, nil)
	netTester.RunTcpListenerDialerTest(l, func() (net.Conn, error) {
		return kmgNet.TfoLazyDial("tcp", "127.0.0.1:50000")
	})
}
