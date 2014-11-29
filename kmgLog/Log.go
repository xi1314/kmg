package kmgLog

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/bronze1man/kmg/kmgFile"
)

type Logger struct {
	LogWriter
}

func (l *Logger) Log(category string, msg string, obj interface{}) (err error) {
	return l.LogWriter.LogWrite(category, LogRow{
		Time: time.Now().Format(time.RFC3339),
		Msg:  msg,
		Obj:  obj,
	})
}

type LogWriter interface {
	LogWrite(category string, row LogRow) (err error)
}

type fileJsonLogWriter struct {
	logDir string
}

func (lw fileJsonLogWriter) LogWrite(category string, row LogRow) (err error) {
	b, err := json.Marshal(row)
	if err != nil {
		return err
	}
	toWrite := append(b, byte('\n'))
	err = kmgFile.AppendFile(filepath.Join(lw.logDir, category+".log"), toWrite)
	if err != nil {
		return err
	}
	return
}

//new file json log, will mkdir if dir not exist.
func NewFileJsonLogger(logDir string) *Logger {
	kmgFile.MustMkdirAll(logDir)
	return &Logger{
		LogWriter: fileJsonLogWriter{
			logDir: logDir,
		},
	}
}

type nullLogWriter struct {
}

func (nl nullLogWriter) LogWrite(category string, row LogRow) (err error) {
	return nil
}
func NewNullJsonLogger() *Logger {
	return &Logger{
		LogWriter: nullLogWriter{},
	}
}

var NullLogger = NewNullJsonLogger()

type stdoutLogWriter struct {
}

func (nl stdoutLogWriter) LogWrite(category string, row LogRow) (err error) {
	b, err := json.Marshal(row)
	if err != nil {
		return err
	}
	_, err = fmt.Printf("%s %s\n", category, b)
	return
}

var StdoutLogger = &Logger{
	LogWriter: stdoutLogWriter{},
}

var DefaultLogger *Logger = StdoutLogger

type LogRow struct {
	Time string
	Msg  string
	Obj  interface{}
}

// 写一条log, category是分类名 msg是消息信息 obj是需要调试的对象
// 要求obj可以被json序列化
// 如果使用文件序列化方案 category 是文件名.
func Log(category string, msg string, obj interface{}) {
	err := DefaultLogger.Log(category, msg, obj)
	if err != nil {
		panic(err)
	}
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
