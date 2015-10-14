package kmgTask

import "sync"

//implement TaskManager
//简单版taskmanager
//不做任何限制
type SimpleTaskManager struct {
	wg sync.WaitGroup //等待任务完成
}

// 添加一个任务
func (t *SimpleTaskManager) AddTask(task Task) {
	t.wg.Add(1)
	go func() {
		//defer t.wg.Done()
		task.Run()
		// panic不会完成.
		t.wg.Done()
	}()
}

func (t *SimpleTaskManager) AddFunc(task func()) {
	t.wg.Add(1)
	go func() {
		//defer t.wg.Done()
		task()
		// panic不会完成.
		t.wg.Done()
	}()
}

//等待所有任务完成
func (t *SimpleTaskManager) Wait() {
	t.wg.Wait()
}

//关闭管理器
//需要等待所有任务完成后,返回
func (t *SimpleTaskManager) Close() {
	t.Wait()
}

func NewSimpleTaskManager() *SimpleTaskManager{
	return &SimpleTaskManager{}
}

func RunTaskRepeat(f func(), num int) {
	var wg sync.WaitGroup
	for i := 0; i < num; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f()
		}()
	}
	wg.Wait()
}

func RunTaskRepeatWithLimitThread(f func(), taskNum int, threadNum int) {
	var wg sync.WaitGroup
	taskChan := make(chan func())
	for i := 0; i < threadNum; i++ {
		go func() {
			for {
				task, ok := <-taskChan
				if ok == false {
					return
				}
				task()
				wg.Done()
			}
		}()
	}
	wg.Add(taskNum)
	for i := 0; i < taskNum; i++ {
		taskChan <- f
	}
	close(taskChan)
	wg.Wait()
}

func RunTask(threadNum int, funcList ...func()) {
	if threadNum > len(funcList) {
		threadNum = len(funcList)
	}
	var wg sync.WaitGroup
	taskChan := make(chan func())
	for i := 0; i < threadNum; i++ {
		go func() {
			for {
				task, ok := <-taskChan
				if ok == false {
					return
				}
				task()
				wg.Done()
			}
		}()
	}
	wg.Add(len(funcList))
	for i := range funcList {
		taskChan <- funcList[i]
	}
	close(taskChan)
	wg.Wait()
}
