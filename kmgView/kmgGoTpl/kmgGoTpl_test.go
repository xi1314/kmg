package kmgGoTpl

import (
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgTest"
	"path/filepath"
	"testing"
)

func TestGoTpl(ot *testing.T) {
	MustBuildTplInDir("testFile")

	files := kmgFile.MustGetAllFiles("testFile")
	for _, file := range files {
		if filepath.Ext(file) != ".gotpl" {
			continue
		}
		generated := kmgFile.MustReadFile(filepath.Join(filepath.Dir(file), kmgFile.GetFileBaseWithoutExt(file)+".go"))
		correct := kmgFile.MustReadFile(filepath.Join(filepath.Dir(file), kmgFile.GetFileBaseWithoutExt(file)+".go.good"))
		kmgTest.Equal(generated, correct, file)
	}
}

func setCurrentAsCorrect() {
	files := kmgFile.MustGetAllFiles("testFile")
	for _, file := range files {
		if filepath.Ext(file) != ".gotpl" {
			continue
		}
		kmgFile.MustCopyFile(filepath.Join(filepath.Dir(file), kmgFile.GetFileBaseWithoutExt(file)+".go"), filepath.Join(filepath.Dir(file), kmgFile.GetFileBaseWithoutExt(file)+".go.good"))
	}
}
