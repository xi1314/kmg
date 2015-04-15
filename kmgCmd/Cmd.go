package kmgCmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

//please use Cmd* function to new a Cmd,do not create one yourself.
type Cmd struct {
	cmd *exec.Cmd
}

//you need at least one args: the path of the command, or it will panic
func CmdSlice(args []string) *Cmd {
	if len(args) == 0 {
		panic("[CmdSlice] need the path of the command")
	}
	return &Cmd{
		cmd: exec.Command(args[0], args[1:]...),
	}
}

func CmdString(cmd string) *Cmd {
	if cmd == "" {
		panic("[CmdString] need the path of the command")
	}
	args := strings.Split(cmd, " ")
	return &Cmd{
		cmd: exec.Command(args[0], args[1:]...),
	}
}

func CmdBash(cmd string) *Cmd {
	return CmdSlice([]string{"bash", "-c", cmd})
}

func (c *Cmd) MustSetEnv(key string, value string) *Cmd {
	err := SetCmdEnv(c.cmd, key, value)
	if err != nil {
		panic(err)
	}
	return c
}

func (c *Cmd) SetDir(path string) *Cmd {
	c.cmd.Dir = path
	return c
}

func (c *Cmd) PrintCmdLine() {
	fmt.Println(">", c.cmd.Path, strings.Join(c.cmd.Args, " "))
}

//回显命令,并且运行,并且和标准输入输出接起来
func (c *Cmd) Run() error {
	c.PrintCmdLine()
	c.cmd.Stdin = os.Stdin
	c.cmd.Stdout = os.Stdout
	c.cmd.Stderr = os.Stderr
	return c.cmd.Run()
}

//get the os/exec.Cmd
func (c *Cmd) GetExecCmd() *exec.Cmd {
	return c.cmd
}

//回显命令,并且运行,返回运行的输出结果.并且把输出结果放在stdout中
func (c *Cmd) RunAndReturnOutput() (b []byte, err error) {
	c.PrintCmdLine()
	b, err = c.cmd.CombinedOutput()
	os.Stdout.Write(b)
	return b, err
}

//不回显命令,运行,并且返回运行的输出结果
func (c *Cmd) StdioRun() error {
	c.cmd.Stdin = os.Stdin
	c.cmd.Stdout = os.Stdout
	c.cmd.Stderr = os.Stderr
	return c.cmd.Run()
}

//回显命令,并且运行,并且忽略输出结果
func (c *Cmd) RunAndNotExitStatusCheck() error {
	err := c.Run()
	if err != nil {
		_, ok := err.(*exec.ExitError)
		if ok {
			return nil
		}
		return err
	}
	return nil
}

func (c *Cmd) MustStdioRun() {
	err := c.StdioRun()
	if err != nil {
		panic(err)
	}
}

func (c *Cmd) MustRunAndReturnOutput() (b []byte) {
	b, err := c.RunAndReturnOutput()
	if err != nil {
		panic(err)
	}
	return b
}

func (c *Cmd) MustRunAndNotExitStatusCheck() {
	err := c.RunAndNotExitStatusCheck()
	if err != nil {
		panic(err)
	}
}

func (c *Cmd) MustRun() {
	err := c.Run()
	if err != nil {
		panic(err)
	}
}
