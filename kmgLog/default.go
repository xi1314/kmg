package kmgLog

func SetDefaultThreadFileLog(logDir string) {
	DefaultLogger = NewThreadFileLogger(logDir)
}

func SetDefaultStdoutAndFileLog(logDir string) {
	DefaultLogger = NewThreadStdoutAndFileLogger(logDir)
}
