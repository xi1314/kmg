package kmgScheduleTask

import (
	"container/heap"
	"fmt"
	"github.com/bronze1man/kmg/encoding/kmgJson"
	"github.com/bronze1man/kmg/kmgErr"
	"github.com/bronze1man/kmg/kmgLog"
	"github.com/bronze1man/kmg/kmgRand"
	"github.com/bronze1man/kmg/kmgSql"
	"github.com/bronze1man/kmg/kmgTime"
	"sync"
	"time"
)

type Task struct {
	Id          string
	FuncName    string
	InMap       map[string]string
	ExecuteTime time.Time
}

type TaskSlice []Task

func (ts *TaskSlice) Len() int           { return len(*ts) }
func (ts *TaskSlice) Less(i, j int) bool { return (*ts)[i].ExecuteTime.Before((*ts)[j].ExecuteTime) }
func (ts *TaskSlice) Swap(i, j int)      { (*ts)[i], (*ts)[j] = (*ts)[j], (*ts)[i] }
func (ts *TaskSlice) Push(x interface{}) { *ts = append(*ts, x.(Task)) }
func (ts *TaskSlice) Pop() interface{} {
	ele := (*ts)[len(*ts)-1]
	*ts = (*ts)[:len(*ts)-1]
	return ele
}

func RegisterTask(task Task) {
	defaultTaskManager.RegisterTask(task)
}

// 任务内部抛出的异常,处理方式和kmgControllerRunner差不多.会全部截取掉,然后写个log
func RegisterTaskFunc(FuncName string, Func func(task Task)) {
	defaultTaskManager.RegisterTaskFunc(FuncName, Func)
}

var initOnce sync.Once

func Init() {
	initOnce.Do(func() {
		defaultTaskManager.Init()
	})
}

func RegisterTable() {
	kmgSql.MustRegisterTable(kmgSql.Table{
		Name: "kmgScheduleTask",
		FieldList: map[string]kmgSql.DbType{
			"Id":          kmgSql.DbTypeString,
			"ExecuteTime": kmgSql.DbTypeDatetime,
			"FuncName":    kmgSql.DbTypeString,
			"InMap":       kmgSql.DbTypeLongBlob,
		},
		PrimaryKey: "Id",
	})
}

var defaultTaskManager = &taskManager{}

//搞这个主要是纯函数接口有全局状态不好测试.
type taskManager struct {
	taskList     TaskSlice
	funcMap      map[string]func(task Task)
	lock         sync.Mutex
	pullDuration time.Duration
	closeChan    chan struct{}
}

// 请先Init 后注册任务
func (tm *taskManager) RegisterTask(task Task) {
	tm.lock.Lock()
	defer tm.lock.Unlock()
	if tm.taskList == nil {
		tm.taskList = make(TaskSlice, 0)
	}
	task.Id = kmgRand.MustCryptoRandToAlphaNum(16)
	tm.taskList.Push(task)
	kmgSql.Insert("kmgScheduleTask", map[string]string{
		"Id":          task.Id,
		"ExecuteTime": kmgTime.DefaultFormat(task.ExecuteTime),
		"InMap":       string(kmgJson.MustMarshal(task.InMap)),
		"FuncName":    task.FuncName,
	})
	kmgLog.Log("kmgScheduleTask", "register", task)
}

// 请先Init 后注册任务
func (tm *taskManager) RegisterTaskFunc(FuncName string, Func func(task Task)) {
	tm.lock.Lock()
	defer tm.lock.Unlock()
	if tm.funcMap == nil {
		tm.funcMap = map[string]func(task Task){}
	}
	tm.funcMap[FuncName] = Func
}

// 请先注册任务 后Init
func (tm *taskManager) Init() {
	tm.lock.Lock()
	defer tm.lock.Unlock()
	if tm.pullDuration == 0 {
		tm.pullDuration = time.Second
	}
	tm.closeChan = make(chan struct{})
	rowList := kmgSql.MustGetAllInTable("kmgScheduleTask")
	if tm.taskList == nil {
		tm.taskList = make(TaskSlice, 0, len(rowList))
	}
	for _, row := range rowList {
		t := Task{
			Id:          row["Id"],
			FuncName:    row["FuncName"],
			ExecuteTime: kmgTime.MustFromMysqlFormatDefaultTZ(row["ExecuteTime"]),
		}
		kmgJson.MustUnmarshal([]byte(row["InMap"]), &t.InMap)
		tm.taskList = append(tm.taskList, t)
	}
	heap.Init(&tm.taskList)
	go tm.run()
}

// 该函数不允许并发调用
func (tm *taskManager) run() {
	for {
		select {
		case <-time.After(tm.pullDuration):
			tm.runOne()
		case <-tm.closeChan:
		}
	}
}

func (tm *taskManager) Close() error {
	tm.closeChan <- struct{}{}
	return nil
}

func (tm *taskManager) runOne() {
	tm.lock.Lock()
	defer tm.lock.Unlock()
	now := time.Now()
	for {
		//检查一下最早的时间是不是已经到了
		if len(tm.taskList) == 0 {
			return
		}
		if tm.taskList[0].ExecuteTime.After(now) {
			return
		}
		task := heap.Pop(&tm.taskList).(Task)
		f := tm.funcMap[task.FuncName]
		if f == nil {
			// 这个是编程错误,程序不要运行了吧.赶紧修bug.
			panic(fmt.Errorf("[kmgScheduleTask.runOne] funcName[%s] not exist", task.FuncName))
		}
		kmgSql.MustDeleteById("kmgScheduleTask", "Id", task.Id)
		//不需要做是否正在运行的东西.正在运行的都从列表中删除了.
		go func(task Task) {
			err := kmgErr.PanicToError(func() {
				f(task)
			})
			if err != nil {
				kmgErr.LogErrorWithStack(err)
			}
		}(task)
		kmgLog.Log("kmgScheduleTask", "run", task)
	}
}
