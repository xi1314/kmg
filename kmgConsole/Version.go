package kmgConsole

import (
	"github.com/bronze1man/kmg/kmgFile"
)

func MustWriteGolangVersion(versionFilePath string, version string) {
	kmgFile.MustWriteFile(versionFilePath, []byte(`package main

import "github.com/bronze1man/kmg/kmgConsole"

func init() {
	kmgConsole.VERSION = "`+version+`"
}
`))
}
