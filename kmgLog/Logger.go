package kmgLog

type Logger interface {
	Log(cat string, objs ...interface{})
}

func NewLogger(w LogWriter) Logger {
	return loggerWithWriter{
		w: w,
	}
}

func NewThreadStdoutAndFileLogger(logDir string) Logger {
	return NewLogger(ThreadLogWriter(MultiLogWriter(StdoutLogWriter, NewFileLogWriter(logDir))))
}

type mulitLogger []LogWriter

func (ml mulitLogger) LogWrite(r LogRow) {
	for i := range ml {
		ml[i].LogWrite(r)
	}
}

func MultiLogWriter(loggers ...LogWriter) LogWriter {
	return mulitLogger(loggers)
}

type threadLogWriter struct {
	logWriter LogWriter
}

func (ml threadLogWriter) LogWrite(r LogRow) {
	go ml.logWriter.LogWrite(r)
}

func ThreadLogWriter(logger LogWriter) LogWriter {
	return threadLogWriter{logWriter: logger}
}
