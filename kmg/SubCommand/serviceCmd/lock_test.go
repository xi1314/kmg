package serviceCmd_test

import (
	"fmt"
	"github.com/bronze1man/kmg/kmg/SubCommand/serviceCmd"
	"testing"
	"time"
	"sync"
)


//期望输出结果应该是：
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
			l := &serviceCmd.FileMutex{}
			l.Lock("abc")
			fmt.Println("Lock", i)
			time.Sleep(time.Second)
			fmt.Println(i)
			l.UnLock()
			fmt.Println("UnLock", i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
