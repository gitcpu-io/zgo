package zgomysql1

import (
	"context"
	"fmt"
	"git.zhugefang.com/gocore/zgo/config"
	"testing"
	"time"
)

const (
	label_bj = "mysql_label_bj"
	label_sh = "mysql_label_sh"
)

func TestMysqlGet(t *testing.T) {

	//-------------test for start engine---------
	hsm := make(map[string][]*config.ConnDetail)
	cd_bj := config.ConnDetail{
		C:        "北京主库-----mysql",
		Host:     "root:123456@(localhost:3306)/spider?charset=utf8&parseTime=True&loc=Local",
		ConnSize: 1,
		PoolSize: 5000,
	}
	cd_sh := config.ConnDetail{
		C:        "上海主库-----mysql",
		Host:     "root:123456@(localhost:3306)/spider?charset=utf8&parseTime=True&loc=Local",
		ConnSize: 1,
		PoolSize: 5000,
	}
	var s1 []*config.ConnDetail
	var s2 []*config.ConnDetail
	s1 = append(s1, &cd_bj)
	s2 = append(s2, &cd_sh)
	hsm = map[string][]*config.ConnDetail{
		label_bj: s1,
		label_sh: s2,
	}
	//----------------------

	InitMysql(hsm) //测试时表示使用mysql，在zgo_start中使用一次

	//测试读取nsq数据，wait for sdk init connection
	time.Sleep(2 * time.Second)

	clientBj, err := GetMysql(label_bj)
	clientSh, err := GetMysql(label_sh)
	if err != nil {
		panic(err)
	}

	var replyChan = make(chan int)
	var countChan = make(chan int)
	l := 5 //暴力测试50000个消息，时间10秒，本本的并发每秒5000

	count := []int{}
	total := []int{}
	stime := time.Now()

	for i := 0; i < l; i++ {
		go func(i int) {
			countChan <- i //统计开出去的goroutine
			if i%2 == 0 {
				//ch := getMysql(label_sh,clientBj,i)
				ch := getMysqlTx(label_sh, clientBj, i)
				reply := <-ch
				replyChan <- reply

			} else {
				//ch := getMysql(label_bj,clientSh,i)
				ch := getMysql(label_bj, clientSh, i)
				reply := <-ch
				replyChan <- reply
			}
		}(i)
	}

	go func() {
		for v := range replyChan {
			if v != 10001 { //10001表示超时
				count = append(count, v) //成功数
			} else {

				fmt.Println("有不成功的")
			}
		}
	}()

	go func() {
		for v := range countChan { //总共的goroutine
			total = append(total, v)
		}
	}()

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

func getMysqlTx(label string, client *zgomysql, i int) chan int {
	fmt.Println("开始")
	//还需要一个上下文用来控制开出去的goroutine是否超时
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	//输入参数：上下文ctx，mysqlChan里面是client的连接，args具体的查询操作参数
	house := &House{}
	args := make(map[string]interface{})
	args["tablename"] = "house"
	args["query"] = " id = ? "
	args["args"] = []interface{}{1}
	args["out"] = house
	ch, _ := client.GetConnChan(label_bj)
	db := <-ch
	a := db.DB()
	a.Begin()
	fmt.Println("开启事务")
	time.Sleep(10 * time.Second)
	a.Query("select * from house")
	fmt.Println("关闭事务")
	db.DB().Close()
	//err := client.Get(ctx, args)
	//if err != nil {
	//	panic(err)
	//}
	out := make(chan int, 1)
	select {
	case <-ctx.Done():
		fmt.Println("超时")
		out <- 10001
		return out
	default:
		fmt.Println(house)
		//fmt.Println(string(bytes), err, "---from mysql successful---")
		out <- 1
	}

	return out

}

func getMysql(label string, client *zgomysql, i int) chan int {
	fmt.Println("开始")
	//还需要一个上下文用来控制开出去的goroutine是否超时
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	//输入参数：上下文ctx，mysqlChan里面是client的连接，args具体的查询操作参数
	house := &House{}
	args := make(map[string]interface{})
	args["tablename"] = "house"
	args["query"] = " id = ? "
	args["args"] = []interface{}{1}
	args["out"] = house
	err := client.Get(ctx, args)
	//err := client.Get(ctx, args)
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
		fmt.Println(house)
		//fmt.Println(string(bytes), err, "---from mysql successful---")
		out <- 1
	}

	return out

}

type House struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}
