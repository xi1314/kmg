package kmgSsh

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgRand"
	"github.com/bronze1man/kmg/kmgSys"
	"strconv"
	"strings"
	"os"
)

type RemoteServer struct {
	Address  string
	Password string
	UserName string //默认 root
	SshPort  int    //默认 22
}

func (r *RemoteServer) String() string {
	if r.SshPort == 0 {
		r.SshPort = 22
	}
	if r.UserName == "" {
		r.UserName = "root"
	}
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
	if remote.Address == "" {
		return
	}
	isReachable, havePermission := AvailableCheck(remote)
	if !isReachable {
		panic("[kmgSsh SshCertCopyLocalToRemote]" + remote.String() + " unreachable!")
	}
	if havePermission {
		return
	}
	if remote.Password == "" {
		kmgCmd.MustRunInBash("ssh-copy-id " + remote.String())
		return
	}
	//p := filepath.Join(kmgSys.GetCurrentUserHomeDir(), ".ssh", "id_rsa.pub")
	RunCmdWithPassword(
		//strings.Join([]string{"ssh-copy-id", "-i", p, remote.String()}, " "),
		"ssh-copy-id "+remote.String(),
		remote.Password,
	)
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

func RunCmdWithPassword(cmd, password string) {
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
	os.Setenv("HOME", kmgSys.GetCurrentUserHomeDir())
	kmgCmd.CmdSlice([]string{tmpPath}).MustStdioRun()
}
