package zgolog

import (
	"fmt"
	"git.zhugefang.com/gocore/zgo.git/config"
	"github.com/go-stack/stack"
	log "github.com/sirupsen/logrus"
)

const (
	project = "project"
	file    = "file"
)

type Logger interface {
	NewLog(projectName string, logLevel string) *zgolog
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

type zgolog struct {
	Project  string
	LogLevel string
	Entry    *log.Logger
}

func Newzgolog() Logger {
	z := &zgolog{
		Project: config.Project,
		Entry:   log.New(),
	}
	return z
}

func (z *zgolog) NewLog(project string, logLevel string) *zgolog {
	return &zgolog{
		Project:  project,
		LogLevel: logLevel,
	}
}

// debug: 使用text格式, Level是Debug, 打印所有级别
// not debug: 使用json格式, level是Info, 不打印Debug级别
func (z *zgolog) SetDebug(level string) *log.Logger {
	l, _ := log.ParseLevel(level)
	logger := z.Entry
	switch l {
	case log.DebugLevel:
		format := new(log.TextFormatter)
		format.ForceColors = true
		format.FullTimestamp = true
		format.TimestampFormat = "2006-01-02 15:04:05"
		logger.Level = log.DebugLevel
		logger.Formatter = format
	case log.ErrorLevel:
		format := new(log.JSONFormatter)
		format.TimestampFormat = "2006-01-02 15:04:05"
		logger.Level = log.ErrorLevel
		logger.Formatter = format
	case log.WarnLevel:
		format := new(log.JSONFormatter)
		format.TimestampFormat = "2006-01-02 15:04:05"
		logger.Level = log.WarnLevel
		logger.Formatter = format
	case log.InfoLevel:
		format := new(log.JSONFormatter)
		format.TimestampFormat = "2006-01-02 15:04:05"
		logger.Level = log.InfoLevel
		logger.Formatter = format
	case log.PanicLevel:
		format := new(log.JSONFormatter)
		format.TimestampFormat = "2006-01-02 15:04:05"
		logger.Level = log.InfoLevel
		logger.Formatter = format
	case log.FatalLevel:
		format := new(log.JSONFormatter)
		format.TimestampFormat = "2006-01-02 15:04:05"
		logger.Level = log.InfoLevel
		logger.Formatter = format
	}
	return nil
}

func (z *zgolog) WithField(key string, value interface{}) *log.Entry {
	return z.withCaller().WithField(key, value)
}

func (z *zgolog) WithFields(fs log.Fields) *log.Entry {
	return z.withCaller().WithFields(fs)
}

func (z *zgolog) withCaller() *log.Entry {
	var value interface{}
	z.SetDebug(config.Loglevel)
	if config.Loglevel == "debug" {
		// 支持goland点击跳转
		value = fmt.Sprintf(" %+v:", stack.Caller(1))
	} else {
		value = fmt.Sprintf("%+v", stack.Caller(2))
	}
	en := z.Entry.WithFields(log.Fields{
		project: config.Project,
		file:    value,
	})
	return en
}

func (z *zgolog) Error(args ...interface{}) {
	z.withCaller().Error(args...)
}

func (z *zgolog) Info(args ...interface{}) {
	z.withCaller().Info(args...)
}

func (z *zgolog) Print(args ...interface{}) {
	z.withCaller().Print(args...)
}

func (z *zgolog) Warn(args ...interface{}) {
	z.withCaller().Warn(args...)
}

func (z *zgolog) Debug(args ...interface{}) {
	z.withCaller().Debug(args...)
}

func (z *zgolog) Errorf(format string, args ...interface{}) {
	z.withCaller().Errorf(format, args...)
}

func (z *zgolog) Infof(format string, args ...interface{}) {
	z.withCaller().Infof(format, args...)
}

func (z *zgolog) Printf(format string, args ...interface{}) {
	z.withCaller().Printf(format, args...)
}

func (z *zgolog) Warnf(format string, args ...interface{}) {
	z.withCaller().Warnf(format, args...)
}

func (z *zgolog) Debugf(format string, args ...interface{}) {
	z.withCaller().Debugf(format, args...)
}
