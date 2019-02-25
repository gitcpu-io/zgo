package zgomongo

import (
	"context"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/json-iterator/go"
	"testing"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	label_bj = "label_bj"
	label_sh = "label_sh"
)

func TestMongoGet(t *testing.T) {

	InitMongo(map[string][]string{
		label_bj: []string{
			"localhost:27017",
		},
		label_sh: []string{
			"localhost:27017",
		},
	}) //测试时表示使用nsq，在zgo_start中使用一次

	clientBj, err := GetMongo(label_bj)
	clientSh, err := GetMongo(label_sh)
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
				//ch := getMongo(label_sh,clientBj,i)
				ch := createMongo(label_sh, clientBj, i)
				reply := <-ch
				replyChan <- reply

			} else {
				//ch := getMongo(label_bj,clientSh,i)
				ch := createMongo(label_bj, clientSh, i)
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
	time.Sleep(2 * time.Second)
}

func getMongo(label string, client *zgomongo, i int) chan int {

	//还需要一个上下文用来控制开出去的goroutine是否超时
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	//输入参数：上下文ctx，mongoChan里面是client的连接，args具体的查询操作参数
	args := make(map[string]interface{})
	args["db"] = "local"
	args["collection"] = "startup_log"
	args["query"] = bson.M{}

	result, err := client.Get(ctx, args)
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
		//fmt.Println(string(bytes), err, "---from mongo successful---")
		out <- 1
	}

	return out

}

type user struct {
	Label string `json:"label"`
	Age   int    `json:"age"`
}

func createMongo(label string, client *zgomongo, i int) chan int {

	//还需要一个上下文用来控制开出去的goroutine是否超时
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	//输入参数：上下文ctx，mongoChan里面是client的连接，args具体的查询操作参数
	args := make(map[string]interface{})
	args["db"] = "test"
	args["collection"] = label
	args["items"] = &user{
		Label: label,
		Age:   i,
	}
	result, err := client.Create(ctx, args)
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
		//fmt.Println(string(bytes), err, "---from mongo successful---")
		out <- 1
	}

	return out

}
