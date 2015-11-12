package kmgSys

import (
	"github.com/bronze1man/kmg/kmgStrings"
	"os"
	"strings"
)

// 确保PATH里面包含 /usr/local/bin 和 /bin
func RecoverPath() {
	pathenv := os.Getenv("PATH")
	pathList := strings.Split(pathenv, ":")
	change := false
	if !kmgStrings.IsInSlice(pathList, "/usr/local/bin") {
		change = true
		pathenv += ":/usr/local/bin"
	}
	if !kmgStrings.IsInSlice(pathList, "/bin") {
		change = true
		pathenv += ":/bin"
	}
	if change {
		os.Setenv("PATH", pathenv)
	}
}
