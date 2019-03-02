package zgoredis

import (
	"context"
	"fmt"
	"git.zhugefang.com/gocore/zgo/config"
	"github.com/json-iterator/go"
	"testing"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	label_bj = "redis_label_bj"
	label_sh = "redis_label_sh"
)

func TestRedisGet(t *testing.T) {
	hsm := make(map[string][]*config.ConnDetail)
	cd_bj := config.ConnDetail{
		C:        "北京主库-----redis1",
		Host:     "localhost",
		Port:     6379,
		ConnSize: 100,
		PoolSize: 200,
		Username: "",
		Password: "",
		Db:       0,
	}
	cd_bj2 := config.ConnDetail{
		C:        "北京从库-----redis2",
		Host:     "localhost",
		Port:     6379,
		ConnSize: 10,
		PoolSize: 200,
		Username: "",
		Password: "",
		Db:       0,
	}
	cd_sh := config.ConnDetail{
		C:        "上海主库-----redis",
		Host:     "localhost",
		Port:     6379,
		ConnSize: 10,
		PoolSize: 200,
		Username: "",
		Password: "",
		Db:       5,
	}
	var s1 []*config.ConnDetail
	var s2 []*config.ConnDetail
	s1 = append(s1, &cd_bj, &cd_bj2)
	s2 = append(s2, &cd_sh)
	hsm = map[string][]*config.ConnDetail{
		label_bj: s1,
		label_sh: s2,
	}

	InitRedis(hsm) //测试时表示使用redis，在zgo_start中使用一次

	clientLocal, err := GetRedis(label_bj)
	clientSpider, err := GetRedis(label_sh)

	fmt.Println(clientLocal)
	fmt.Println(clientSpider)

	if err != nil {
		panic(err)
	}

	//LpushCheck(label_bj, clientLocal, 0)
	//
	//getSet(label_sh, clientSpider, 0)

	var replyChan = make(chan int)
	var countChan = make(chan int)
	l := 100000 //暴力测试50000个消息，时间10秒，本本的并发每秒5000

	count := []int{}
	total := []int{}
	stime := time.Now()

	for i := 0; i < l; i++ {
		go func(i int) {
			countChan <- i //统计开出去的goroutine
			if i%2 == 0 {
				ch := getSet(label_bj, clientLocal, i)
				reply := <-ch
				replyChan <- reply

			} else {
				ch := getSet(label_sh, clientSpider, i)
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

func getSet(label string, client *zgoredis, i int) chan int {
	//还需要一个上下文用来控制开出去的goroutine是否超时
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var fooVal string
	_, err := client.Do(ctx, &fooVal, "GET", "foo")
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
		fmt.Println(fooVal)
		out <- 1
	}

	return out
}

func LpushCheck(label string, client *zgoredis, i int) chan int {
	//还需要一个上下文用来控制开出去的goroutine是否超时
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	result, err := client.Do(ctx, nil, "Lpush", "china_online_list", "somebody111")
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
		_, err := json.Marshal(result)
		if err != nil {
			panic(err)
		}
		out <- 1
	}

	return out
}
