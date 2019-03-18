package zgomysql

import (
	"context"
	"errors"
	"fmt"
	"git.zhugefang.com/gocore/zgo/config"
	"testing"
	"time"
)

const (
	label_bj = "mysql_sell_1"
	label_sh = "mysql_sell_2"
	l        = 30 //暴力测试50000个消息，时间10秒，本本的并发每秒5000
)

// 实体类
type House struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (u *House) BeforeUpdate() (err error) {
	if u.Id == 0 {
		err = errors.New("user id is 0")
	}
	fmt.Println(u.Id)
	return
}

// 创建链接
func newService() (Mysqler, error) {
	//-------------test for start engine---------
	hsm := make(map[string][]*config.ConnDetail)
	cd_bjw := config.ConnDetail{
		C:           "北京主库-----mysql",
		T:           "w",
		Host:        "localhost",
		Port:        3307,
		Username:    "root",
		Password:    "root",
		DbName:      "spider",
		MaxOpenConn: 5,
		MaxIdleSize: 5,
	}
	cd_bjr := config.ConnDetail{
		C:           "北京主库-----mysql",
		Host:        "localhost",
		Port:        3307,
		Username:    "root",
		Password:    "root",
		DbName:      "spider",
		T:           "r",
		MaxOpenConn: 5,
		MaxIdleSize: 5,
	}
	cd_shw := config.ConnDetail{
		C:           "上海主库-----mysql",
		Host:        "localhost",
		Port:        3308,
		Username:    "root",
		Password:    "root",
		DbName:      "spider_sh",
		T:           "w",
		MaxOpenConn: 5,
		MaxIdleSize: 5,
	}
	cd_shr := config.ConnDetail{
		C:           "上海主库-----mysql",
		Host:        "localhost",
		Port:        3308,
		Username:    "root",
		Password:    "root",
		DbName:      "spider_sh",
		T:           "r",
		MaxOpenConn: 5,
		MaxIdleSize: 5,
	}
	//config.Conf.CityDbConfig := map[string]map[string]string{"sell": {"bj": "1", "sh": "2"}}

	var s1 []*config.ConnDetail
	var s2 []*config.ConnDetail

	s1 = append(s1, &cd_bjw)
	s1 = append(s1, &cd_bjr)
	s2 = append(s2, &cd_shw)
	s2 = append(s2, &cd_shr)

	hsm = map[string][]*config.ConnDetail{
		label_sh: s2,
		label_bj: s1,
	}
	//----------------------
	InitMysql(hsm) //测试时表示使用mysql，在zgo_start中使用一次
	//测试读取nsq数据，wait for sdk init connection
	time.Sleep(2 * time.Second)
	return GetMysql(label_bj)
}

