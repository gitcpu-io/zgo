package zgomysql

import (
	"context"
	"fmt"
	"git.zhugefang.com/gocore/zgo.git/config"
	"testing"
	"time"
)

const (
	label_bj = "mysql_sell_r1"
	label_sh = "mysql_sell_r2"
)

// 实体类
type House struct {
	Name string
	Id   int
}

// 创建链接

// 查询
func TestMysqlGet(t *testing.T) {

	//-------------test for start engine---------
	hsm := make(map[string][]*config.ConnDetail)
	cd_bj := config.ConnDetail{
		C:           "北京主库-----mysql",
		Host:        "root:123456@(localhost:3306)/spider?charset=utf8&parseTime=True&loc=Local",
		MaxOpenConn: 5,
		MaxIdleSize: 5,
	}
	cd_sh := config.ConnDetail{
		C:           "上海主库-----mysql",
		Host:        "root:123456@(localhost:3306)/spider_sh?charset=utf8&parseTime=True&loc=Local",
		MaxOpenConn: 5,
		MaxIdleSize: 5,
	}
	cityDbConfig := map[string]map[string]string{"sell": {"bj": "1"}}

	var s1 []*config.ConnDetail
	var s2 []*config.ConnDetail

	s1 = append(s1, &cd_bj)
	s2 = append(s2, &cd_sh)

	hsm = map[string][]*config.ConnDetail{
		label_bj: s1,
		label_sh: s2,
	}
	//----------------------

	InitMysqlService(hsm, cityDbConfig) //测试时表示使用mysql，在zgo_start中使用一次

	//测试读取nsq数据，wait for sdk init connection
	time.Sleep(2 * time.Second)

	clientBj := MysqlService()

	clientSh := MysqlService()
	//if err != nil {
	//	panic(err)
	//}

	var replyChan = make(chan int)
	var countChan = make(chan int)
	l := 5000 //暴力测试50000个消息，时间10秒，本本的并发每秒5000

	count := []int{}
	total := []int{}
	stime := time.Now()

	for i := 0; i < l; i++ {
		go func(i int) {
			countChan <- i //统计开出去的goroutine
			if i%2 == 0 {
				//ch := getMongo(label_sh,clientBj,i)
				ch := getMysql(clientBj, i)
				reply := <-ch
				replyChan <- reply

			} else {
				//ch := getMongo(label_bj,clientSh,i)
				ch := getMysql(clientSh, i)
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

	//for _, v := range count {
	//	if v != 1 {
	//		fmt.Println("有不成功的")
	//	}
	//}

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

func getMysql(client MysqlServiceInterface, i int) chan int {

	//还需要一个上下文用来控制开出去的goroutine是否超时
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// 开始查询
	// 获取对应的label
	label, err := client.GetLabelByCity("bj", "sell", "r")
	if err != nil {
		panic(err)
	}
	// 获取re对象
	re, err := client.NewRs(label)
	if err != nil {
		panic(err)
	}

	// 获取链接池对象
	pool := re.GetPool()
	dbName, _ := client.GetDbByCityBiz("sh", "sell")
	house := &House{}
	pool.Table(dbName+".house").Where(" id = ? ", 1).First(house)
	fmt.Println(0, house)

	//直接用resource查询
	house1 := &House{}
	args := make(map[string]interface{})
	args["tablename"] = "house"
	args["query"] = " id = ? "
	args["args"] = []interface{}{1}
	args["out"] = house1
	re.Get(ctx, args)
	fmt.Println(1, house1)

	out := make(chan int, 1)
	select {
	case <-ctx.Done():
		fmt.Println("超时")
		out <- 10001
		return out
	default:
		//fmt.Println(string(bytes), err, "---from mongo successful---")
		fmt.Println(house)
		out <- 1
	}
	return out
}
