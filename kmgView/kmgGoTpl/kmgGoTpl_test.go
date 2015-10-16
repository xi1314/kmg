package kmgGoTpl

import (
	"path/filepath"
	"testing"

	"github.com/bronze1man/kmg/kmgFile"
	"fmt"
	"bytes"
)

func TestGoTpl(ot *testing.T) {
	MustBuildTplInDir("testFile")
	files := kmgFile.MustGetAllFiles("testFile")
	AllCurrect:=true
	ErrorMsg:=""
	for _, file := range files {
		if filepath.Ext(file) != ".gotplhtml" {
			continue
		}
		generated := kmgFile.MustReadFile(filepath.Join(filepath.Dir(file), kmgFile.GetFileBaseWithoutExt(file)+".go"))
		correct,err := kmgFile.ReadFile(filepath.Join(filepath.Dir(file), kmgFile.GetFileBaseWithoutExt(file)+".go.good"))
		if err!=nil{
			ErrorMsg+=fmt.Sprintf("%s read good file fail err [%s]\n",file,err)
			AllCurrect = false
			continue
		}
		if !bytes.Equal(generated,correct){
			ErrorMsg+=fmt.Sprintf("%s generated not equal correct\n",file)
			AllCurrect = false
			continue
		}
	}
	if !AllCurrect{
		panic(ErrorMsg)
	}
}

func setCurrentAsCorrect() {
	files := kmgFile.MustGetAllFiles("testFile")
	for _, file := range files {
		if filepath.Ext(file) != ".gotplhtml" {
			continue
		}
		kmgFile.MustCopyFile(filepath.Join(filepath.Dir(file), kmgFile.GetFileBaseWithoutExt(file)+".go"), filepath.Join(filepath.Dir(file), kmgFile.GetFileBaseWithoutExt(file)+".go.good"))
	}
}
