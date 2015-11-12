// +build linux darwin

package kmgProcessMutex

import (
	"github.com/bronze1man/kmg/kmgErr"
	"os"
	"path/filepath"
	"syscall"
)

type FileMutex struct {
	isOwner  bool
	filePath string
	Name     string
	f        *os.File
}

// flock 加锁，会将锁（系统级别）挂在某个文件上，只要有进程给某个文件挂上了锁，则其他进程（包括本进程）就必须解锁
func (fm *FileMutex) Lock() {
	kmgErr.PanicIfError(syscall.Flock(int(fm.getFd()), int(syscall.LOCK_EX)))
}

// flock 加的锁，有两种方式解除锁：
// 1.加锁的进程退出了，锁自动释放
// 2.在进程内，使用 fd0 挂上的锁，显式的使用 fd0 解锁；
// 注意：同一个进程内，fd0 和 fd1 指向同一个文件（fd1 不是 fd0 的副本），若是 fd0 挂上的锁，必须用 fd0 来解锁，fd1 无法解锁
func (fm *FileMutex) UnLock() {
	kmgErr.PanicIfError(syscall.Flock(fm.getFd(), int(syscall.LOCK_UN)))
}

func (fm *FileMutex) getFd() int {
	if fm.f != nil {
		return int(fm.f.Fd())
	}
	if fm.Name == "" {
		panic(`no specialed file [Name] for locking，Example：l := &kmgProcessMutex.FileMutex{Name: "abc"}`)
	}
	fm.filePath = filepath.Join("/tmp", fm.Name)
	var err error
	fm.f, err = os.OpenFile(fm.filePath, os.O_CREATE|os.O_EXCL, os.FileMode(0777))
	if err == nil {
		return int(fm.f.Fd())
	}
	if os.IsExist(err) {
		fm.f, err = os.OpenFile(fm.filePath, os.O_RDONLY, os.FileMode(0777))
	} else {
		kmgErr.LogErrorWithStack(err)
	}
	return int(fm.f.Fd())
}
