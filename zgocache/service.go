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

var (
	expire int
	label  string
)

func InitCache(cacheCh chan *config.CacheConfig) CacheServiceInterface {
	go func() { //接收到etcd变化后，触发label和expire的值

		for v := range cacheCh {

			label = v.Label
			expire = v.Expire

			fmt.Println(label, expire, "-----etcd tiger cache value----")

		}

	}()

	hm := config.Cache
	expire := 86400
	label := ""
	dbtype := "pika"
	if hm.Start != 1 {
		return GetCache(0, dbtype, label, expire)
	}
	if hm.Expire != 0 {
		expire = hm.Expire
	}
	if hm.Label != "" {
		label = hm.Label
	}
	if hm.DbType != "" {
		dbtype = hm.DbType
	}
	return GetCache(1, dbtype, label, expire)
}

// 创建service对象的方法
func GetCache(start int, dbtype string, label string, expire int) CacheServiceInterface {
	if start == 1 {
		if dbtype == "pika" {
			// todo 找不到pika后异常处理
			service := zgopika.Pika(label)
			return &zgocache{
				label,
				dbtype,
				service,
				expire,
				start,
			}
		}
	} else {
		fmt.Println("未配置缓存")
		return &zgocache{
			"",
			"",
			nil,
			0,
			start,
		}
	}
	log.Fatalf("缓存数据库类型不支持")
	return nil
}

// 对外接口
type CacheServiceInterface interface {
	NewPikaCacheService(label string, expire int) CacheServiceInterface
	Decorate(fn CacheFunc) CacheFunc
	TimeOutDecorate(fn CacheFunc) CacheFunc
}

// 缓存装饰器接收的函数类型
type CacheFunc func(ctx context.Context, param map[interface{}]interface{}) (interface{}, error)

// 函数返回值 channel接收时候用到
type funResult struct {
	result interface{}
	err    error
}

// 缓存入redis结构体
type cacheResult struct {
	result interface{}
	time   int
}

//zgocache 结构体
type zgocache struct {
	label   string
	dbtype  string
	service zgopika.Pikaer
	expire  int
	start   int
}

// 缓存装饰器
func (z *zgocache) Decorate(fn CacheFunc) CacheFunc {
	return func(ctx context.Context, param map[interface{}]interface{}) (interface{}, error) {

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
		data, err := z.getData(ctx, key, field)
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
func (z *zgocache) TimeOutDecorate(fn CacheFunc) CacheFunc {
	return func(ctx context.Context, param map[interface{}]interface{}) (interface{}, error) {
		fmt.Println("进入TimeOutDecorate")
		if z.start != 1 {
			return fn(ctx, param)
		}
		ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
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
			data, err := z.getData(ctx, key, field)
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
func (z *zgocache) NewPikaCacheService(label string, expire int) CacheServiceInterface {
	return GetCache(z.start, "pika", label, expire)
}

func (z *zgocache) getData(ctx context.Context, key string, field string) (interface{}, error) {
	// 根据项目名，包名，类名，函数名称，拼接，然后数据库
	//project := config.Project
	//fn := runtime.FuncForPC(reflect.ValueOf(a).Pointer()).Name()
	//path := reflect.TypeOf(a).PkgPath()
	value, err := z.service.Hget(ctx, key, field)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	} else if value == nil {
		return nil, errors.New("缓存数据为空")
	} else {
		data := &cacheResult{}
		jsoniter.UnmarshalFromString(value.(string), data)
		if data.time < time.Now().Second()-z.expire {
			return nil, errors.New("缓存已失效")
		}
		return data.result, nil
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
			z.service.Hset(ctx, key, field, value)
		}
	}(ctx)

}

func (z *zgocache) getKey(fn CacheFunc) string {
	key := "GOCache_" + config.Project + "_" + runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	return key
}
