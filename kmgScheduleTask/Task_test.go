package kmgScheduleTask

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgSql"
	"github.com/bronze1man/kmg/kmgTime"
	"testing"
	"time"
)

func TestSelectUsage(t *testing.T) {
	c := make(chan int)
	go foo(c)
	select {
	case <-c:
		fmt.Println("Task foo Done!")
	case <-time.After(5 * time.Second):
		fmt.Println("Timed out!")
	}
}

func foo(ch chan int) {
	time.Sleep(10 * time.Second)
	ch <- 1
}

func TestRun(t *testing.T) {
	kmgTime.NowTime = kmgTime.NewFixedNower(kmgTime.MustFromMysqlFormatDefaultTZ("2000-01-15 00:00:00"))
	RegisterTask(Task{
		ExecuteTime: kmgTime.MustFromMysqlFormatDefaultTZ("2000-06-15 00:00:00"),
	})
	RegisterTask(Task{
		ExecuteTime: kmgTime.MustFromMysqlFormatDefaultTZ("1999-06-15 00:00:00"),
		IsDone:      true,
	})
	RegisterTask(Task{
		ExecuteTime:           kmgTime.MustFromMysqlFormatDefaultTZ("1999-06-15 00:00:00"),
		MaxAllowNumberOfRetry: 3,
		Func: func(isDone chan bool, task Task) {
			time.Sleep(2 * time.Second)
			fmt.Println("hi task")
			isDone <- false
		},
	})
	RegisterTask(Task{
		ExecuteTime: kmgTime.MustFromMysqlFormatDefaultTZ("2999-06-15 00:00:00"),
		IsDone:      true,
	})
	run()
}

func TestStart(t *testing.T) {
	kmgTime.NowTime = kmgTime.NewFixedNower(kmgTime.MustFromMysqlFormatDefaultTZ("2000-01-15 00:00:00"))
	RegisterTask(Task{
		ExecuteTime:           kmgTime.MustFromMysqlFormatDefaultTZ("1999-06-15 00:00:00"),
		MaxAllowNumberOfRetry: 3,
		Func: FooTaskFunc,
	})
	RegisterTaskFunc(FooTaskFunc)
	Start()
	fmt.Println("Start Done!")
}

func FooTaskFunc(isDone chan bool, task Task) {
	time.Sleep(2 * time.Second)
	fmt.Println("hi task")
	isDone <- false
}

func TestReflectFunc(t *testing.T) {
	RegisterTaskFunc(boo)
	for n, f := range taskFuncMap {
		fmt.Println(n)
		f(make(chan bool), Task{})
	}
}

func boo(isDone chan bool, task Task) {
	fmt.Println(12)
}

func setDb() {
	kmgSql.MustLoadTestConfig()
	kmgSql.MustSetTableDataYaml(`
KmgScheduleTask:
  - IsDone: true
    ExecuteTime: "1999-06-15 00:00:00"
    Parameter: "{\"a\":1,\"b\":2}"
    MaxAllowExecuteSecond: 0
    MaxAllowNumberOfRetry: 0
`)
}
