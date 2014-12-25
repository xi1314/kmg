package kmgIo

import (
	"fmt"
	"time"
)

func FmtByteSpeed(byteNum int, dur time.Duration) string {
	bytePerSecond := float64(byteNum) / (float64(dur) / float64(time.Second))
	if bytePerSecond > 1e9 {
		return fmt.Sprintf("%.2fGB/s", bytePerSecond/(1024*1024*1024))
	}
	if bytePerSecond > 1e6 {
		return fmt.Sprintf("%.2fMB/s", bytePerSecond/(1024*1024))
	}
	if bytePerSecond > 1e3 {
		return fmt.Sprintf("%.2fKB/s", bytePerSecond/1024)
	}
	return fmt.Sprintf("%.2fB/s", bytePerSecond)
}

func FmtByteNum(byteNum int) string {
	if byteNum > 1e9 {
		return fmt.Sprintf("%.2fGB", float64(byteNum)/(1024*1024*1024))
	}
	if byteNum > 1e6 {
		return fmt.Sprintf("%.2fMB", float64(byteNum)/(1024*1024))
	}
	if byteNum > 1e3 {
		return fmt.Sprintf("%.2fKB", float64(byteNum)/(1024))
	}
	return fmt.Sprintf("%dB", byteNum)
}
