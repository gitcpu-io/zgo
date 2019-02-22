package zgo_log

import (
	"fmt"
	"github.com/go-stack/stack"
	log "github.com/sirupsen/logrus"
	"os"
)

const (
	project = "project"
	file    = "file"
)

type Logger *log.Logger

var logger *log.Logger

type zlog struct {
	project  string
	logLevel log.Level
}

func Newzlog() *zlog {
	return &zlog{}
}

func (z *zlog) InitLog(projectName string, logLevel string) {
	z.project = projectName
	z.setDebug(logLevel)
	//return logger
}

// debug: 使用text格式, Level是Debug, 打印所有级别
// not debug: 使用json格式, level是Info, 不打印Debug级别
func (z *zlog) setDebug(level string) {
	l, err := log.ParseLevel(level)
	if err != nil {
		fmt.Errorf("请输入有效果的日志等级")
		return
	}
	z.logLevel = l

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
	}
}

func (z *zlog) WithField(key string, value interface{}) *log.Entry {
	return z.withCaller().WithField(key, value)
}

func (z *zlog) WithFields(fs log.Fields) *log.Entry {
	return z.withCaller().WithFields(fs)
}

func (z *zlog) withCaller() *log.Entry {
	var key = z.project
	var value interface{}
	if z.logLevel == log.DebugLevel {
		// 支持goland点击跳转
		value = fmt.Sprintf(" %+v:", stack.Caller(2))
	} else {
		value = fmt.Sprintf("%+v", stack.Caller(2))
	}

	return logger.WithFields(log.Fields{
		project: key,
		file:    value,
	})
}

func (z *zlog) Error(args ...interface{}) {
	z.withCaller().Error(args...)
}

func (z *zlog) Info(args ...interface{}) {
	z.withCaller().Info(args...)
}

func (z *zlog) Print(args ...interface{}) {
	z.withCaller().Print(args...)
}

func (z *zlog) Warn(args ...interface{}) {
	z.withCaller().Warn(args...)
}

func (z *zlog) Debug(args ...interface{}) {
	z.withCaller().Debug(args...)
}

func (z *zlog) Errorf(format string, args ...interface{}) {
	z.withCaller().Errorf(format, args...)
}

func (z *zlog) Infof(format string, args ...interface{}) {
	z.withCaller().Infof(format, args...)
}

func (z *zlog) Printf(format string, args ...interface{}) {
	z.withCaller().Printf(format, args...)
}

func (z *zlog) Warnf(format string, args ...interface{}) {
	z.withCaller().Warnf(format, args...)
}

func (z *zlog) Debugf(format string, args ...interface{}) {
	z.withCaller().Debugf(format, args...)
}

func init() {
	logger = &log.Logger{
		Out:       os.Stdout,
		Formatter: nil,
		Hooks:     make(log.LevelHooks),
		Level:     0,
	}
	//不能启用这个，否则会打印出zgo中的，而不是具体项目中的代码行数
	//logger.SetReportCaller(true)
}
