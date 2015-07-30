package kmgDebug

import (
	"fmt"
	"sync"
)

type RWMutex struct {
	rw   sync.RWMutex
	Name string
}

func (rw *RWMutex) RLock() {
	fmt.Println("[RWMutex]", rw.Name, "RLock start")
	rw.rw.RLock()
	fmt.Println("[RWMutex]", rw.Name, "RLock end")
}
func (rw *RWMutex) RUnlock() {
	fmt.Println("[RWMutex]", rw.Name, "RUnlock start")
	rw.rw.RUnlock()
	fmt.Println("[RWMutex]", rw.Name, "RUnlock end")
}
func (rw *RWMutex) Lock() {
	fmt.Println("[RWMutex]", rw.Name, "Lock start")
	rw.rw.Lock()
	fmt.Println("[RWMutex]", rw.Name, "Lock end")
}
func (rw *RWMutex) Unlock() {
	fmt.Println("[RWMutex]", rw.Name, "Unlock start")
	rw.rw.Unlock()
	fmt.Println("[RWMutex]", rw.Name, "Unlock end")
}

type UselessRWMutex struct{}

func (rw UselessRWMutex) RLock() {
}
func (rw UselessRWMutex) RUnlock() {
}
func (rw UselessRWMutex) Lock() {
}
func (rw UselessRWMutex) Unlock() {
}
