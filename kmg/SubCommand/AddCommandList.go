package SubCommand

import "github.com/bronze1man/kmg/kmgConsole"

func AddCommandList() {
	kmgConsole.AddAction(kmgConsole.Command{
		Name: "Make",
		Desc: `run a project defined command
保证在项目根目录下运行
使用普通空格分割方法定义命令
将命令输出结果log到文件中`,
		Runner: makeCmd,
	})
	kmgConsole.AddCommandWithName("NewPassword", NewPassword)
	kmgConsole.AddCommandWithName("HttpsCertCsr",httpsCertCsrCLI)
}
