package zgo_log

import (
	log "github.com/sirupsen/logrus"
)

type Entry *log.Entry

type zlog struct {
	E Entry
}

func Newzlog() *zlog {
	return &zlog{
		E: log.NewEntry(log.New()),
	}
}

func (z *zlog) InitLog() Entry {

	z.E.Logger.SetFormatter(&log.JSONFormatter{
		//PrettyPrint: true,
		//DisableTimestamp: false,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	z.E.Logger.SetReportCaller(true)

	return z.E
}
