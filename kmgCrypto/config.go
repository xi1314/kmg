package kmgCrypto

import "crypto/sha512"

// 这个作为当前应用的一个psk,而进行使用,开新项目的时候建议重新生成一个psk.
var DefaultPsk = [64]byte{0xe9, 0xf6, 0x6c, 0x4f, 0xa4, 0xee, 0x88, 0xc8}

// 可以在新项目开头注册一次.
// example:
//   kmgCrypto.SetDefaultPskFromString("4tLW/1FvSbwgc/mOrdtMSzcSYx7WWtI1Nn2uBJ5e/FXnW8XcPp9L45p/ahInsadGVF8Xsol1SnX4\nunlWzqAOUg==\n")
// use   python -c 'import os;print os.urandom(64).encode("base64")'  to get a new psk.
func SetDefaultPskFromString(s string) {
	DefaultPsk = sha512.Sum512([]byte(s))
	for _, f := range pskChangeCallbackList {
		f()
	}
}

//长度需要小于64
func GetPskFromDefaultPsk(length int, name string) []byte {
	psk := sha512.Sum512(append(DefaultPsk[:], []byte(name)...))
	return psk[:length]
}

var pskChangeCallbackList = []func(){}

func RegisterPskChangeCallback(f func()) {
	pskChangeCallbackList = append(pskChangeCallbackList, f)
}
