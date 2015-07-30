package kmgSys

func UsagePreInstall() {}

func Memory() (used float64, total int) {
	return float64(0), 0
}

func Cpu() (used float64, numOfCore int) {
	return float64(0), 0
}

func Disk() (used float64, total int) {
	return float64(0), 0
}

func NetworkRXTX(deviceName string) (rx int, tx int) {
	return 0, 0
}
func NetworkConnection() (connectionCount int) {
	return 0
}

func IKEUserCount() int {
	return 0
}
