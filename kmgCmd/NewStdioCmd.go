package kmgCmd

import (
	"io"
	"os"
	"os/exec"
	"strings"
)

func RunOsStdioCmd(name string, args ...string) error {
	return NewOsStdioCmd(name, args...).Run()
}

func MustRunOsStdioCmd(name string, args ...string) {
	err := NewOsStdioCmd(name, args...).Run()
	if err != nil {
		panic(err)
	}
}

// 运行命令,并且把命令输入输出和当前的输入输出接起来(bash的默认调用方式)
// 传入一个字符串是命令,不允许参数中包含空格
// 请不要在参数两边加入单引号.如果需要请使用 RunOsStdioCmd 直接传入.
// 不能在里面使用bash的各种连接符之类的.
// 只能写一条命令
func RunOsStdioCmdString(cmd string) error {
	return NewOsStdioCmdString(cmd).Run()
}

func MustRunOsStdioCmdString(cmd string) {
	err := NewOsStdioCmdString(cmd).Run()
	if err != nil {
		panic(err)
	}
	return
}

// 运行命令,不检查返回状态值,(但是其他类型的错误仍然检查)
func MustRunOsStdioCmdStringNotExitStatusCheck(cmd string) {
	err := NewOsStdioCmdString(cmd).Run()
	if err != nil {
		_, ok := err.(*exec.ExitError)
		if ok {
			return
		}
		panic(err)
	}
	return
}

// 运行命令,返回命令输出的内容
func MustRunOsStdioCmdStringCombinedOutput(cmd string) []byte {
	args := strings.Split(cmd, " ")
	cmdObj := exec.Command(args[0], args[1:]...)
	out, err := cmdObj.CombinedOutput()
	if err != nil {
		panic(err)
	}
	return out
}

// 运行bash命令,和bash一样需要进行转义
func RunOsStdioCmdInBash(cmd string) error {
	return NewOsStdioCmd("bash", "-c", cmd).Run()
}

func MustRunOsStdioCmdInBash(cmd string) {
	err := NewOsStdioCmd("bash", "-c", cmd).Run()
	if err != nil {
		panic(err)
	}
	return
}

func NewOsStdioCmdString(cmd string) *exec.Cmd {
	args := strings.Split(cmd, " ")
	return NewOsStdioCmd(args[0], args[1:]...)
}

func NewOsStdioCmd(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

type Stdio interface {
	GetStdin() io.ReadCloser
	GetStdout() io.WriteCloser
	GetStderr() io.WriteCloser
}

func NewStdioCmd(stdio Stdio, name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Stdin = stdio.GetStdin()
	cmd.Stdout = stdio.GetStdout()
	cmd.Stderr = stdio.GetStderr()
	return cmd
}

var OsStdio = osStdio{}

type osStdio struct{}

func (io osStdio) GetStdin() io.ReadCloser {
	return os.Stdin
}

func (io osStdio) GetStdout() io.WriteCloser {
	return os.Stdout
}
func (io osStdio) GetStderr() io.WriteCloser {
	return os.Stderr
}
