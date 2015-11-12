package kmgFileToXcode

import (
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgConfig/defaultEnv"
	"github.com/bronze1man/kmg/kmgStrings"
)

func AddFileToXcode(FilePath string, ProjectPath string) []byte {
	dir := defaultEnv.Env().ProjectPath + "/"
	if !kmgStrings.IsStartWith(FilePath, "~") && !kmgStrings.IsStartWith(FilePath, "/") {
		FilePath = dir + FilePath
	}
	if !kmgStrings.IsStartWith(ProjectPath, "~") && !kmgStrings.IsStartWith(ProjectPath, "/") {
		ProjectPath = dir + ProjectPath
	}
	cmd := kmgCmd.CmdBash("export LANG=UTF-8;ruby AddFileToXcode.rb " + FilePath + " " + ProjectPath)
	cmd.SetDir(dir + "src/github.com/bronze1man/kmg/kmgFileToXcode")
	out := cmd.MustRunAndReturnOutput()
	return out
}
func AddFilesToXcode(FilePaths []string, ProjectPath string) {
	for _, s := range FilePaths {
		AddFileToXcode(s, ProjectPath)
	}
}
