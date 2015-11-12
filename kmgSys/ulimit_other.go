// +build !linux

package kmgSys

func SetCurrentMaxFileNum(limit uint64) (err error) {
	return ErrPlatformNotSupport
}
