package kmgPlatform

import "runtime"

type Platform struct {
	Os   string
	Arch string
}

func (p Platform) Compatible(other Platform) bool {
	return p == other
}

func (p Platform) String() string {
	return p.Os + "_" + p.Arch
}

func (p Platform) GetExeSuffix() string {
	if p.Os == "windows" {
		return p.Os + "_" + p.Arch + ".exe"
	}
	return p.Os + "_" + p.Arch
}

var LinuxAmd64 = Platform{Os: "linux", Arch: "amd64"}
var DarwinAmd64 = Platform{Os: "darwin", Arch: "amd64"}
var WindowsAmd64 = Platform{Os: "windows", Arch: "amd64"}

//编译这个软件的平台
func GetCompiledPlatform() Platform {
	return Platform{
		Os:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}
}

func IsLinux() bool {
	return runtime.GOOS == "linux"
}

func IsDarwin() bool {
	return runtime.GOOS == "darwin"
}

func IsWindows() bool {
	return runtime.GOOS == "windows"
}
