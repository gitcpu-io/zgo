package zgolog

import (
	"context"
	"fmt"
	"git.zhugefang.com/gocore/zgo/config"
	"git.zhugefang.com/gocore/zgo/zgofile"
	"git.zhugefang.com/gocore/zgo/zgokafka"
	"git.zhugefang.com/gocore/zgo/zgonsq"
	"git.zhugefang.com/gocore/zgo/zgoutils"
	"strings"
	"time"
)

/*
@Time : 2019-03-11 20:01
@Author : rubinus.chu
@File : resource
@project: zgo
*/

var LogWatch = make(chan *config.CacheConfig, 1)

var LogStore *logStore

const (
	dbtNsq   = "nsq"
	dbtKafka = "kafka"
	dbtFile  = "file"
)

type logStore struct {
	DbType string
	Label  string
	Start  int
}

func NewLogStore() *logStore {
	return &logStore{}
}

func StartLogStoreWatcher() {
	go func() {

		for {
			select {
			case v := <-LogWatch:
				LogStore.DbType = v.DbType
				LogStore.Label = v.Label
				LogStore.Start = v.Start
			}

		}

	}()
}

func (ls *logStore) StartQueue() {
	for v := range LbodyCh {
		if ls.Start == 1 {
			topic := config.Conf.Project

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			body, err := zgoutils.Utils.Marshal(v)
			if err != nil {
				fmt.Println("error logstore")
			}

			switch ls.DbType {

			case dbtNsq:
				nq, _ := zgonsq.GetNsq(ls.Label)
				_, err = nq.Producer(context.TODO(), topic, body)
				if err != nil {
					fmt.Println(ls.Label, "==nsq==", err)
				}

			case dbtKafka:
				kq, _ := zgokafka.GetKafka(ls.Label)

				_, err = kq.Producer(ctx, topic, body)
				if err != nil {
					fmt.Println(ls.Label, "==kafka==", err)
				}

			case dbtFile:
				if ls.Label == "" || !strings.HasPrefix(ls.Label, "/") {
					ls.Label = config.Conf.File.Home
				}
				f := zgofile.New(ls.Label)
				input := strings.NewReader(string(body) + "\r\n")
				_, err = f.Append("/"+zgoutils.Utils.FormatFromUnixTimeShort(-1)+"/"+topic+".log", input)
				if err != nil {
					fmt.Println(ls.Label, "==file==", err)
				}

			}

		}
	}

}
