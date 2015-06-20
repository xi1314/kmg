package kmgScheduleTask

import (
	"fmt"
	"github.com/bronze1man/kmg/encoding/kmgJson"
	"github.com/bronze1man/kmg/kmgSql"
	"github.com/bronze1man/kmg/kmgTime"
	"reflect"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type Status string

const (
	MySQLTableName               string        = "KmgScheduleTask"
	DefaultMaxAllowExecuteTime   time.Duration = 10 * time.Minute
	DefaultMaxAllowNumberOfRetry int           = 3
)

type Task struct {
	Id                    int
	ExecuteTime           time.Time                         // 计划执行的时间点
	Func                  func(isDone chan bool, task Task) // Task 要执行的函数
	Parameter             map[string]string                 // Task 的参数
	IsDone                bool                              // Task 是否已经已经完成
	MaxAllowExecuteSecond time.Duration
	MaxAllowNumberOfRetry int // Task 允许的最大重试次数
	NumberOfTry           int // Task 已重试的次数 TODO 暂时不在数据库保存该值
}

func (task *Task) doneAndPersist() {
	row := map[string]string{}
	if task.Id != 0 {
		row["Id"] = strconv.Itoa(task.Id)
	}
	row["ExecuteTime"] = kmgTime.DefaultFormat(task.ExecuteTime)
	row["Func"] = getFuncFullName(task.Func)
	row["Parameter"] = kmgJson.MustMarshalToString(task.Parameter)
	row["IsDone"] = "1"
	row["MaxAllowExecuteSecond"] = task.MaxAllowExecuteSecond.String()
	row["MaxAllowNumberOfRetry"] = strconv.Itoa(task.MaxAllowNumberOfRetry)
	kmgSql.MustReplaceById(MySQLTableName, "Id", row)
}

var taskSlice []*Task

var taskFuncMap map[string]func(isDone chan bool, task Task) = map[string]func(isDone chan bool, task Task){}

func getFuncFromString(fullName string) func(isDone chan bool, task Task) {
	function, ok := taskFuncMap[fullName]
	if ok {
		return function
	}
	return nil
}

func getFuncFullName(function interface{}) string {
	v := reflect.ValueOf(function)
	f := runtime.FuncForPC(v.Pointer())
	return f.Name()
}

func recoverFromDb() {
	rowList := kmgSql.MustQuery("SELECT * FROM `"+MySQLTableName+"` WHERE IsDone<>?", "1")
	for _, row := range rowList {
		task := &Task{
			ExecuteTime: kmgTime.MustFromMysqlFormatDefaultTZ(row["ExecuteTime"]),
			Func:        getFuncFromString(row["Func"]),
			IsDone:      false,
		}
		task.Parameter = make(map[string]string)
		kmgJson.MustUnmarshal([]byte(row["Parameter"]), task.Parameter)
		id, err := strconv.Atoi(row["Id"])
		if err != nil {
			id = 0
			log(task, err.Error())
		}
		task.Id = id
		maxAllowExecuteSecond, err := strconv.Atoi(row["MaxAllowExecuteSecond"])
		if err != nil {
			maxAllowExecuteSecond = 0
			log(task, err.Error())
		}
		task.MaxAllowExecuteSecond = time.Duration(maxAllowExecuteSecond)
		maxAllowNumberOfRetry, err := strconv.Atoi(row["MaxAllowNumberOfRetry"])
		if err != nil {
			maxAllowNumberOfRetry = 0
			log(task, err.Error())
		}
		task.MaxAllowNumberOfRetry = maxAllowNumberOfRetry
		taskSlice = append(taskSlice, task)
	}
}

func run() {
	for _, task := range taskSlice {
		if task.IsDone {
			//TODO 将其从 taskSlice 中删除，提高一点点性能
			continue
		}
		now := kmgTime.NowTime.Now()
		d := now.Sub(task.ExecuteTime)
		if int(d) < 0 {
			continue
		}
		if task.Func == nil {
			log(task, "task has no Func to execute")
			task.IsDone = true // TODO 以后都不再理它，避免打一堆没用的Log
			continue
		}
		exec(task)
	}
	time.Sleep(time.Second)
	run()
}

func exec(task *Task) {
	maxAllowNumberOfRetry := DefaultMaxAllowNumberOfRetry
	if task.MaxAllowNumberOfRetry > 0 {
		maxAllowNumberOfRetry = task.MaxAllowNumberOfRetry
	}
	if task.NumberOfTry >= maxAllowNumberOfRetry { //TODO 似乎这里不需要判断
		return
	}
	task.NumberOfTry++
	timeout := DefaultMaxAllowExecuteTime
	if task.MaxAllowExecuteSecond > 0 {
		timeout = task.MaxAllowExecuteSecond * time.Second
	}
	isDoneCh := make(chan bool)
	go task.Func(isDoneCh, *task)
	select {
	case isDone := <-isDoneCh:
		fmt.Println(isDone)
		if isDone {
			task.doneAndPersist()
			return
		}
		if task.NumberOfTry >= maxAllowNumberOfRetry {
			task.doneAndPersist()
			log(task, "task failed "+strconv.Itoa(task.NumberOfTry)+" times")
			return
		}
		exec(task)
	case <-time.After(timeout):
		task.doneAndPersist()
		log(task, "time out")
	}
}

func log(task *Task, info string) {
	fmt.Println(info)
}

var initOnce sync.Once

func Start() {
	initOnce.Do(func() {
		kmgSql.MustRegisterTable(kmgSql.Table{
			Name: MySQLTableName,
			FieldList: map[string]kmgSql.DbType{
				"Id":                    kmgSql.DbTypeIntAutoIncrement,
				"ExecuteTime":           kmgSql.DbTypeDatetime,
				"Func":                  kmgSql.DbTypeString,
				"Parameter":             kmgSql.DbTypeLongBlob,
				"IsDone":                kmgSql.DbTypeBool,
				"MaxAllowExecuteSecond": kmgSql.DbTypeString,
				"MaxAllowNumberOfRetry": kmgSql.DbTypeInt,
			},
			PrimaryKey: "Id",
		})
		kmgSql.SyncDbCommand()
		recoverFromDb()
		run()
	})
}

//可以在任何位置注册
func RegisterTask(task Task) {
	taskSlice = append(taskSlice, &task)
}

//必须在 Start 前注册好一些数据库中可能保存的 TaskFunc
func RegisterTaskFunc(taskFunc interface{}) {
	v := reflect.ValueOf(taskFunc)
	f := runtime.FuncForPC(v.Pointer())
	taskFuncMap[f.Name()] = v.Interface().(func(isDone chan bool, task Task))
}
