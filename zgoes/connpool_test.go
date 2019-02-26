package zgoes

import (
	"context"
	"fmt"
	"git.zhugefang.com/gocore/zgo.git/config"
	"testing"
	"time"
)

const (
	label_sell = "es_label_sell"
	label_rent = "es_label_rent"
)

func TestEsSearch(t *testing.T) {
	hsm := make(map[string][]*config.ConnDetail)
	cd_bj := config.ConnDetail{
		C:        "北京主库-----es1",
		Uri:      "http://101.201.28.195:9200",
		Host:     "http://101.201.28.195",
		Port:     9200,
		ConnSize: 50,
		PoolSize: 20000,
	}
	cd_bj2 := config.ConnDetail{
		C:        "北京主库-----es2",
		Uri:      "http://101.201.28.195:9200",
		Host:     "http://101.201.28.195",
		Port:     9200,
		ConnSize: 50,
		PoolSize: 20000,
	}
	cd_sh := config.ConnDetail{
		C:        "上海主库-----es",
		Uri:      "http://101.201.28.195:9200",
		Host:     "http://101.201.28.195",
		Port:     9200,
		ConnSize: 50,
		PoolSize: 20000,
	}
	var s1 []*config.ConnDetail
	var s2 []*config.ConnDetail
	s1 = append(s1, &cd_bj, &cd_bj2)
	s2 = append(s2, &cd_sh)
	hsm = map[string][]*config.ConnDetail{
		label_sell: s1,
		label_rent: s2,
	}
	InitEs(hsm)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	args := map[string]interface{}{}
	args["index"] = "active_bj_house_sell"
	args["table"] = "spider"
	args["dsl"] = `{"query": {"match_all": {}}}`

	sellR, _ := GetEs(label_sell)
	result, err := sellR.Search(ctx, args)

	fmt.Print(result, err)

	//InitEs(map[string][]config.ConnDetail{
	//	label_sell: []string{"localhost:27017"},
	//	label_rent: []string{"localhost:27017"},
	//}) //测试时表示使用nsq，在zgo_start中使用一次

}
