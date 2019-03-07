package config

import (
	"testing"
)

/*
@Time : 2019-03-04 15:16
@Author : rubinus.chu
@File : etcd_test
@project: zgo
*/

func TestWatcher(t *testing.T) {
	c, _ := CreateClient()
	key := "zgo/nsq/label_bj"
	Watcher(c, key)
}
