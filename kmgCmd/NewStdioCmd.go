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

func RunOsStdioCmdString(cmd string) error {
	return NewOsStdioCmdString(cmd).Run()
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
