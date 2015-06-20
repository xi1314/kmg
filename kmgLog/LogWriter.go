package kmgLog

import (
	"encoding/json"
	"fmt"
	"github.com/bronze1man/kmg/encoding/kmgJson"
	"os"
	"time"
)

type LogWriter interface {
	LogWrite(r LogRow)
}

type loggerWithWriter struct {
	w LogWriter
}

func (logger loggerWithWriter) Log(cat string, data ...interface{}) {
	logRow := LogRow{
		Cat:  cat,
		Time: time.Now().Format(time.RFC3339),
		Data: make([]json.RawMessage, len(data)),
	}

	for i := range data {
		var err error
		logRow.Data[i], err = json.Marshal(data[i])
		if err != nil {
			fmt.Fprintln(os.Stderr, "[loggerWithWriter] json.Marshal fail", err)
			return
		}
	}
	logger.w.LogWrite(logRow)
}

type LogRow struct {
	Cat  string
	Time string
	Data []json.RawMessage
}

func (r LogRow) Marshal() (b []byte, err error) {
	return kmgJson.MarshalIndent(r)
}
