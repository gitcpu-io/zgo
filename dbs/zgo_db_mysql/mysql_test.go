package zgo_db_mysql

import (
	"context"
	"fmt"
	"testing"
	"time"
)

//Define Resourcer struct
type House struct {
	Id   int
	Name string
}

//测试 查询数据库
func TestInit(t *testing.T) {
	replyChan := make(chan interface{})
	start := time.Now().Nanosecond()
	var l = 2
	var count = 0
	for i := 0; i < l; i++ {
		go func(i int) {
			ch := queryMysql(i)

			reply := <-ch
			replyChan <- reply
		}(i)
	}

	for v := range replyChan {
		fmt.Println(v)
		count++
		if count == l {
			break
		}
	}
	fmt.Println("count:", count)
	fmt.Println(time.Now().Nanosecond() - start)
	//db := dbs["sell"]
	//// 创建
	//db.Create(&House{Name: "test"})
	//
	//// 读取
	//var house House
	//db.First(&house, 1) // 查询id为1的product
	//fmt.Println(&house)
	//db.First(&house, "name = ?", "test") // 查询code为l1212的product
	//fmt.Println(&house)
	//// 更新 - 更新product的price为2000
	//db.Model(&house).Update("name", "22222")
	//fmt.Println(&house)
	// 删除 - 删除product
	//db.Delete(&house)

}

func queryMysql(n int) chan interface{} {
	ctx := context.Background()
	var house House
	outch := First(ctx, house)
	fmt.Println(n)
	return outch
}
