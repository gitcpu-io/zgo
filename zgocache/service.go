package zgocache

import (
	"context"
	"errors"
	"fmt"
	"git.zhugefang.com/gocore/zgo/config"
	"git.zhugefang.com/gocore/zgo/zgopika"
	"github.com/json-iterator/go"
	"log"
	"reflect"
	"runtime"
	"time"
)

func InitCache(cacheCh chan *config.CacheConfig) chan Cacher {
	out := make(chan Cacher)
	go func() { //接收到etcd变化后，触发label和expire的值
		for v := range cacheCh {
			fmt.Printf("Label:%v; Rate:%v; DbType:%v; TcType:%v; Start:%v; -----etcd tiger cache value----\n", v.Label, v.Rate, v.DbType, v.TcType, v.Start)
			hm := config.Conf.Cache
			rate := hm.Rate
			dbtype := hm.DbType
			tcType := hm.TcType
			label := hm.Label
			start := hm.Start
			out <- GetCache(start, dbtype, label, rate, tcType)
		}
	}()

	go func() { //接收到etcd变化后，触发label和expire的值
		hm := config.Conf.Cache
		fmt.Printf("Label:%v; Rate:%v; DbType:%v; TcType:%v; Start:%v; -----etcd tiger cache value----\n", hm.Label, hm.Rate, hm.DbType, hm.TcType, hm.Start)
		rate := hm.Rate
		dbtype := hm.DbType
		tcType := hm.TcType
		label := hm.Label
		start := hm.Start
		out <- GetCache(start, dbtype, label, rate, tcType)
	}()
	return out
}

/*
 GetCache 创建service对象的方法
 1.start 是否开启
 2.dbtype 是否开启
 3.label 是否开启
 4.expire 是否开启
 5.tcType 超时对象
*/
func GetCache(start int, dbtype string, label string, rate int, tcType int) Cacher {
	if start == 1 {
		if dbtype == "pika" {
			// todo 找不到pika后异常处理
			service := zgopika.Pika(label)
			return &zgocache{
				start,
				label,
				dbtype,
				service,
				tcType,
				rate,
			}
		}
	} else {
		fmt.Println("未配置缓存")
		return &zgocache{
			0,
			label,
			dbtype,
			nil,
			tcType,
			rate,
		}
	}
	log.Fatalf("缓存数据库类型不支持")
	return nil
}

// 对外接口
type Cacher interface {
	NewPikaCacheService(label string, expire int, tcType int) Cacher
	Decorate(fn CacheFunc, expire int) CacheFunc
	TimeOutDecorate(fn CacheFunc, timeout int) CacheFunc
}

// 缓存装饰器接收的函数类型
type CacheFunc func(ctx context.Context, param map[string]interface{}) (interface{}, error)

// 函数返回值 channel接收时候用到
type funResult struct {
	result interface{}
	err    error
}

// 缓存入redis结构体
type cacheResult struct {
	Result interface{}
	Time   int
}

//zgocache 结构体
type zgocache struct {
	start   int
	label   string
	dbtype  string
	service zgopika.Pikaer
	tcType  int // 1 降级缓存 2 正常缓存
	rate    int // 失效时间 倍率
}

// 缓存装饰器
func (z *zgocache) Decorate(fn CacheFunc, expire int) CacheFunc {
	return func(ctx context.Context, param map[string]interface{}) (interface{}, error) {
		fmt.Println("进入Decorate")
		if z.start != 1 {
			return fn(ctx, param)
		}
		key := z.getKey(fn)
		field, err := jsoniter.MarshalToString(param)
		if err != nil {
			// field转换失败 直接走函数获取数据
			fmt.Println(err.Error())
			return fn(ctx, param)
		}

		// 获取缓存
		data, err := z.getData(ctx, key, field, expire)
		if err != nil { // 有异常 或者 没有缓存
			// 执行函数获取数据
			data, err = fn(ctx, param)
			if data != nil && err == nil {
				// 正常返回结果 存入缓存
				z.setData(ctx, key, field, data)
			}
			// 返回结果
			return data, err
		}
		return data, err
	}
}

// 降级缓存装饰器
func (z *zgocache) TimeOutDecorate(fn CacheFunc, timeout int) CacheFunc {
	return func(ctx context.Context, param map[string]interface{}) (interface{}, error) {
		if z.tcType == 2 {
			return z.Decorate(fn, 0)(ctx, param)
		}
		fmt.Println("进入TimeOutDecorate")
		if z.start != 1 {
			return fn(ctx, param)
		}
		ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
		ch := make(chan *funResult)

		// 执行
		go func(ctx context.Context) {
			result, err := fn(ctx, param)
			ch <- &funResult{
				result,
				err,
			}
		}(ctxTimeout)

		// 缓存结果
		field, fieldErr := jsoniter.MarshalToString(param)
		key := z.getKey(fn)

		select {
		case <-ctxTimeout.Done():
			cancel()
			fmt.Println("超时获取缓存")
			// 拼接key 获取缓存返回
			// 失败返回
			if fieldErr != nil {
				fmt.Println(fieldErr.Error())
				return nil, fieldErr
			}
			// 返回
			data, err := z.getData(ctx, key, field, 0)
			return data, err

		case value, ok := <-ch:
			if ok {
				fmt.Println("获取成功")
				// 查询成功返回数据 并 塞入缓存
				z.setData(ctx, key, field, value)
				return value, nil
			}

		}

		return nil, errors.New("操作失败")
	}
}

// 创建新的缓存
func (z *zgocache) NewPikaCacheService(label string, expire int, tcType int) Cacher {
	return GetCache(z.start, "pika", label, expire, tcType)
}

func (z *zgocache) getData(ctx context.Context, key string, field string, expire int) (interface{}, error) {
	// 根据项目名，包名，类名，函数名称，拼接，然后数据库
	//project := config.Project
	//fn := runtime.FuncForPC(reflect.ValueOf(a).Pointer()).Name()
	//path := reflect.TypeOf(a).PkgPath()

	fmt.Println("取", key, ":", field)
	value, err := z.service.Hget(ctx, key, field)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	} else if value == nil || value == "" {
		fmt.Println("没有缓存")
		return nil, errors.New("缓存数据为空")
	} else {
		data := &cacheResult{}
		jsoniter.UnmarshalFromString(value.(string), data)
		if expire != 0 {
			if data.Time < time.Now().Second()-expire*z.rate {
				return nil, errors.New("缓存已失效")
			}
		}
		fmt.Println("有缓存")
		return data.Result, nil
	}
}

func (z *zgocache) setData(ctx context.Context, key string, field string, data interface{}) {
	// 开goroutine 存缓存
	go func(ctx context.Context) {
		d := &cacheResult{data, time.Now().Second()}
		value, err := jsoniter.MarshalToString(d)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("缓存放入失败")
		} else {

			fmt.Println("存：", key, ":", field)
			_, err := z.service.Hset(ctx, key, field, value)
			fmt.Println(err.Error())
		}
	}(ctx)

}

func (z *zgocache) getKey(fn CacheFunc) string {
	key := "GOCache:" + config.Conf.Project + ":" + runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	return key
}
