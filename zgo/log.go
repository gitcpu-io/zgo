package zgo

import (
	"git.zhugefang.com/gocore/zgo.git/logic/zgo_log"
)

//init Log

var Log logger

func init() { //初始化Log
	Log = zgo_log.Newzglog()
}

type logger interface {
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
