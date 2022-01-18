package zgolog

import "testing"

/*
@Time : 2019-03-06 14:56
@Author : rubinus.chu
@File : service_test
@project: zgo
*/

func TestNewzgolog(t *testing.T) {
  for i := 0; i < 10000; i++ {
    go func(i int) {
      zl := InitLog("")
      zl.Info("info")
      zl.Infof("infof:%s", i)
      zl.Error("error")
      zl.Errorf("errorf:%s", i)
      zl.Warn("warn")
      zl.Warnf("warnf:%s", i)
      zl.Debug("debug")
      zl.Debugf("debuf:%s", i)
    }(i)
  }

}
