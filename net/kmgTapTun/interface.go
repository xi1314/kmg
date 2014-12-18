package kmgTapTun

import (
	"errors"
	"io"
)

var ErrPlatformNotSupport = errors.New("tun/tap: platform is not support")

// Interface is a TUN/TAP interface.
type Interface interface {
	IsTUN() bool
	IsTAP() bool
	Name() string
	io.Writer
	io.Reader
	io.Closer
}
