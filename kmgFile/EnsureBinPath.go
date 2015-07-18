package kmgFile

import (
	"github.com/bronze1man/kmg/kmgTime"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"github.com/bronze1man/kmg/kmgCmd"
)

func MustEnsureBinPath(finalPath string) {
	basePath := filepath.Base(finalPath)
	path, err := exec.LookPath(basePath)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		panic(err)
	}
	if path != finalPath {
		backPathDir := "/var/backup/bin/" + basePath + time.Now().Format(kmgTime.FormatFileName)
		MustMkdirAll(backPathDir)
		kmgCmd.MustRun("mv " + path + " " + backPathDir)
	}
}
