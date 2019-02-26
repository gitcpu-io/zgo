package zgoes

import (
	"context"
	"fmt"
	"git.zhugefang.com/gocore/zgo.git/config"
	"testing"
	"time"
)

const (
	label_new  = "new"
	label_sell = "sell"
	label_rent = "rents"
)

func TestEsSearch(t *testing.T) {
	hsm := make(map[string][]config.ConnDetail)
	new := config.ConnDetail{
		C:        "北京主库-----es1",
		Uri:      "http://101.201.28.195:9200",
		Host:     "http://101.201.28.195",
		Port:     9200,
		ConnSize: 50,
		PoolSize: 20000,
	}
	sell := config.ConnDetail{
		C:        "北京主库-----es2",
		Uri:      "http://101.201.28.195:9200",
		Host:     "http://101.201.28.195",
		Port:     9200,
		ConnSize: 50,
		PoolSize: 20000,
	}
	rent := config.ConnDetail{
		C:        "上海主库-----es",
		Uri:      "http://101.201.28.195:9200",
		Host:     "http://101.201.28.195",
		Port:     9200,
		ConnSize: 50,
		PoolSize: 20000,
	}
	var s1 []config.ConnDetail
	var s2 []config.ConnDetail
	var s3 []config.ConnDetail
	s1 = append(s1, new)
	s2 = append(s2, sell)
	s3 = append(s3, rent)
	hsm = map[string][]config.ConnDetail{
		label_new:  s1,
		label_sell: s2,
		label_rent: s3,
	}
	InitEs(hsm)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	args := map[string]interface{}{}
	args["index"] = "active_bj_house_sell"
	args["table"] = "spider"
	args["dsl"] = `{"query": {"match_all": {}}}`

	sellR, _ := GetEs(label_sell)
	start := time.Now()
	out := make(chan interface{}, 1000)
	go func() {
		for i := 0;i < 1000 ; i++  {
			result, err := sellR.Search(ctx, args)
			if err != nil{
				//fmt.Println(err)
				//panic(err)
				t.Fatal(err)
			}
			out <- result
		}
	}()
	//result, err := sellR.Search(ctx, args)
	for i := 0;i < 1000 ; i++  {
		fmt.Println(<- out)

	}
	finish := time.Now()
	fmt.Println(finish.Sub(start))

}

//func BenchmarkGetEs(b *testing.B) {
//	hsm := make(map[string][]config.ConnDetail)
//	new := config.ConnDetail{
//		C:        "北京主库-----es1",
//		Uri:      "http://101.201.28.195:9200",
//		Host:     "http://101.201.28.195",
//		Port:     9200,
//		ConnSize: 50,
//		PoolSize: 20000,
//	}
//	sell := config.ConnDetail{
//		C:        "北京主库-----es2",
//		Uri:      "http://101.201.28.195:9200",
//		Host:     "http://101.201.28.195",
//		Port:     9200,
//		ConnSize: 50,
//		PoolSize: 20000,
//	}
//	rent := config.ConnDetail{
//		C:        "上海主库-----es",
//		Uri:      "http://101.201.28.195:9200",
//		Host:     "http://101.201.28.195",
//		Port:     9200,
//		ConnSize: 50,
//		PoolSize: 20000,
//	}
//	var s1 []config.ConnDetail
//	var s2 []config.ConnDetail
//	var s3 []config.ConnDetail
//	s1 = append(s1, new)
//	s2 = append(s2, sell)
//	s3 = append(s3, rent)
//	hsm = map[string][]config.ConnDetail{
//		label_new:  s1,
//		label_sell: s2,
//		label_rent: s3,
//	}
//	InitEs(hsm)
//	sellR, _ := GetEs(label_sell)
//	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
//	args := map[string]interface{}{}
//	args["index"] = "active_bj_house_sell"
//	args["table"] = "spider"
//	args["dsl"] = `{"query": {"match_all": {}}}`
//	// 必须循环 b.N 次 。 这个数字 b.N 会在运行中调整，以便最终达到合适的时间消耗。方便计算出合理的数据。 （ 免得数据全部是 0 ）
//	for i:=0 ; i<b.N ; i++ {
//		_, err := sellR.Search(ctx, args)
//		if err != nil{
//			panic(err)
//		}
//	}
//}