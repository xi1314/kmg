package kmgLog

import "github.com/bronze1man/kmg/kmgTime"

// 写一条log, category是分类名 data是需要调试的对象
// 实际格式,用户可以随意设定
// data 必须为可以json.Marshal的数据
// 错误处理方式为 输出到stderr (不用panic,是因为log坏了,程序应该还可以继续执行没那么重要,不忽略,是因为忽略错误后果更严重)
// 这个算是LogToRow的语法糖.
func Log(cat string, data ...interface{}) {
	logRow := LogRow{
		Cat:  cat,
		Time: kmgTime.NowFromDefaultNower(),
		Data: data,
		//Data: make([]json.RawMessage, len(data)),
	}
	defaultWriter(logRow)
	//defaultLogger(cat, data...)
}

func LogToRow(row LogRow) {
	defaultWriter(row)
}

var defaultWriter LogWriter = StdoutLogWriter

/*
example:
	kmgLog.SetLogWriter(kmgLog.StdoutLogWriter)
*/
func SetLogWriter(w LogWriter) {
	defaultWriter = w
}

func SetDefaultThreadFileLog(logDir string) {
	defaultWriter = ThreadLogWriter(NewFileLogWriter(logDir))
}

func SetDefaultStdoutAndFileLog(logDir string) {
	defaultWriter = ThreadLogWriter(MultiLogWriter(StdoutLogWriter, NewFileLogWriter(logDir)))

}
