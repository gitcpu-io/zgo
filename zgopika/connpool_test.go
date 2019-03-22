package zgopika

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
	bj_rw = "pika_label_rw"
	bj_r  = "pika_label_r"
)

func TestPikaGet(t *testing.T) {
	hsm := make(map[string][]*config.ConnDetail)
	//cd_bj := config.ConnDetail{
	//	C:        "北京主库-----redis1",
	//	Host:     "localhost",
	//	Port:     6379,
	//	ConnSize: 100,go
	//	PoolSize: 20000,
	//	Username: "",
	//	Password: "",
	//	Db:       5,
	//}
	cd_bj_rw := config.ConnDetail{
		C:        "北京主库-----pika1",
		Host:     "101.201.110.130",
		Port:     49221,
		ConnSize: 400,
		PoolSize: 200,
		Username: "",
		Password: "",
		Prefix:   "rent:",
	}

	cd_bj_r := config.ConnDetail{
		C:        "北京从库-----pika2",
		Host:     "101.201.110.130",
		Port:     49221,
		ConnSize: 400,
		PoolSize: 200,
		Username: "",
		Password: "",
		Prefix:   "rent:",
	}

	//cd_sh := config.ConnDetail{
	//	C:        "上海主库-----pika",
	//	Host:     "localhost",
	//	Port:     6379,
	//	ConnSize: 10,
	//	PoolSize: 20000,
	//	Username: "",
	//	Password: "",
	//	Db:       5,
	//}
	var s1 []*config.ConnDetail
	var s2 []*config.ConnDetail
	s1 = append(s1, &cd_bj_rw)
	s2 = append(s2, &cd_bj_r)
	hsm = map[string][]*config.ConnDetail{
		bj_rw: s1,
		bj_r:  s2,
	}

	InitPika(hsm) //测试时表示使用redis，在zgo_start中使用一次

	//测试读取nsq数据，wait for sdk init connection
	//time.Sleep(2 * time.Second)

	clientBjRw, err := GetPika(bj_rw)
	clientBjR, err := GetPika(bj_r)

	if err != nil {
		panic(err)
	}

	//ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	//defer cancel()

	//1.set get
	//res, err := clientBjRw.Set(ctx, "china666", 111111)
	//
	//res2, err := clientBjRw.Get(ctx, "china666")
	//
	//fmt.Println(res)
	//
	//fmt.Println(res2)
	//2.expire
	//res, err := clientBjRw.Expire(ctx, "china666", 20000)
	//fmt.Println(res)

	//3.hset hget hlen hdel
	//res, err := clientBjRw.Hset(ctx, "china_lining", "liuwei", 9999)
	//fmt.Println(res)
	//
	//res1, err := clientBjRw.Hget(ctx, "china_lining", "liuwei")
	//fmt.Println(res1)

	//res1, err := clientBjRw.Hlen(ctx, "china_lining")
	//fmt.Println(res1)

	//res1, err := clientBjRw.Hdel(ctx, "china_lining", "liuwei")
	//fmt.Println(res1)

	//res1, err := clientBjRw.Hgetall(ctx, "china_lining")
	//fmt.Println(res1)

	//res1, err := clientBjRw.Del(ctx, "china_lining")
	//fmt.Println(res1)

	//4.lpush rpush llen lrange lpop rpop

	//res, err := clientBjRw.Lpush(ctx, "china_list", 23232343)
	//fmt.Println(res)
	//
	//res2, err := clientBjRw.Llen(ctx, "china_list")
	//fmt.Println(res2)

	//res, err := clientBjRw.Lrange(ctx, "china_list", 0, 10)
	//fmt.Println(res)

	//res2, err := clientBjRw.Rpop(ctx, "china_list")
	//fmt.Println(res2)

	//5.sadd Scard Smembers Sismember
	//res, err := clientBjRw.Sadd(ctx, "china_member", 1143402)
	//fmt.Println(res)

	//res, err := clientBjRw.Scard(ctx, "china_member")
	//fmt.Println(res)

	//res, err := clientBjRw.Srem(ctx, "china_member", 1143402)
	//fmt.Println(res)

	//res, err := clientBjRw.Smembers(ctx, "china_member")
	//fmt.Println(res)

	//res, err := clientBjRw.Sismember(ctx, "china_member", 1113402)
	//fmt.Println(res)

	//res, err := clientBjRw.Exists(ctx, "china_member")
	//fmt.Println(res)

	//res, err := clientBjRw.Keys(ctx, "*")
	//fmt.Println(res)

	//res, err := clientBjRw.Ttl(ctx, "china_member")
	//fmt.Println(res)

	//res, err := clientBjRw.Type(ctx, "china_member")
	//fmt.Println(res)

	var replyChan = make(chan int)
	var countChan = make(chan int)
	l := 10000 //暴力测试50000个消息，时间10秒，本本的并发每秒5000

	count := []int{}
	total := []int{}
	stime := time.Now()

	for i := 0; i < l; i++ {
		go func(i int) {
			countChan <- i //统计开出去的goroutine
			if i%2 == 0 {
				//ch := getMongo(label_sh,clientBj,i)
				//ch := LpushCheck(bj_rw, clientBjRw, 0)
				//ch := setSet(bj_rw, clientBjRw, i)
				ch := hetSet(bj_rw, clientBjRw, i)
				//ch := hgetSet(bj_rw, clientBjRw, i)
				//ch := getSet(bj_rw, clientBjRw, i)
				reply := <-ch
				replyChan <- reply

			} else {
				//ch := getMongo(label_bj,clientSh,i)
				//ch := getSet(bj_r, clientBjR, 0)
				//ch := setSet(bj_r, clientBjR, i)
				//ch := hgetSet(bj_r, clientBjR, i)
				//ch := getSet(bj_r, clientBjR, i)
				ch := hetSet(bj_r, clientBjR, i)
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

func hetSet(label string, client *zgopika, i int) chan int {
	//还需要一个上下文用来控制开出去的goroutine是否超时
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	key := fmt.Sprintf("foo_china_%d", i)

	name := "foo"

	value := "wwwwwwwwwwwwwww"
	for i := 0; i < 7; i++ {
		value = value + value
	}

	_, err := client.Hset(ctx, key, name, value)
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
		//fmt.Println("ok........")
		out <- 1
	}

	return out
}

func setSet(label string, client *zgopika, i int) chan int {
	//还需要一个上下文用来控制开出去的goroutine是否超时
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	key := fmt.Sprintf("foo_%d", i)

	value := "wwwwwwwwwwwwwww"
	for i := 0; i < 7; i++ {
		value = value + value
	}

	_, err := client.Set(ctx, key, value)
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
		//fmt.Println("ok........")
		out <- 1
	}

	return out
}

func getSet(label string, client *zgopika, i int) chan int {
	//还需要一个上下文用来控制开出去的goroutine是否超时
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	key := fmt.Sprintf("foo_%d", i)
	result, err := client.Get(ctx, key)
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
		fmt.Println(result)
		out <- 1
	}

	return out
}

func hgetSet(label string, client *zgopika, i int) chan int {
	//还需要一个上下文用来控制开出去的goroutine是否超时
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	key := "foo_china"

	name := fmt.Sprintf("foo_%d", i)

	result, err := client.Hget(ctx, key, name)
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
		fmt.Println(result)
		out <- 1
	}

	return out
}

func LpushCheck(label string, client *zgopika, i int) chan int {
	//还需要一个上下文用来控制开出去的goroutine是否超时
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	result, err := client.Lpush(ctx, "china_online_list", "somebody111")
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
