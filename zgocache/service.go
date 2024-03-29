package zgocache

import (
  "context"
  "errors"
  "fmt"
  "github.com/gitcpu-io/zgo/config"
  "github.com/gitcpu-io/zgo/zgoredis"
  "github.com/gitcpu-io/zgo/zgoutils"
  "github.com/json-iterator/go"
  "log"
  "time"
)

func InitCacheByEtcd(v *config.CacheConfig) chan Cacher {
  out := make(chan Cacher)
  go func() { //接收到etcd变化后，触发label和expire的值
    //fmt.Printf("Label:%v; Rate:%v; DbType:%v; TcType:%v; Start:%v; -----etcd tiger cache value----\n", v.Label, v.Rate, v.DbType, v.TcType, v.Start)
    rate := v.Rate
    dbtype := v.DbType
    tcType := v.TcType
    label := v.Label
    start := v.Start
    out <- GetCache(start, dbtype, label, rate, tcType)
  }()
  return out
}

func InitCache() chan Cacher {
  out := make(chan Cacher)
  go func() {
    hm := config.Conf.Cache
    //fmt.Printf("Label:%v; Rate:%v; DbType:%v; TcType:%v; Start:%v; -----etcd tiger cache value----\n", hm.Label, hm.Rate, hm.DbType, hm.TcType, hm.Start)
    rate := hm.Rate
    dbtype := hm.DbType
    tcType := hm.TcType
    label := hm.Label
    start := hm.Start
    out <- GetCache(start, dbtype, label, rate, tcType)
  }()
  return out
}

type dbServicer interface {
  Hget(ctx context.Context, key string, field string) (interface{}, error)
  Hset(ctx context.Context, key string, field string, value interface{}) (int, error)
  Hdel(ctx context.Context, key string, field string) (int, error)
}

