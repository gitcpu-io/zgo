package zgolog

import (
	"fmt"
	"git.zhugefang.com/gocore/zgo/config"
	"git.zhugefang.com/gocore/zgo/zgoutils"
	"github.com/go-stack/stack"
	log "github.com/sirupsen/logrus"
)

const (
	project = "project"
	file    = "file"
)

var LbodyCh = make(chan *logBody, 2000)

type Logger interface {
	NewLog() *zgolog
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

type logBody struct {
	File    string `json:"file"`
	Project string `json:"project"`
	Time    string `json:"time"`
	Msg     string `json:"msg"`
	Level   string `json:"level"`
}

func InitLog(project string) *zgolog {

	return &zgolog{
		Project:  project,
		LogLevel: config.Loglevel,
		Entry:    log.New(),
	}
}

func (z *zgolog) NewLog() *zgolog {
	l := ""
	if z.LogLevel == "" {
		l = config.Loglevel
	} else {
		l = z.LogLevel
	}
	return &zgolog{
		Project:  config.Project,
		LogLevel: l,
		Entry:    log.New(),
	}
}

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
	default:
		format := new(log.JSONFormatter)
		format.TimestampFormat = "2006-01-02 15:04:05"
		logger.Level = log.DebugLevel
		logger.Formatter = format
	}
	return nil
}

func (z *zgolog) withCaller() (*log.Entry, interface{}) {
	var value interface{}
	z.SetDebug(config.Loglevel)
	if config.Loglevel == "debug" {
		// 支持goland点击跳转
		value = fmt.Sprintf(" %+v:", stack.Caller(2))
	} else {
		value = fmt.Sprintf("%+v", stack.Caller(2))
	}
	p := config.Project
	if p == "" {
		p = "zgo"
	}
	en := z.Entry.WithFields(log.Fields{
		project: p,
		file:    value,
	})
	return en, value
}

func (z *zgolog) Error(args ...interface{}) {

	en, value := z.withCaller()

	lb := logBody{
		Project: config.Project,
		File:    value.(string),
		Msg:     fmt.Sprint(args...),
		Time:    zgoutils.Utils.FormatFromUnixTime(-1),
		Level:   "error",
	}
	LbodyCh <- &lb
	en.Error(args...)
}

func (z *zgolog) Info(args ...interface{}) {
	en, value := z.withCaller()

	lb := logBody{
		Project: config.Project,
		File:    value.(string),
		Msg:     fmt.Sprint(args...),
		Time:    zgoutils.Utils.FormatFromUnixTime(-1),
		Level:   "info",
	}
	LbodyCh <- &lb
	en.Info(args...)
}

func (z *zgolog) Print(args ...interface{}) {
	en, value := z.withCaller()

	lb := logBody{
		Project: config.Project,
		File:    value.(string),
		Msg:     fmt.Sprint(args...),
		Time:    zgoutils.Utils.FormatFromUnixTime(-1),
		Level:   "print",
	}
	LbodyCh <- &lb
	en.Print(args...)
}

func (z *zgolog) Warn(args ...interface{}) {
	en, value := z.withCaller()

	lb := logBody{
		Project: config.Project,
		File:    value.(string),
		Msg:     fmt.Sprint(args...),
		Time:    zgoutils.Utils.FormatFromUnixTime(-1),
		Level:   "warn",
	}
	LbodyCh <- &lb
	en.Warn(args...)
}

func (z *zgolog) Debug(args ...interface{}) {
	en, value := z.withCaller()

	lb := logBody{
		Project: config.Project,
		File:    value.(string),
		Msg:     fmt.Sprint(args...),
		Time:    zgoutils.Utils.FormatFromUnixTime(-1),
		Level:   "debug",
	}
	LbodyCh <- &lb
	en.Debug(args...)
}

func (z *zgolog) Errorf(format string, args ...interface{}) {
	en, value := z.withCaller()

	lb := logBody{
		Project: config.Project,
		File:    value.(string),
		Msg:     fmt.Sprint(args...),
		Time:    zgoutils.Utils.FormatFromUnixTime(-1),
		Level:   "error",
	}
	LbodyCh <- &lb
	en.Errorf(format, args...)
}

func (z *zgolog) Infof(format string, args ...interface{}) {
	en, value := z.withCaller()

	lb := logBody{
		Project: config.Project,
		File:    value.(string),
		Msg:     fmt.Sprint(args...),
		Time:    zgoutils.Utils.FormatFromUnixTime(-1),
		Level:   "info",
	}
	LbodyCh <- &lb
	en.Infof(format, args...)
}

func (z *zgolog) Printf(format string, args ...interface{}) {
	en, value := z.withCaller()

	lb := logBody{
		Project: config.Project,
		File:    value.(string),
		Msg:     fmt.Sprint(args...),
		Time:    zgoutils.Utils.FormatFromUnixTime(-1),
		Level:   "print",
	}
	LbodyCh <- &lb
	en.Printf(format, args...)

}

func (z *zgolog) Warnf(format string, args ...interface{}) {
	en, value := z.withCaller()

	lb := logBody{
		Project: config.Project,
		File:    value.(string),
		Msg:     fmt.Sprint(args...),
		Time:    zgoutils.Utils.FormatFromUnixTime(-1),
		Level:   "warn",
	}
	LbodyCh <- &lb
	en.Warnf(format, args...)
}

func (z *zgolog) Debugf(format string, args ...interface{}) {
	en, value := z.withCaller()

	lb := logBody{
		Project: config.Project,
		File:    value.(string),
		Msg:     fmt.Sprint(args...),
		Time:    zgoutils.Utils.FormatFromUnixTime(-1),
		Level:   "debug",
	}
	LbodyCh <- &lb
	en.Debugf(format, args...)
}
