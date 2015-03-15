package goCommand

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgConsole"
	"github.com/bronze1man/kmg/kmgFile"
	"os"
)

func init() {
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "GoCrossCompile",
		Desc:   "cross compile target in current project",
		Runner: runGoCrossCompile,
	})
}

/*
GoCrossComplie [gofile]
the output file will put into $project_root/bin/name_GOOS_GOARCH[.exe]
*/
func runGoCrossCompile() {
	command := GoCrossCompile{}
	flag.StringVar(&command.outputPath, "o", "", "output file dir(file name come from source file name),default to $project_root/bin")
	flag.StringVar(&command.version, "v", "", "version string in output file name")
	flag.Parse()

	if len(os.Args) <= 1 {
		kmgConsole.ExitOnErr(fmt.Errorf("need gofile parameter"))
		return
	}
	targetFile := flag.Arg(0)
	kmgc, err := kmgConfig.LoadEnvFromWd()
	kmgConsole.ExitOnErr(err)
	targetName := kmgFile.GetFileBaseWithoutExt(targetFile)
	if command.outputPath == "" {
		command.outputPath = filepath.Join(kmgc.ProjectPath, "bin")
	}
	for _, target := range kmgc.CrossCompileTarget {
		fileName := ""
		if command.version == "" {
			fileName = targetName + "_" + target.GetGOOS() + "_" + target.GetGOARCH()
		} else {
			fileName = targetName + "_" + command.version + "_" + target.GetGOOS() + "_" + target.GetGOARCH()
		}

		if target.GetGOOS() == "windows" {
			fileName = fileName + ".exe"
		}
		outputFilePath := filepath.Join(command.outputPath, fileName)
		cmd := kmgCmd.NewOsStdioCmd("go", "build", "-i", "-o", outputFilePath, targetFile)
		kmgCmd.SetCmdEnv(cmd, "GOOS", target.GetGOOS())
		kmgCmd.SetCmdEnv(cmd, "GOARCH", target.GetGOARCH())
		kmgCmd.SetCmdEnv(cmd, "GOPATH", kmgc.GOPATHToString())
		err = cmd.Run()
		kmgConsole.ExitOnErr(err)
	}
	return
}

type GoCrossCompile struct {
	outputPath string
	version    string
}
