package kmgLog

import (
	"fmt"
	"os"
)

var StdoutLogger Logger = NewLogger(stdoutLogWriter{})
var StdoutLogWriter LogWriter = stdoutLogWriter{}

type stdoutLogWriter struct {
}

func (nl stdoutLogWriter) LogWrite(r LogRow) {
	b, err := r.Marshal()
	if err != nil {
		fmt.Fprintln(os.Stderr, "[StdoutLogger] Marshal fail", err)
		return
	}
	_, err = fmt.Printf("%s\n", b)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[StdoutLogger] printf fail", err)
		return
	}
	return
}
