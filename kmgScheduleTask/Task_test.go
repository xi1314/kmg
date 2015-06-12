package kmgScheduleTask

import (
	"github.com/bronze1man/kmg/kmgSql"
	"github.com/bronze1man/kmg/kmgTest"
	"strconv"
	"sync"
	"testing"
	"time"
)

func init() {
	kmgSql.MustLoadTestConfig()

}

func TestRun(ot *testing.T) {
	kmgSql.MustSetTableDataYaml(`
kmgScheduleTask: []
`)
	tm := taskManager{}
	tm.pullDuration = 10 * time.Millisecond
	tm.Init()
	defer tm.Close()
	var a = 1
	c := make(chan struct{})
	tm.RegisterTaskFunc("abc", func(task Task) {
		a++
		kmgTest.Equal(task.InMap["a"], "1")
		c <- struct{}{}
	})
	for i := 0; i < 2; i++ {
		a = 1
		tm.RegisterTask(Task{
			FuncName:    "abc",
			ExecuteTime: time.Now().Add(time.Millisecond),
			InMap: map[string]string{
				"a": "1",
			},
		})
		kmgTest.Equal(a, 1)
		<-c
		kmgTest.Equal(a, 2)
	}
}

func TestRunRecoverFromDb(ot *testing.T) {
	kmgSql.MustSetTableDataYaml(`
kmgScheduleTask: []
`)
	tm := taskManager{}
	tm.pullDuration = 10 * time.Millisecond
	tm.Init()
	a := 1
	wg := sync.WaitGroup{}
	tm.RegisterTaskFunc("abc", func(task Task) {
		a = 2
		kmgTest.Equal(task.InMap["a"], "1")
		wg.Done()
	})
	wg.Add(1)
	tm.RegisterTask(Task{
		FuncName:    "abc",
		ExecuteTime: time.Now().Add(20 * time.Millisecond),
		InMap: map[string]string{
			"a": "1",
		},
	})
	tm.Close()
	rowList := kmgSql.MustGetAllInTable("kmgScheduleTask")
	kmgTest.Equal(len(rowList), 1)

	tm = taskManager{}
	tm.pullDuration = 10 * time.Millisecond
	tm.Init()
	tm.RegisterTaskFunc("abc", func(task Task) {
		a = 3
		kmgTest.Equal(task.InMap["a"], "1")
		wg.Done()
	})
	wg.Wait()

	kmgTest.Equal(a, 3)
	tm.RegisterTask(Task{
		FuncName:    "abc",
		ExecuteTime: time.Now().Add(30 * time.Millisecond),
		InMap: map[string]string{
			"a": "1",
		},
	})
	tm.Close()
	rowList = kmgSql.MustGetAllInTable("kmgScheduleTask")
	kmgTest.Equal(len(rowList), 1)

	a = 1
	time.Sleep(20 * time.Millisecond)
	wg.Add(1)
	tm = taskManager{}
	tm.pullDuration = time.Millisecond

	tm.RegisterTaskFunc("abc", func(task Task) {
		a = 4
		kmgTest.Equal(task.InMap["a"], "1")
		wg.Done()
	})
	tm.Init()
	time.Sleep(20 * time.Millisecond)

	wg.Wait()
	kmgTest.Equal(a, 4)
}

func TestHeap(ot *testing.T) {
	kmgSql.MustSetTableDataYaml(`
kmgScheduleTask: []
`)
	tm := taskManager{}
	tm.pullDuration = 10 * time.Millisecond
	tm.Init()
	a := 0
	c := make(chan struct{})
	tm.RegisterTaskFunc("abc", func(task Task) {
		kmgTest.Equal(task.InMap["a"], strconv.Itoa(a))
		a++
		c <- struct{}{}
	})

	for i := 0; i < 5; i++ {
		tm.RegisterTask(Task{
			FuncName:    "abc",
			ExecuteTime: time.Now().Add(time.Duration(i) * 3 * time.Millisecond),
			InMap: map[string]string{
				"a": strconv.Itoa(i),
			},
		})
	}

	for i := 0; i < 5; i++ {
		<-c
	}
}
