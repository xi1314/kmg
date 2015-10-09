package main

import (
	"fmt"
	"github.com/bronze1man/kmg/kmg/SubCommand/serviceCmd"
	"time"
)

func main() {
	t := time.Now()
	time.Sleep(time.Second * 3)
	serviceCmd.ServiceStartSuccess()
	fmt.Println("ServiceStartSuccess", time.Now().Sub(t).String())
}
