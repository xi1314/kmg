package serviceCmd

import (
	"encoding/json"
	"fmt"
	"github.com/bronze1man/kmg/kmgConsole"
	"os"
)

//kmg service.process start {"Name":"xxx",}
func processCmd() {
	if len(os.Args) < 3 {
		kmgConsole.ExitOnErr(fmt.Errorf(`example:
kmg service.process start {"Name":"xxx",}`))
	}
	cmd := os.Args[1]
	dataJson := os.Args[2]
	s := &Service{}
	err := json.Unmarshal([]byte(dataJson), s)
	kmgConsole.ExitOnErr(err)
	switch cmd {
	case "start":

	case "stop":
	case "restart":
	case "status":
	default:
		kmgConsole.ExitOnErr(fmt.Errorf(`not support command.`))
	}
}

func processIsRuning() {

}
