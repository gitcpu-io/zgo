package zgo_db_mongo

import (
	"context"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/json-iterator/go"
	"testing"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func TestInit(t *testing.T) {
	//测试读取mongodb数据，wait for sdk init connection
	time.Sleep(2 * time.Second)

	var replyChan = make(chan int)
	var countChan = make(chan int)
	l := 50000 //暴力测试50000个查询，时间10秒，本本的并发每秒5000

	count := []int{}
	total := []int{}
	stime := time.Now()

	for i := 0; i < l; i++ {
		go func(i int) {
			countChan <- i //统计开出去的goroutine
			ch := QueryMongo(i)
			reply := <-ch
			replyChan <- reply
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

			fmt.Printf("总消耗时间：%s, 成功：%d, 总共开出来的goroutine：%d", timeLen, len(count), len(total))
			break
		}

		select {
		case <-time.Tick(time.Duration(1000 * time.Millisecond)):
			fmt.Println("处理进度每1000毫秒", len(count))

		}
	}
	time.Sleep(2 * time.Second)
}

func QueryMongo(i int) chan int {
	//这里需要一个mongoChan
	mongoChan := MongoClientChan()
	//c := <-mongoChan //原生使用client connection

	//还需要一个上下文用来控制开出去的goroutine是否超时
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	//输入参数：上下文ctx，mongoChan里面是client的连接，args具体的查询操作参数

	result, err := dealQuery(ctx, mongoChan, "local", "startup_log", bson.M{})
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

func dealQuery(ctx context.Context, ch chan *mgo.Session, db, collection string, query bson.M) (interface{}, error) {
	s := <-ch
	var res interface{}
	//stime := time.Now()
	s.DB(db).C(collection).Find(&query).One(&res)
	//fmt.Println(time.Now().Sub(stime))
	return res, nil
}
