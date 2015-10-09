package serviceCmd

import (
	"fmt"
	"github.com/bronze1man/kmg/errors"
	"github.com/bronze1man/kmg/kmgCrypto"
)

type StartStatus string

const (
	StartStatusSuccess StartStatus = "StartStatusSuccess"
	StartStatusFail                = "StartStatusFail"
)

func ServiceStartSuccess() {
	client := NewClient_ServiceRpc("http://"+rpcAddress, rpcPsk)
	err := client.Send(StartStatusSuccess)
	fmt.Println("src/github.com/bronze1man/kmg/kmg/SubCommand/serviceCmd/rpc.go ServiceStartSuccess", err)
}

func ServiceStartFail() {
	client := NewClient_ServiceRpc("http://"+rpcAddress, rpcPsk)
	err := client.Send(StartStatusFail)
	fmt.Println("src/github.com/bronze1man/kmg/kmg/SubCommand/serviceCmd/rpc.go ServiceStartSuccess", err)
}

//TODO 对于不会发 RPC 的进程，应该可以用一个 flag 来表明不需要等待 RPC 响应
func waitRpcRespond() error {
	lock := &FileMutex{}
	lock.Lock("kmg_service_lock")
	ListenAndServe_ServiceRpc(rpcAddress, &ServiceRpc{}, rpcPsk)
	startStatus := <-statusChannel
	lock.UnLock()
	if startStatus == StartStatusSuccess {
		return nil
	}
	return errors.New("StartFail")
}

var statusChannel = make(chan StartStatus)

var rpcPsk = kmgCrypto.Get32PskFromString("w4n4ts28cq")

var rpcAddress = "127.0.0.1:2777"

type ServiceRpc struct{}

func (sr *ServiceRpc) Send(status StartStatus) {
	statusChannel <- status
}
