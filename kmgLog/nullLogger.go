package kmgLog

var NullLogger Logger = nullLogger{}

type nullLogger struct {
}

func (nl nullLogger) Log(category string, data ...interface{}) {
	return
}
