// +build !linux,!darwin

package kmgSys

// Create a new TAP interface whose name is ifName.
// If ifName is empty, a default name (tap0, tap1, ... ) will be assigned.
// ifName should not exceed 16 bytes.
func NewTap(ifName string) (ifce TunTapInterface, err error) {
	return nil, ErrPlatformNotSupport
}

// Create a new TUN interface whose name is ifName.
// If ifName is empty, a default name (tap0, tap1, ... ) will be assigned.
// ifName should not exceed 16 bytes.
func NewTun(ifName string) (ifce TunTapInterface, err error) {
	return nil, ErrPlatformNotSupport
}

func NewTunNoName() (ifce TunTapInterface, err error) {
	return nil, ErrPlatformNotSupport
}
