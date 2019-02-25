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
			"label_bj",
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
