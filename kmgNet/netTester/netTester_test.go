package netTester

import (
	"github.com/bronze1man/kmg/kmgTest"
	"net"
	"testing"
)

func TestRunTcpListenerDialerTest(ot *testing.T) {
	t := kmgTest.NewTestTools(ot)
	l, err := net.Listen("tcp", "127.0.0.1:50000")
	t.Equal(err, nil)
	RunTcpListenerDialerTest(l, func() (net.Conn, error) {
		return net.Dial("tcp", "127.0.0.1:50000")
	})
}
