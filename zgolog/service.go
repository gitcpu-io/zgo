package zgolog

import (
	"errors"
	"fmt"
	"github.com/go-stack/stack"
	log "github.com/sirupsen/logrus"
)

const (
	project = "project"
	file    = "file"
)

type Logger interface {
	NewLog(projectName string, logLevel string) error
	Error(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Warn(args ...interface{})
	Debug(args ...interface{})
	Errorf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
}

var zlogger *Logger

type zglog struct {
	project  string
	logLevel log.Level
}

func Newzglog(project, loglevel string) *zglog {
	var l log.Level
	switch loglevel {
	case "info":
		l = log.InfoLevel
	case "debug":
		l = log.DebugLevel
	}
	return &zglog{
		project:  project,
		logLevel: l,
	}
}

func (z *zglog) NewLog(projectName string, logLevel string) error {
	//zlogger = &log.Logger{
	//	Out:       os.Stdout,
	//	Formatter: nil,
	//	Hooks:     make(log.LevelHooks),
	//	Level:     0,
	//}
	z.project = projectName
	err := z.setDebug(logLevel)
	if err != nil {
		return err
	}
	return nil
}

// debug: 使用text格式, Level是Debug, 打印所有级别
// not debug: 使用json格式, level是Info, 不打印Debug级别
func (z *zglog) setDebug(level string) error {
	l, err := log.ParseLevel(level)
	if err != nil {
		return errors.New("请输入有效的日志等级")
	}
	z.logLevel = l

	switch l {
	case log.DebugLevel:
		format := new(log.TextFormatter)
		format.ForceColors = true
		format.FullTimestamp = true
		format.TimestampFormat = "2006-01-02 15:04:05"
		//log.Level = log.DebugLevel
		//logger.Formatter = format
	case log.ErrorLevel:
		format := new(log.JSONFormatter)
		format.TimestampFormat = "2006-01-02 15:04:05"
		//logger.Level = log.ErrorLevel
		//logger.Formatter = format
	case log.WarnLevel:
		format := new(log.JSONFormatter)
		format.TimestampFormat = "2006-01-02 15:04:05"
		//logger.Level = log.WarnLevel
		//logger.Formatter = format
	case log.InfoLevel:
		format := new(log.JSONFormatter)
		format.TimestampFormat = "2006-01-02 15:04:05"
		//logger.Level = log.InfoLevel
		//logger.Formatter = format
	}
	return nil
}

func (z *zglog) WithField(key string, value interface{}) *log.Entry {
	return z.withCaller().WithField(key, value)
}

func (z *zglog) WithFields(fs log.Fields) *log.Entry {
	return z.withCaller().WithFields(fs)
}

func (z *zglog) withCaller() *log.Entry {
	var key = z.project
	var value interface{}
	if z.logLevel == log.DebugLevel {
		// 支持goland点击跳转
		value = fmt.Sprintf(" %+v:", stack.Caller(2))
	} else {
		value = fmt.Sprintf("%+v", stack.Caller(2))
	}

	return log.WithFields(log.Fields{
		project: key,
		file:    value,
	})
}

func (z *zglog) Error(args ...interface{}) {
	z.withCaller().Error(args...)
}

func (z *zglog) Info(args ...interface{}) {
	z.withCaller().Info(args...)
}

func (z *zglog) Print(args ...interface{}) {
	z.withCaller().Print(args...)
}

func (z *zglog) Warn(args ...interface{}) {
	z.withCaller().Warn(args...)
}

func (z *zglog) Debug(args ...interface{}) {
	z.withCaller().Debug(args...)
}

func (z *zglog) Errorf(format string, args ...interface{}) {
	z.withCaller().Errorf(format, args...)
}

func (z *zglog) Infof(format string, args ...interface{}) {
	z.withCaller().Infof(format, args...)
}

func (z *zglog) Printf(format string, args ...interface{}) {
	z.withCaller().Printf(format, args...)
}

func (z *zglog) Warnf(format string, args ...interface{}) {
	z.withCaller().Warnf(format, args...)
}

func (z *zglog) Debugf(format string, args ...interface{}) {
	z.withCaller().Debugf(format, args...)
}
