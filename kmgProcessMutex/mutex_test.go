package kmgProcessMutex_test

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgProcessMutex"
	"sync"
	"testing"
)

//期望输出结果应该是 Lock N 和 UnLock N 成对出现，一对 Lock 和 UnLock 之间是的操作是原子的，不会 Data Race：
//Lock n0
//n0
//UnLock n0
//Lock n9
//n9
//UnLock n9
//...
func TestFileMutex(t *testing.T) {
	wg := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			l := &kmgProcessMutex.FileMutex{Name: "abc"}
			l.Lock()
			fmt.Println("Lock", i)
			l.UnLock()
			fmt.Println("UnLock", i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
