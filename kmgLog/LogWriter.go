package kmgLog

import (
	"encoding/json"
	"fmt"
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
		Time: time.Now().Format(time.RFC3339Nano),
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
	return json.Marshal(r)
}

func (r LogRow) UnmarshalData(index int, obj interface{}) (err error) {
	return json.Unmarshal(r.Data[index], obj)
}
func (r LogRow) MustUnmarshalData(index int, obj interface{}) {
	err := json.Unmarshal(r.Data[index], obj)
	if err != nil {
		panic(err)
	}
}
