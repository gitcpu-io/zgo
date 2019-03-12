package zgolog

import (
	"fmt"
	"git.zhugefang.com/gocore/zgo/config"
	"git.zhugefang.com/gocore/zgo/zgoutils"
	"github.com/go-stack/stack"
	log "github.com/sirupsen/logrus"
	"strings"
)

const (
	project = "project"
	file    = "file"
)

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
	res      LogStorer
}

//var Log = Newzgolog()

func Newzgolog() Logger {
	z := &zgolog{
		Project:  config.Project,
		LogLevel: config.Loglevel,
		Entry:    log.New(),
		res:      NewLogStore("file", "/tmp", 1),
	}
	return z
}

func InitLog(project, label, dbType string, start int) *zgolog {
	res := NewLogStore(dbType, label, start)
	return &zgolog{
		Project:  project,
		LogLevel: config.Loglevel,
		Entry:    log.New(),
		res:      res,
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
		res:      NewLogStore("file", "/tmp", 1),
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
	default:
		format := new(log.JSONFormatter)
		format.TimestampFormat = "2006-01-02 15:04:05"
		logger.Level = log.DebugLevel
		logger.Formatter = format
	}
	return nil
}

//func (z *zgolog) WithField(key string, value interface{}) *log.Entry {
//	return z.withCaller().WithField(key, value)
//}
//
//func (z *zgolog) WithFields(fs log.Fields) *log.Entry {
//	return z.withCaller().WithFields(fs)
//}

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

	//
	//k, _ := zgokafka.GetKafka("kafka_label_bj")
	//_,err = k.Producer(context.TODO(), "kafka_label_bj", []byte(fmt.Sprint(args...)))
	//if err != nil {
	//	z.withCaller().Error(args...)
	//	return
	//}

	en, value := z.withCaller()

	//f := zgofile.NewLocal(".")
	bstr := strings.Builder{}
	bstr.WriteString("(" + config.Project + ")")
	bstr.WriteString("(" + value.(string) + ")")
	bstr.WriteString("(error)")
	bstr.WriteString(fmt.Sprint(args...))
	bstr.WriteString(zgoutils.Utils.FormatFromUnixTime(-1))
	bstr.WriteString("\r\n")
	//input := strings.NewReader(bstr.String())
	//_, err := f.Append("/"+zgoutils.Utils.FormatFromUnixTimeShort(-1)+"/kafka_label_bj.txt", input)

	_, err := z.res.Save(z.Project, []byte(bstr.String()))
	if err != nil {
		en.Error(args...)
		return
	}

	en.Error(args...)
}

func (z *zgolog) Info(args ...interface{}) {
	//z.withCaller().Info(args...)
	en, _ := z.withCaller()
	en.Info(args...)
}

func (z *zgolog) Print(args ...interface{}) {
	//z.withCaller().Print(args...)
	en, _ := z.withCaller()
	en.Print(args...)
}

func (z *zgolog) Warn(args ...interface{}) {
	//z.withCaller().Warn(args...)
	en, _ := z.withCaller()
	en.Warn(args...)
}

func (z *zgolog) Debug(args ...interface{}) {
	//z.withCaller().Debug(args...)
	en, _ := z.withCaller()
	en.Debug(args...)
}

func (z *zgolog) Errorf(format string, args ...interface{}) {
	//z.withCaller().Errorf(format, args...)
	en, _ := z.withCaller()
	en.Errorf(format, args...)
}

func (z *zgolog) Infof(format string, args ...interface{}) {
	//z.withCaller().Infof(format, args...)
	en, _ := z.withCaller()
	en.Infof(format, args...)
}

func (z *zgolog) Printf(format string, args ...interface{}) {
	//z.withCaller().Printf(format, args...)
	en, _ := z.withCaller()
	en.Printf(format, args...)

}

func (z *zgolog) Warnf(format string, args ...interface{}) {
	//z.withCaller().Warnf(format, args...)
	en, _ := z.withCaller()
	en.Warnf(format, args...)
}

func (z *zgolog) Debugf(format string, args ...interface{}) {
	//z.withCaller().Debugf(format, args...)
	en, _ := z.withCaller()
	en.Debugf(format, args...)
}
