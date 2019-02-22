package zgo_log

import (
	"testing"
)

func TestZlog_Info(t *testing.T) {

	l := Newzlog()
	l.InitLog("newProject", "debug")
	l.Info("33333")
}
