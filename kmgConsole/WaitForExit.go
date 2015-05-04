package kmgConsole

import (
	"os"
	"os/signal"
	"syscall"
)

//等用户按CTRL+c退出,或者等到被kill掉
func WaitForExit() {
	//不要在这个地方检查WaitForExit和AddExitAction一起使用,因为程序自身会进行调用
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-ch
}

var exitActionList = []func(){}

//如果你使用了这个命令,主线程不会退出,而是会等到用户或者系统发送退出命令才会退出.
//使用命令调用时,注册退出动作,
//如果你使用了这个东西,请不要再使用WaitForExit了
//TODO 这个地方过于不直观
func AddExitAction(f func()) {
	exitActionList = append(exitActionList, f)
}

//调用这个函数来保证使用AddExitAction方法来注册进程退出请求.
// 使用这个来保证kmgConsole.Main一定会等待进程结束
func UseExitActionRegister() {
	AddExitAction(func() {})
}
