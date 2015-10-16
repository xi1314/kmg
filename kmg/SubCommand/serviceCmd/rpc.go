package serviceCmd

import (
	"fmt"
	"github.com/bronze1man/kmg/errors"
	"github.com/bronze1man/kmg/kmgCrypto"
	"github.com/bronze1man/kmg/kmgProcessMutex"
	"time"
)

type StartStatus string

const (
	StartStatusSuccess StartStatus = "StartStatusSuccess"
	StartStatusFail                = "StartStatusFail"
)

var statusChannel = make(chan StartStatus)

var rpcPsk = kmgCrypto.Get32PskFromString("w4n4ts28cq")

var rpcAddress = "127.0.0.1:2777"

type ServiceRpc struct{}

func (sr *ServiceRpc) Send(status StartStatus) {
	statusChannel <- status
}

//TODO 对于不会发 RPC 的进程，应该可以用一个 flag 来表明不需要等待 RPC 响应
func waitRpcRespond() chan error {
	returnChan := make(chan error)
	lock := &kmgProcessMutex.FileMutex{Name: "kmg_service_lock"}
	lock.Lock()
	ListenAndServe_ServiceRpc(rpcAddress, &ServiceRpc{}, rpcPsk)
	go func() {
		startStatus := <-statusChannel
		lock.UnLock()
		if startStatus == StartStatusSuccess {
			returnChan <- nil
		}
		returnChan <- errors.New("StartFail")
	}()
	return returnChan
}

func ServiceStartSuccess() {
	time.Sleep(time.Millisecond * 100)
	client := NewClient_ServiceRpc("http://"+rpcAddress, rpcPsk)
	err := client.Send(StartStatusSuccess)
	if err != nil {
		fmt.Println("src/github.com/bronze1man/kmg/kmg/SubCommand/serviceCmd/rpc.go ServiceStartSuccess", err)
	}
}

func ServiceStartFail() {
	time.Sleep(time.Second)
	client := NewClient_ServiceRpc("http://"+rpcAddress, rpcPsk)
	err := client.Send(StartStatusFail)
	if err != nil {
		fmt.Println("src/github.com/bronze1man/kmg/kmg/SubCommand/serviceCmd/rpc.go ServiceStartFail", err)
	}
}