/*
 GetCache 创建service对象的方法
 1.start 是否开启
 2.dbtype 缓存底层数据库类型
 3.label 标签
 4.expire 是否开启
 5.tcType 超时对象
*/
func GetCache(start int, dbtype string, label string, rate int, tcType int) Cacher {
  if start == 1 {
    if dbtype == "pika" {
      // todo 找不到pika后异常处理
      var service dbServicer = zgoredis.Redis(label)
      return &zgocache{
        start,
        label,
        dbtype,
        service,
        tcType,
        rate,
      }
    } else if dbtype == "redis" {
      var service dbServicer = zgoredis.Redis(label)
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
    //fmt.Println("未配置缓存")
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
  //NewPikaCacheService(label string, expire int, tcType int) Cacher
  Decorate(fn CacheFunc, expire int) CacheFunc
  DelCache(ctx context.Context, key, field string) (int, error)
  TimeOutDecorate(fn CacheFunc, timeout int) CacheFunc
}

// 缓存装饰器接收的函数类型
type CacheFunc func(ctx context.Context, param map[string]interface{}, obj interface{}) error

// 函数返回值 channel接收时候用到
//type funResult struct {
//	err    error
//}

// 缓存入redis结构体
type cacheResult struct {
  Result interface{}
  Time   int64
}

//zgocache 结构体
type zgocache struct {
  start   int
  label   string
  dbtype  string
  service dbServicer
  tcType  int // 1 降级缓存 2 正常缓存
  rate    int // 失效时间 倍率
}

// Decorate 缓存装饰器
// 1.fn 真正执行的方法，必须符合CacheFunc类型
// 2.expire 超时时间 单位s
func (z *zgocache) Decorate(fn CacheFunc, expire int) CacheFunc {
  return func(ctx context.Context, param map[string]interface{}, obj interface{}) error {
    cacheModel := ""
    if v, ok := param["cacheModel"]; ok && v.(string) != "" {
      cacheModel = v.(string)
    }
    // start 是否启动 1 启动，0 停用
    if z.start != 1 || cacheModel == "" {
      return fn(ctx, param, obj)
    }
    key := z.getKey(cacheModel)
    field := ""
    var err error
    if v, ok := param["cacheField"]; ok && v.(string) != "" {
      field = v.(string)
    } else {
      field, err = zgoutils.Utils.MarshalMap(param)
      if err != nil {
        // field转换失败 直接走函数获取数据
        fmt.Println(err.Error())
        return fn(ctx, param, obj)
      }

    }
    // 获取缓存
    err = z.getData(ctx, key, field, expire, obj)
    if err != nil { // 有异常 或者 没有缓存
      // 执行函数获取数据
      err = fn(ctx, param, obj)
      if obj != nil && err == nil {
        // 正常返回结果 存入缓存
        z.setData(ctx, key, field, obj)
      }
      // 返回结果
      return err
    }
    return err
  }
}

func (z *zgocache) DelCache(ctx context.Context, cacheModel, cacheField string) (int, error) {
  if z.start != 1 {
    return 0, nil
  }
  key := z.getKey(cacheModel)
  return z.service.Hdel(ctx, key, cacheField)
}

// TimeOutDecorate 降级缓存装饰器
func (z *zgocache) TimeOutDecorate(fn CacheFunc, timeout int) CacheFunc {
  return func(ctx context.Context, param map[string]interface{}, obj interface{}) error {
    cacheModel := ""
    if v, ok := param["cacheModel"]; ok {
      cacheModel = v.(string)
    }
    // start 是否启动 1 启动，0 停用
    if z.start != 1 && cacheModel == "" {
      return fn(ctx, param, obj)
    }
    // 当调用方法后续服务异常时，通过etcd配置修改tcType为2。可转为走正常缓存逻辑，并且没有失效时间。
    if z.tcType == 2 {
      return z.Decorate(fn, 0)(ctx, param, obj)
    }
    ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
    defer cancel()
    ch := make(chan error)

    // 执行
    go func(ctx context.Context) {
      err := fn(ctx, param, obj)
      ch <- err
    }(ctxTimeout)

    // 缓存结果
    field, fieldErr := zgoutils.Utils.MarshalMap(param)
    key := z.getKey(cacheModel)

    select {
    case <-ctxTimeout.Done():
      // 拼接key 获取缓存返回
      // 失败返回
      if fieldErr != nil {
        fmt.Println(fieldErr.Error())
        return fieldErr
      }
      // 返回
      err := z.getData(ctx, key, field, 0, obj)
      return err

    case value, ok := <-ch:
      if ok {
        // 查询成功返回数据 并 塞入缓存
        if value != nil {
          return value
        }
        z.setData(ctx, key, field, obj)
        return nil
      }

    }

    return errors.New("操作失败")
  }
}

// NewPikaCacheService 创建新的缓存
//func (z *zgocache) NewPikaCacheService(label string, expire int, tcType int) Cacher {
//	return GetCache(z.start, "pika", label, expire, tcType)
//}

func (z *zgocache) getData(ctx context.Context, key string, field string, expire int, obj interface{}) error {
  // 根据项目名，包名，类名，函数名称，拼接，然后数据库
  //project := config.Project
  //fn := runtime.FuncForPC(reflect.ValueOf(a).Pointer()).Name()
  //path := reflect.TypeOf(a).PkgPath()
  value, err := z.service.Hget(ctx, key, field)
  if err != nil {
    fmt.Println(err.Error())
    return err
  } else if value == nil || value == "" {
    return errors.New("缓存数据为空")
  } else {
    data := cacheResult{Result: obj}
    err := jsoniter.UnmarshalFromString(value.(string), &data)
    if err != nil {
      fmt.Println(err.Error())
      return err
    }
    if expire != 0 {
      if data.Time < time.Now().Unix()-int64(expire)*int64(z.rate) {
        return errors.New("缓存已失效")
      }
    }
    return nil
  }
}

func (z *zgocache) setData(ctx context.Context, key string, field string, data interface{}) {
  // 开goroutine 存缓存
  d := &cacheResult{data, time.Now().Unix()}
  value, err := jsoniter.MarshalToString(d)

  go func(ctx context.Context) {
    defer func() {
      err := recover()
      if err != nil {
        fmt.Println("添加缓存 gorouitne panic Error")
      }
    }()
    if err != nil {
      fmt.Println(err.Error())
    } else {
      _, err := z.service.Hset(ctx, key, field, value)
      if err != nil {
        fmt.Println(err)
      }
    }
  }(ctx)
}

func (z *zgocache) getKey(model string) string {
  key := ":cache:" + config.Conf.Project + ":" + model //+ ":" + runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
  return key
}
