package goCmd

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

/*
GoCrossComplie [gofile]
the output file will put into $project_root/bin/name_GOOS_GOARCH[.exe]
*/
func runGoCrossCompile() {
	command := GoCrossCompile{}
	flag.StringVar(&command.outputPath, "o", "", "output file dir(file name come from source file name),default to $project_root/bin")
	flag.StringVar(&command.version, "v", "", "version string in output file name")
	flag.StringVar(&command.platform, "platform", "", "platform(default use .kmg.yml config)")
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
	targetList := []kmgConfig.CompileTarget{}
	if command.platform == "" {
		targetList = kmgc.CrossCompileTarget
	} else {
		targetList = []kmgConfig.CompileTarget{kmgConfig.CompileTarget(command.platform)}
	}
	for _, target := range targetList {
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
		err := kmgCmd.CmdSlice([]string{"go", "build", "-i", "-o", outputFilePath, targetFile}).
			MustSetEnv("GOOS", target.GetGOOS()).
			MustSetEnv("GOARCH", target.GetGOARCH()).
			MustSetEnv("GOPATH", kmgc.GOPATHToString()).
			Run()
		kmgConsole.ExitOnErr(err)
	}
	return
}

type GoCrossCompile struct {
	outputPath string
	version    string
	platform   string
}
