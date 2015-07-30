package kmgSshCertCopy

import (
	"fmt"

	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgRand"
)

type RemoteServer struct {
	Ip       string
	Password string
}

func CopyLocalToRemote(remoteServerList []RemoteServer) {
	for _, remote := range remoteServerList {
		runCmdWithPassword(fmt.Sprintf("ssh-copy-id root@"+remote.Ip), remote.Password)
	}
}

func CopyCertToRemote(cert string, remoteList []RemoteServer) {
}

func runCmdWithPassword(cmd, password string) {
	cmdTpl := `#!/usr/bin/expect -f
spawn %s
expect "assword:"
send "%s\r"
interact`
	cmd = fmt.Sprintf(cmdTpl, cmd, password)
	tmpName := kmgRand.MustCryptoRandToReadableAlphaNum(5)
	tmpPath := "/tmp/" + tmpName
	kmgFile.MustAppendFile(tmpPath, []byte(cmd))
	defer kmgFile.MustDelete(tmpPath)
	kmgCmd.CmdSlice([]string{tmpPath}).MustStdioRun()
}
