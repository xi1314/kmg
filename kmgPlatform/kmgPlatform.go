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
	return p.Os + "-" + p.Arch
}

var LinuxAmd64 = Platform{Os: "linux", Arch: "amd64"}
var DarwinAmd64 = Platform{Os: "darwin", Arch: "amd64"}

//编译这个软件的平台
func GetCompiledPlatform() Platform {
	return Platform{
		Os:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}
}