package kmgFileToXcode
import (
	"github.com/bronze1man/kmg/kmgCmd"
	"os"
	"strings"
)
func AddFileToXcode(FilePath string,ProjectPath string)[]byte{
	cmd := kmgCmd.CmdSlice(append([]string{"ruby","AddFileToXcode.rb",FilePath,ProjectPath},os.Args[1:]...))
    dir:=changeDir()
	cmd.SetDir(dir + "/src/github.com/bronze1man/kmg/kmgFileToXcode")
	out := cmd.MustCombinedOutput()
	return out
}

func changeDir() string {
	dir, _ := os.Getwd()
	tmp := strings.SplitN(dir, "src/", -1)
	if len(tmp) > 1 {
		return tmp[0]
	} else {
		return dir
	}
}
