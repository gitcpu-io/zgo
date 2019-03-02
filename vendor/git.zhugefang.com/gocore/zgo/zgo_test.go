package zgo

import (
	"fmt"
	"testing"
	"time"
)

func TestEngine(t *testing.T) {

	Engine(&Options{
		Env: "local",
		Nsq: []string{
			"nsq_label_bj",
			//"label_sh",
		},
		Es: []string{
			"es_new_write",
			//"label_sh",
		},
	})

	for {
		select {
		case <-time.Tick(time.Duration(5) * time.Second):
			fmt.Println("start engine for test")
		}
	}
}
