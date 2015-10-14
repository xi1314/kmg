package main

import (
	"fmt"
	"github.com/bronze1man/kmg/kmg/SubCommand/serviceCmd"
	"time"
)

func main() {
	t := time.Now()
	serviceCmd.ServiceStartSuccess()
	//serviceCmd.ServiceStartFail()
	fmt.Println("ServiceStartSuccess", time.Now().Sub(t).String())
}
