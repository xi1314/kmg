package kmgLog

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/bronze1man/kmg/kmgFile"
	"os"
)

type Logger interface {
	Log(category string, objs ...interface{})
}

// 写一条log, category是分类名 data是需要调试的对象
// 实际格式,用户可以随意设定
// data 必须为可以json.Marshal的数据
// 错误处理方式为 输出到stderr (不用panic,是因为log坏了,程序应该还可以继续执行没那么重要,不忽略,是因为忽略错误后果更严重)
func Log(category string, data ...interface{}) {
	DefaultLogger.Log(category, data...)
}

type LogRow struct {
	Cat  string
	Time string
	Data []json.RawMessage
}

var DefaultLogger Logger = StdoutLogger

var StdoutLogger Logger = stdoutLogger{}

type stdoutLogger struct {
}

func (nl stdoutLogger) Log(category string, data ...interface{}) {
	b, err := logToJson(category, data)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[StdoutLogger] logToJson fail", err)
		return
	}
	_, err = fmt.Printf("%s\n", b)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[StdoutLogger] printf fail", err)
		return
	}
	return
}

var NullLogger Logger = nullLogger{}

type nullLogger struct {
}

func (nl nullLogger) Log(category string, data ...interface{}) {
	return
}

type fileLoger struct {
	logDir string
}

func (lw fileLoger) Log(category string, data ...interface{}) {
	b, err := logToJson(category, data)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[fileLoger] logToJson fail", err)
		return
	}
	toWrite := append(b, byte('\n'))
	err = kmgFile.AppendFile(filepath.Join(lw.logDir, category+".log"), toWrite)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[fileLoger] logToJson fail", err)
		return
	}
}

type mulitLogger []Logger

func (ml mulitLogger) Log(category string, data ...interface{}) {
	for i := range ml {
		ml[i].Log(category, data...)
	}
}

func MultiLogger(loggers ...Logger) Logger {
	return mulitLogger(loggers)
}

type threadLogger struct {
	logger Logger
}

func (ml threadLogger) Log(category string, data ...interface{}) {
	go ml.Log(category, data...)
}

func ThreadLogger(logger Logger) Logger {
	return threadLogger{logger: logger}
}

//new file log, will mkdir if dir not exist.
// usage:
// 		kmgLog.DefaultLogger = kmgLog.NewFileLogger("log")
func NewFileLogger(logDir string) Logger {
	kmgFile.MustMkdirAll(logDir)
	return fileLoger{
		logDir: logDir,
	}
}

func NewStdoutAndFileLogger(logDir string) Logger {
	return MultiLogger(StdoutLogger, NewFileLogger(logDir))
}

func logToJson(category string, data []interface{}) (b []byte, err error) {
	logRow := LogRow{
		Cat:  category,
		Time: time.Now().Format(time.RFC3339),
		Data: make([]json.RawMessage, len(data)),
	}
	for i := range data {
		logRow.Data[i], err = json.Marshal(data[i])
		if err != nil {
			return
		}
	}
	return json.Marshal(logRow)
}

// 如果f发生错误,写一条log
func LogErrCallback(category string, context interface{}, f func() error) error {
	err := f()
	if err == nil {
		return nil
	}
	Log(category, err.Error(), context)
	return err
}
