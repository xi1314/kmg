package kmgLog

func SetDefaultThreadFileLog(logDir string) {
	DefaultLogger = NewThreadFileLogger(logDir)
}
