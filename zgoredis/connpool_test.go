package zgoredis

import (
	"context"
	"fmt"
	"github.com/json-iterator/go"
	"testing"
	"time"
)

const (
	local1 = "local"
	spider = "spider"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func TestRedisGet(t *testing.T) {

	InitRedis(map[string][]string{
		local1: []string{
			"127.0.0.1:6379",
		},
		spider: []string{
			"127.0.0.1:6379",
		},
	}) //测试时表示使用nsq，在zgo_start中使用一次
	//
	clientLocal, err := GetRedis(local1)
	clientSpider, err := GetRedis(spider)

	if err != nil {
		panic(err)
	}

	fmt.Println(clientLocal)
	fmt.Println(clientSpider)

	if err != nil {
		panic(err)
	}

	//测试读取nsq数据，wait for sdk init connection
	time.Sleep(2 * time.Second)

	var replyChan = make(chan int)
	var countChan = make(chan int)
	l := 10 //暴力测试50000个消息，时间10秒，本本的并发每秒5000

	count := []int{}
	total := []int{}
	stime := time.Now()

	for i := 0; i < l; i++ {
		go func(i int) {
			countChan <- i //统计开出去的goroutine
			if i%2 == 0 {
				//ch := getSet(local1, clientLocal, i)
				ch := setSet(local1, clientLocal, i)
				reply := <-ch
				replyChan <- reply

			} else {
				ch := getSet(spider, clientSpider, i)
				//ch := setSet(spider, clientSpider, i)
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
	result, err := client.Do(ctx, &fooVal, "GET", "foo")
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

func setSet(label string, client *zgoredis, i int) chan int {
	//还需要一个上下文用来控制开出去的goroutine是否超时
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	//var fooVal string
	result, err := client.Do(ctx, nil, "SET", "foo", "someval")
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
