package kmgProcessMutex

import "os"

type FileMutex struct {
	isOwner  bool
	filePath string
	Name     string
	f        *os.File
}

func (fm *FileMutex) Lock() {
	panic("Not Support for Windows")
}

func (fm *FileMutex) UnLock() {
	panic("Not Support for Windows")
}
