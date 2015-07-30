package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/bronze1man/kmg/kmgConsole"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgTime"
)

var logFilePath string

func main() {
	os.Stderr.WriteString("Stderr: start\n")
	kmgFile.MustMkdirAll("/var/ServiceTester")
	pid := os.Getpid()
	logFilePath = filepath.Join("/var/ServiceTester", time.Now().Format(kmgTime.FormatFileName)+"_"+strconv.Itoa(pid)+".log")
	log("start")
	go func() {
		for {
			time.Sleep(time.Second)
			log("running")
		}
	}()
	kmgConsole.WaitForExit()
	log("stop")
}
func log(msg string) {
	wMsg := fmt.Sprintf("%s %s\n", time.Now(), msg)
	os.Stdout.WriteString(wMsg)
	kmgFile.MustAppendFile(logFilePath, []byte(wMsg))
}
