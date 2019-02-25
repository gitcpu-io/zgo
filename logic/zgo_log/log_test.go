package zgo_log

import (
	"testing"
)

func TestZlog_Info(t *testing.T) {

	l := Newzglog()
	l.NewLog("newProject", "debug")
	l.Info("33333")
}

func BenchmarkZglog_Info(b *testing.B) {
	l := Newzglog()
	l.NewLog("newProject", "debug")
	l.Info("33333")
}
