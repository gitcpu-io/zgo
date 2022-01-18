package zgoes

import (
  "context"
  "fmt"
  "github.com/gitcpu-io/zgo/config"
  "github.com/gitcpu-io/zgo/zgoutils"
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
    Host:     "http://101.201.28.195",
    Port:     9200,
    ConnSize: 50,
    PoolSize: 20000,
  }
  cd_bj2 := config.ConnDetail{
    C:        "北京主库-----es2",
    Host:     "http://101.201.28.195",
    Port:     9200,
    ConnSize: 50,
    PoolSize: 20000,
  }
  cd_sh := config.ConnDetail{
    C:        "上海主库-----es",
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

  sellR, _ := GetEs(label_sell)

  var replyChan = make(chan int)
  var countChan = make(chan int)
  l := 100 //暴力测试50000个消息，时间10秒，本本的并发每秒5000

  count := []int{}
  total := []int{}
  stime := time.Now()

  for i := 0; i < l; i++ {
    go func(i int) {
      countChan <- i //统计开出去的goroutine
      if i%2 == 0 {
        //ch := getMongo(label_sh, clientBj, i)
        ch := search(label_sell, sellR, i)
        reply := <-ch
        replyChan <- reply

      } else {
        //ch := getMongo(label_bj,clientSh,i)
        //ch := createMongo(label_bj, clientSh, i)
        ch := search(label_sell, sellR, i)
        reply := <-ch
        replyChan <- reply
      }
    }(i)
  }

  go func() {
    for v := range replyChan {
      if v != 10001 { //10001表示超时
        count = append(count, v) //成功数
      }
    }
  }()

  go func() {
    for v := range countChan { //总共的goroutine
      total = append(total, v)
    }
  }()

  for _, v := range count {
    if v != 1 {
      fmt.Println("有不成功的")
    }
  }

  for {
    if len(count) == l {
      var timeLen time.Duration
      timeLen = time.Now().Sub(stime)

      fmt.Printf("总消耗时间：%s, 成功：%d, 总共开出来的goroutine：%d\n", timeLen, len(count), len(total))
      break
    }

    select {
    case <-time.Tick(time.Duration(1000 * time.Millisecond)):
      fmt.Println("处理进度每1000毫秒", len(count))

    }
  }

}

func search(label string, client *zgoes, i int) chan int {
  //还需要一个上下文用来控制开出去的goroutine是否超时
  ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
  defer cancel()
  //输入参数：上下文ctx，mongoChan里面是client的连接，args具体的查询操作参数
  args := map[string]interface{}{}
  index := "active_bj_house_sell"
  table := "spider"
  dsl := `{
	  "query": {
   		 "bool" : {
    	    "filter": {
        	"term" : { "_id" : "5065" }
     		 }
   		 }
  		}
	}`
  result, err := client.SearchDsl(ctx, index, table, dsl, args)
  if err != nil {
    panic(err)
  }

  out := make(chan int, 1)
  select {
  case <-ctx.Done():
    fmt.Println("超时")
    out <- 10001
    return out
  default:
    _, err := zgoutils.Utils.Marshal(result)
    if err != nil {
      panic(err)
    }
    //fmt.Println(string(bytes), err, "---from mongo successful---")
    out <- 1
  }

  return out

}

//func TestGetById(t *testing.T) {
//	hsm := make(map[string][]*config.ConnDetail)
//	cd_bj := config.ConnDetail{
//		C:        "北京主库-----es1",
//		Uri:      "http://101.201.119.240:9200",
//		Host:     "http://101.201.28.195",
//		Port:     9200,
//		ConnSize: 50,
//		PoolSize: 20000,
//	}
//	cd_bj2 := config.ConnDetail{
//		C:        "北京主库-----es2",
//		Uri:      "http://101.201.119.240:9200",
//		Host:     "http://101.201.28.195",
//		Port:     9200,
//		ConnSize: 50,
//		PoolSize: 20000,
//	}
//	cd_sh := config.ConnDetail{
//		C:        "上海主库-----es",
//		Uri:      "http://101.201.119.240:9200",
//		Host:     "http://101.201.28.195",
//		Port:     9200,
//		ConnSize: 50,
//		PoolSize: 20000,
//	}
//	var s1 []*config.ConnDetail
//	var s2 []*config.ConnDetail
//	s1 = append(s1, &cd_bj, &cd_bj2)
//	s2 = append(s2, &cd_sh)
//	hsm = map[string][]*config.ConnDetail{
//		label_sell: s1,
//		label_rent: s2,
//	}
//
//	InitEs(hsm)
//	sellR, _ := GetEs(label_sell)
//	args := map[string]interface{}{}
//	index := "active_bj_house_sell"
//	table := "spider"
//	dsl := `{
//	  "query": {
//   		 "bool" : {
//    	    "filter": {
//        	"term" : { "_id" : "5055" }
//     		 }
//   		 }
//  		}
//	}`
//	res, _ := sellR.SearchDsl(context.TODO(), index, table, dsl, args)
//	s, _ := zgoutils.NewUtils().Marshal(res)
//	fmt.Print(string(s))
//}
