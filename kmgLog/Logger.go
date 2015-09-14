package kmgLog

func MultiLogWriter(loggers ...LogWriter) LogWriter {
	return func(r LogRow) {
		for i := range loggers {
			loggers[i](r)
		}
	}
}

func ThreadLogWriter(logger LogWriter) LogWriter {
	return func(r LogRow) {
		go logger(r)
	}
}
