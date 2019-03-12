package zgo

import (
	"context"
	"fmt"
	"git.zhugefang.com/gocore/zgo/config"
	"git.zhugefang.com/gocore/zgo/zgofile"
	"git.zhugefang.com/gocore/zgo/zgokafka"
	"git.zhugefang.com/gocore/zgo/zgolog"
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

func StartQueue() {
	go func() {
		for v := range zgolog.LbodyCh {
			//fmt.Println("from chan ============:", v)
			lbb, err := zgoutils.Utils.Marshal(v)
			_, err = LogStore.Deal(config.Project, lbb)
			if err != nil {
				fmt.Println("error logstore")
			}
		}
	}()

}

type LogStorer interface {
	Deal(topic string, body []byte) (int, error)
}

type logStore struct {
	DbType string `json:"dbType"`
	Label  string `json:"label"`
	Start  int    `json:"start"`
}

func NewLogStore(dbType string, label string, start int) LogStorer {
	nls := &logStore{
		DbType: dbType,
		Label:  label,
		Start:  start,
	}
	return nls
}

func (ls *logStore) Deal(topic string, body []byte) (int, error) {
	//fmt.Println(topic, ls.DbType, ls.Label, ls.Start, "=====当前日志存储方式======")

	if ls.Start != 1 {
		return 0, nil
	}

	switch ls.DbType {
	case "nsq":
		n, err := zgonsq.GetNsq(ls.Label)
		if err != nil {
			fmt.Println(ls.Label, "==nsq==", err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		ui8, err := n.Producer(ctx, topic, body)
		if err != nil {
			fmt.Println(ls.Label, "==nsq==", err)
		}
		pint := int(<-ui8)
		return pint, err

	case "kafka":
		k, err := zgokafka.GetKafka(ls.Label)
		if err != nil {
			fmt.Println(ls.Label, "==kafka==", err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		ui8, err := k.Producer(ctx, topic, body)
		if err != nil {
			fmt.Println(ls.Label, "==kafka==", err)
		}
		pint := int(<-ui8)
		return pint, err

	case "file":
		input := strings.NewReader(string(body) + "\r\n")
		f := zgofile.NewLocal(ls.Label)
		pn, err := f.Append("/"+zgoutils.Utils.FormatFromUnixTimeShort(-1)+"/"+topic+".log", input)
		pint := int(pn)
		return pint, err
	}
	return 0, nil
}