func testMysql(fn func(client Mysqler, i int, city string) chan int) {
	clientBj, _ := newService()

	var replyChan = make(chan int)
	var countChan = make(chan int)

	count := []int{}
	total := []int{}
	totalFinish := []int{}
	stime := time.Now()

	for i := 0; i < l; i++ {
		go func(i int) {
			countChan <- i //统计开出去的goroutine
			//ch := getMongo(label_sh,clientBj,i)
			ch := fn(clientBj, i, "bj")
			reply := <-ch
			replyChan <- reply

			//ch := getMongo(label_bj,clientSh,i)
			ch1 := fn(clientBj, i, "sh")
			reply1 := <-ch1
			replyChan <- reply1
		}(i)
	}

	go func() {
		for v := range replyChan {
			totalFinish = append(totalFinish, v)
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
		if len(totalFinish) == 2*l {
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

// 创建
func TestMysqlCreate(t *testing.T) {
	testMysql(createMysql)
}

// 查询
func TestMysqlGet(t *testing.T) {
	testMysql(getMysql)
}

// 更新
func TestMysqlUpdate(t *testing.T) {
	testMysql(updateMysql)
}

// 列表
func TestMysqlList(t *testing.T) {
	testMysql(listMysql)
}

// count
func TestMysqlCount(t *testing.T) {
	testMysql(countMysql)
}

// 删除
func TestMysqlDelete(t *testing.T) {
	testMysql(deleteMysql)
}

func getMysql(client Mysqler, i int, city string) chan int {
	//还需要一个上下文用来控制开出去的goroutine是否超时
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	client.MysqlServiceByCityBiz(city, "sell")
	dbName, _ := client.GetDbByCityBiz(city, "sell")
	// 开始查询
	house1 := &House{}
	args := make(map[string]interface{})
	args["table"] = dbName + ".house"
	args["query"] = " id = ? "
	args["args"] = []interface{}{i}
	//args["args"] = []interface{}{1}
	args["obj"] = house1
	client.Get(ctx, args)

	out := make(chan int, 1)
	select {
	case <-ctx.Done():
		fmt.Println("超时")
		out <- 10001
		return out
	default:
		//fmt.Println(string(bytes), err, "---from mongo successful---")
		fmt.Println(house1)
		out <- 1
	}
	return out
}

func createMysql(ms Mysqler, i int, city string) chan int {
	//还需要一个上下文用来控制开出去的goroutine是否超时
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// 开始查询
	house1 := &House{Name: city + ":house:" + string(i)}
	args := make(map[string]interface{})
	args["table"] = "house"
	args["obj"] = house1
	label, _ := ms.GetLabelByCityBiz(city, "sell")
	client, err := GetMysql(label)
	client.Create(ctx, args)
	if err != nil {
		fmt.Println(err.Error())
	}
	out := make(chan int, 1)
	select {
	case <-ctx.Done():
		fmt.Println("超时")
		out <- 10001
		return out
	default:
		//fmt.Println(string(bytes), err, "---from mongo successful---")
		fmt.Println(house1)
		out <- 1
	}
	return out
}

func listMysql(ms Mysqler, i int, city string) chan int {
	//还需要一个上下文用来控制开出去的goroutine是否超时
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// 开始查询
	//house1 := &House{Name:city+":house:"+string(i)}
	args := make(map[string]interface{})
	args["table"] = "house"
	args["query"] = " id < ? "
	args["args"] = []interface{}{int(i)}
	args["limit"] = 30
	args["offset"] = 0
	args["order"] = " id desc "
	obj := make([]House, 0)
	args["obj"] = &obj
	label, _ := ms.GetLabelByCityBiz(city, "sell")
	client, err := GetMysql(label)
	client.List(ctx, args)
	if err != nil {
		fmt.Println(err.Error())
	}
	out := make(chan int, 1)
	select {
	case <-ctx.Done():
		fmt.Println("超时")
		out <- 10001
		return out
	default:
		//fmt.Println(string(bytes), err, "---from mongo successful---")
		fmt.Println(obj)
		out <- 1
	}
	return out
}

func updateMysql(ms Mysqler, i int, city string) chan int {
	//还需要一个上下文用来控制开出去的goroutine是否超时
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// 开始查询
	house1 := &House{Id: i}
	args := make(map[string]interface{})
	args["table"] = "house"
	args["obj"] = house1
	label, _ := ms.GetLabelByCityBiz(city, "sell")
	client, err := GetMysql(label)
	data := map[string]interface{}{"name": city + ":uHouse:" + string(i)}
	args["data"] = data
	c, _ := client.UpdateOne(ctx, args)
	fmt.Println("更新了" + string(c) + "条数据")
	if err != nil {
		fmt.Println(err.Error())
	}
	out := make(chan int, 1)
	select {
	case <-ctx.Done():
		fmt.Println("超时")
		out <- 10001
		return out
	default:
		//fmt.Println(string(bytes), err, "---from mongo successful---")
		fmt.Println(house1)
		out <- 1
	}
	return out
}

func countMysql(ms Mysqler, i int, city string) chan int {
	//还需要一个上下文用来控制开出去的goroutine是否超时
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// 开始查询
	count := 0
	args := make(map[string]interface{})
	args["table"] = "house"
	args["count"] = &count
	args["query"] = " id < ? "
	args["args"] = []interface{}{int(i)}
	label, _ := ms.GetLabelByCityBiz(city, "sell")
	client, err := GetMysql(label)
	client.Count(ctx, args)
	if err != nil {
		fmt.Println(err.Error())
	}
	out := make(chan int, 1)
	select {
	case <-ctx.Done():
		fmt.Println("超时")
		out <- 10001
		return out
	default:
		//fmt.Println(string(bytes), err, "---from mongo successful---")
		fmt.Println(count)
		fmt.Println("查询到了" + string(count) + "条数据")
		out <- 1
	}
	return out
}

func deleteMysql(ms Mysqler, i int, city string) chan int {
	//还需要一个上下文用来控制开出去的goroutine是否超时
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// 开始查询
	house1 := &House{Id: i}
	args := make(map[string]interface{})
	args["table"] = "house"
	args["obj"] = house1
	args["id"] = i
	label, _ := ms.GetLabelByCityBiz(city, "sell")
	client, err := GetMysql(label)
	cn, _ := client.DeleteOne(ctx, args)
	if err != nil {
		fmt.Println(err.Error())
	}
	out := make(chan int, 1)
	select {
	case <-ctx.Done():
		fmt.Println("超时")
		out <- 10001
		return out
	default:
		//fmt.Println(string(bytes), err, "---from mongo successful---")
		fmt.Println("删除了" + string(cn) + "条数据")
		out <- 1
	}
	return out
}
