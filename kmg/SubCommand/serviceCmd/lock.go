package serviceCmd

import (
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgRand"
	"path/filepath"
	"time"
)

type FileMutex struct {
	isOwner  bool
	filePath string
}

func (fm *FileMutex) Lock(name string) {
	fm.filePath = filepath.Join("/tmp", name)
	if !fm.isLock() {
		fm.ownAndLock()
		return
	}
	if fm.isOwner {
		return
	} else {
		fm.waitToUnLock()
		fm.ownAndLock()
	}
}

func (fm *FileMutex) UnLock() {
	if !fm.isOwner {
		return
	}
	fm.isOwner = false
	kmgFile.MustDelete(fm.filePath)
}

func (fm *FileMutex) waitToUnLock() {
	i := 0
	for {
		i++
		if !fm.isLock() {
			break
		} else {
			time.Sleep(time.Second)
		}
		if i > 60 {
			panic("FileMutex wait for UnLock timeout")
		}
	}
}

func (fm *FileMutex) isLock() bool {
	time.Sleep(time.Duration(kmgRand.IntBetween(0, 1000)) * time.Millisecond)
	if kmgFile.MustFileExist(fm.filePath) {
		return true
	}
	return false
}

func (fm *FileMutex) ownAndLock() {
	kmgFile.MustWriteFile(fm.filePath, []byte("1"))
	fm.isOwner = true
}
