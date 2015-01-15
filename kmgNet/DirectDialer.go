package kmgNet

import (
	"net"
)

type DirectDialer interface {
	DirectDial() (net.Conn, error)
}

type DirectDialerFunc func() (net.Conn, error)

func (f DirectDialerFunc) Dial() (net.Conn, error) {
	return f()
}
