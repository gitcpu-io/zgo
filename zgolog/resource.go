package zgolog

import (
	"context"
	"fmt"
	"git.zhugefang.com/gocore/zgo/zgofile"
	"git.zhugefang.com/gocore/zgo/zgokafka"
	"git.zhugefang.com/gocore/zgo/zgonsq"
	"git.zhugefang.com/gocore/zgo/zgoutils"
	"strings"
)

/*
@Time : 2019-03-11 20:01
@Author : rubinus.chu
@File : resource
@project: zgo
*/

type LogStorer interface {
	Save(project string, body []byte) (int, error)
}

type logStore struct {
	DbType string `json:"dbType"`
	Label  string `json:"label"`
	Start  int    `json:"start"`
}

func NewLogStore(dbType string, label string, start int) LogStorer {
	return &logStore{DbType: dbType, Label: label, Start: start}
}

func (ls *logStore) Save(project string, body []byte) (int, error) {
	if ls.Start != 1 {
		return 0, nil
	}
	fmt.Println(project, ls.DbType, ls.Label, ls.Start, "===========")

	switch ls.DbType {
	case "nsq":
		n, _ := zgonsq.GetNsq(ls.Label)
		ui8, err := n.Producer(context.TODO(), project, body)
		pint := int(<-ui8)
		return pint, err

	case "kafka":
		k, _ := zgokafka.GetKafka(ls.Label)
		ui8, err := k.Producer(context.TODO(), project, body)
		pint := int(<-ui8)
		return pint, err

	case "file":
		input := strings.NewReader(string(body))
		f := zgofile.NewLocal(ls.Label)
		pn, err := f.Append("/"+zgoutils.Utils.FormatFromUnixTimeShort(-1)+"/"+project+".txt", input)
		pint := int(pn)
		return pint, err
	}
	return 0, nil
}
