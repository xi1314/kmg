package kmgSys

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgRand"
	"strconv"
	"strings"
)

type RemoteServer struct {
	Address  string
	Password string
	UserName string //默认 root
	SshPort  int    //默认 22
}

func completeRemote(remote *RemoteServer) {
	if remote.SshPort == 0 {
		remote.SshPort = 22
	}
	if remote.UserName == "" {
		remote.UserName = "root"
	}
}

func (r *RemoteServer) String() string {
	completeRemote(r)
	cmd := []string{}
	cmd = append(cmd, "-p", strconv.Itoa(r.SshPort), r.UserName+"@"+r.Address)
	return strings.Join(cmd, " ")
}

func SshCertCopyLocalToRemoteRoot(remoteAddress string) {
	SshCertCopyLocalToRemote(&RemoteServer{
		Address: remoteAddress,
	})
}

func SshCertCopyLocalToRemote(remote *RemoteServer) {
	if IsLocalSshCertCopyToRemote(remote) {
		return
	}
	if remote.Password == "" {
		kmgCmd.MustRunInBash("ssh-copy-id " + remote.String())
		return
	}
	runCmdWithPassword(
		"ssh-copy-id "+remote.String(),
		remote.Password,
	)
}

func IsLocalSshCertCopyToRemoteRoot(remoteAddress string) bool {
	return IsLocalSshCertCopyToRemote(&RemoteServer{
		Address: remoteAddress,
	})
}

func IsLocalSshCertCopyToRemote(remote *RemoteServer) bool {
	cmd := "ssh -o StrictHostKeyChecking=no -o PasswordAuthentication=no " + remote.String() + " echo OK"
	err := kmgCmd.Run(cmd)
	if err != nil {
		return false
	}
	b := kmgCmd.MustRunAndReturnOutput(cmd)
	if strings.Contains(string(b), "OK") {
		return true
	}
	return false
}

func SshCertCopyCertToRemote(cert string, remoteList []RemoteServer) {
	kmgFile.MustWriteFile("/tmp/cert", []byte(cert))
	defer kmgFile.MustDelete("/tmp/cert")
	for _, remote := range remoteList {
		cmd := []string{"ssh", "-t", "-p", strconv.Itoa(remote.SshPort), remote.UserName + "@" + remote.Address}
		kmgCmd.CmdSlice([]string{"scp", "-P", strconv.Itoa(remote.SshPort), "/tmp/cert", remote.UserName + "@" + remote.Address + ":~/"}).MustRun()
		authorizedKeysByte := kmgCmd.CmdSlice(append(cmd, "cat ~/.ssh/authorized_keys")).MustRunAndReturnOutput()
		if strings.Contains(string(authorizedKeysByte), cert) {
			fmt.Println("Cert has contains in authorized_keys")
			continue
		}
		kmgCmd.CmdSlice(append(cmd, "mkdir .ssh;cat cert >> .ssh/authorized_keys;rm cert")).MustRun()
	}
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
