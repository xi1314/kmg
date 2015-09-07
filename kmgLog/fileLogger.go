package kmgLog

import (
	"fmt"
	"path/filepath"

	"github.com/bronze1man/kmg/kmgFile"
)

//new file log, will mkdir if dir not exist.
// usage:
// 		kmgLog.DefaultLogger = kmgLog.NewFileLogger("log")
func NewFileLogWriter(logDir string) LogWriter {
	kmgFile.MustMkdirAll(logDir)
	return fileLoger{
		logDir: logDir,
	}
}

type fileLoger struct {
	logDir string
}

func (lw fileLoger) LogWrite(r LogRow) {
	b, err := r.Marshal()
	if err != nil {
		fmt.Println("[fileLoger] logToJson fail", err)
		return
	}
	toWrite := append(b, byte('\n'))
	err = kmgFile.AppendFile(filepath.Join(lw.logDir, r.Cat+".log"), toWrite)
	if err != nil {
		fmt.Println("[fileLoger] logToJson fail", err)
		return
	}
}
