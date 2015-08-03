package kmgSys

func UsagePreInstall() {}

func Memory() (used float64, total int) {
	panic("not support os x")
	return float64(0), 0
}

func Cpu() (used float64, numOfCore int) {
	panic("not support os x")
	return float64(0), 0
}

func Disk() (used float64, total int) {
	panic("not support os x")
	return float64(0), 0
}

func FindDeviceNameByIp(ip string) string {
	panic("not support os x")
	return ""
}

func NetworkRXTX(deviceName string) (rx int, tx int) {
	panic("not support os x")
	return 0, 0
}
func NetworkConnection() (connectionCount int) {
	panic("not support os x")
	return 0
}

func IKEUserCount() int {
	panic("not support os x")
	return 0
}
