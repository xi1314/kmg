package kmgNet

import "net"

type ConnWrapper func(conn net.Conn) (net.Conn, error)

func ListenerConnWrapper(l net.Listener, connWrapper ConnWrapper) net.Listener {
	return &listenerConnWrapper{
		Listener:    l,
		connWrapper: connWrapper,
	}
}

type listenerConnWrapper struct {
	net.Listener
	connWrapper ConnWrapper
}

func (l listenerConnWrapper) Accept() (c net.Conn, err error) {
	c, err = l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	c, err = l.connWrapper(c)
	return
}
