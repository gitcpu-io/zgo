package zgo

import "git.zhugefang.com/gocore/zgo.git/logic/zgo_log"

//init Log
var Log logger

func init() { //初始化Log
	Log = zgo_log.Newzlog()
}

type logger interface {
	InitLog() zgo_log.Entry
}
