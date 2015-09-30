package kmgProcess

import (
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
)

func TestExtractProcessList(t *testing.T) {
	pList := Extract(` 5393 /usr/bin/InlProxy vpnAccount -redisAddr=127.0.0.1:30001
 5424 /usr/bin/xasdf -block
 5455 /usr/bin/asdfkja;sdflk -nl asdfadsf;ak
 5590 /usr/bin/abc
 6625 InlProxy
 7254 InlProxy abc
27939 /usr/bin/InlProxy`)
	kmgTest.Equal(len(pList), 7)
	kmgTest.Equal(pList[1].Id, 5424)
	kmgTest.Equal(pList[5].Command, "InlProxy abc")
}

func TestDiff(t *testing.T) {
	notExpect, notRunning := Diff([]*Process{
		{Command: "a"},
		{Command: "b"},
		{Command: "a"},
		{Command: "c"},
	}, []*Process{
		{Command: "a"},
		{Command: "a"},
		{Command: "c"},
		{Command: "e"},
		{Command: "e"},
	})
	kmgTest.Equal(len(notExpect), 2)
	kmgTest.Equal(notExpect[0].Command, "e")
	kmgTest.Equal(notExpect[1].Command, "e")
	kmgTest.Equal(len(notRunning), 1)
	kmgTest.Equal(notRunning[0].Command, "b")
}
