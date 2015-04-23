package kmgCmd

import (
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgTime"
	"os"
	"os/exec"
	"path/filepath"
	"time"
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
		kmgFile.MustMkdirAll(backPathDir)
		MustRun("mv " + path + " " + backPathDir)
	}
}
