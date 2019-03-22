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
	New() *zgolog
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type zgolog struct {
	Project string
	Entry   *log.Logger
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
		Project: project,
		Entry:   log.New(),
	}
}

func (z *zgolog) New() *zgolog {
	return &zgolog{
		Project: config.Conf.Project,
		Entry:   log.New(),
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
	ll := config.Levels[config.Conf.Log.LogLevel]
	z.SetDebug(ll)
	if config.Conf.Log.LogLevel == 0 {
		// 支持goland点击跳转
		value = fmt.Sprintf(" %+v:", stack.Caller(2))
	} else {
		value = fmt.Sprintf("%+v", stack.Caller(2))
	}
	p := z.Project
	if p == "" {
		p = "zgo"
	}
	en := z.Entry.WithFields(log.Fields{
		project: p,
		file:    value,
	})
	return en, value
}

func (z *zgolog) Debug(args ...interface{}) {
	en, value := z.withCaller()

	//fmt.Println(config.Conf.Log.LogLevel, config.Debug, "---debug")
	if config.Conf.Log.LogLevel <= config.Debug {
		lb := logBody{
			Project: z.Project,
			File:    value.(string),
			Msg:     fmt.Sprint(args...),
			Time:    zgoutils.Utils.FormatFromUnixTime(-1),
			Level:   config.Levels[config.Debug],
		}
		LbodyCh <- &lb

	}

	en.Debug(args...)
}

func (z *zgolog) Info(args ...interface{}) {
	en, value := z.withCaller()
	//fmt.Println(config.Conf.Log.LogLevel, config.Info, "---info")

	if config.Conf.Log.LogLevel <= config.Info {
		lb := logBody{
			Project: z.Project,
			File:    value.(string),
			Msg:     fmt.Sprint(args...),
			Time:    zgoutils.Utils.FormatFromUnixTime(-1),
			Level:   config.Levels[config.Info],
		}
		LbodyCh <- &lb
	}

	en.Info(args...)
}

func (z *zgolog) Warn(args ...interface{}) {
	en, value := z.withCaller()
	//fmt.Println(config.Conf.Log.LogLevel, config.Warn, "---warn")

	if config.Conf.Log.LogLevel <= config.Warn {
		lb := logBody{
			Project: z.Project,
			File:    value.(string),
			Msg:     fmt.Sprint(args...),
			Time:    zgoutils.Utils.FormatFromUnixTime(-1),
			Level:   config.Levels[config.Warn],
		}
		LbodyCh <- &lb
	}

	en.Warn(args...)
}

func (z *zgolog) Error(args ...interface{}) {

	en, value := z.withCaller()
	//fmt.Println(config.Conf.Log.LogLevel, config.Error, "---error")

	if config.Conf.Log.LogLevel <= config.Error {
		lb := logBody{
			Project: z.Project,
			File:    value.(string),
			Msg:     fmt.Sprint(args...),
			Time:    zgoutils.Utils.FormatFromUnixTime(-1),
			Level:   config.Levels[config.Error],
		}
		LbodyCh <- &lb

	}

	en.Error(args...)
}

func (z *zgolog) Debugf(format string, args ...interface{}) {
	en, value := z.withCaller()
	if config.Conf.Log.LogLevel <= config.Debug {
		lb := logBody{
			Project: z.Project,
			File:    value.(string),
			Msg:     fmt.Sprint(args...),
			Time:    zgoutils.Utils.FormatFromUnixTime(-1),
			Level:   config.Levels[config.Debug],
		}
		LbodyCh <- &lb
	}

	en.Debugf(format, args...)
}

func (z *zgolog) Infof(format string, args ...interface{}) {
	en, value := z.withCaller()
	if config.Conf.Log.LogLevel <= config.Info {
		lb := logBody{
			Project: z.Project,
			File:    value.(string),
			Msg:     fmt.Sprint(args...),
			Time:    zgoutils.Utils.FormatFromUnixTime(-1),
			Level:   config.Levels[config.Info],
		}
		LbodyCh <- &lb
	}

	en.Infof(format, args...)
}

func (z *zgolog) Warnf(format string, args ...interface{}) {
	en, value := z.withCaller()
	if config.Conf.Log.LogLevel <= config.Warn {
		lb := logBody{
			Project: z.Project,
			File:    value.(string),
			Msg:     fmt.Sprint(args...),
			Time:    zgoutils.Utils.FormatFromUnixTime(-1),
			Level:   config.Levels[config.Warn],
		}
		LbodyCh <- &lb
	}

	en.Warnf(format, args...)
}

func (z *zgolog) Errorf(format string, args ...interface{}) {
	en, value := z.withCaller()
	if config.Conf.Log.LogLevel <= config.Error {
		lb := logBody{
			Project: z.Project,
			File:    value.(string),
			Msg:     fmt.Sprint(args...),
			Time:    zgoutils.Utils.FormatFromUnixTime(-1),
			Level:   config.Levels[config.Error],
		}
		LbodyCh <- &lb
	}

	en.Errorf(format, args...)
}
