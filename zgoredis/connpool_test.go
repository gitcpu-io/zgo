package zgoredis

import (
  "context"
  "fmt"
  "github.com/gitcpu-io/zgo/config"
  "github.com/json-iterator/go"
  "github.com/mediocregopher/radix/v4"
  "testing"
  "time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
  label_bj = "redis_label_bj"
  label_sh = "redis_label_sh"
)

func TestRedisGet(t *testing.T) {
  cd_bj := config.ConnDetail{
    C:        "北京主库-----redis1",
    Host:     "localhost",
    Port:     6380,
    ConnSize: 10,
    PoolSize: 200,
    Username: "",
    Password: "",
    Db:       0,
  }
  cd_bj2 := config.ConnDetail{
    C:        "北京从库-----redis2",
    Host:     "localhost",
    Port:     6380,
    ConnSize: 10,
    PoolSize: 200,
    Username: "",
    Password: "",
    Db:       0,
  }
  cd_sh := config.ConnDetail{
    C:        "上海主库-----redis",
    Host:     "localhost",
    Port:     6380,
    ConnSize: 10,
    PoolSize: 200,
    Username: "",
    Password: "",
    Db:       5,
  }
  var s1 []*config.ConnDetail
  var s2 []*config.ConnDetail
  s1 = append(s1, &cd_bj, &cd_bj2)
  s2 = append(s2, &cd_sh)
  hsm := map[string][]*config.ConnDetail{
    label_bj: s1,
    label_sh: s2,
  }

  InitRedis(hsm) //测试时表示使用redis，在origin中使用一次

  clientLocal, err := GetRedis(label_bj)
  if err != nil {
    panic(err)
  }
  clientSpider, err := GetRedis(label_sh)

  fmt.Println(clientLocal)
  fmt.Println(clientSpider)

  if err != nil {
    panic(err)
  }

  //ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
  //defer cancel()

  //top
  //ZRANK zadd zrange zincrby zrem zreverange ZSCORE
  //res, err := clientLocal.Zrank(ctx, "page_rank", "bing.com")
  //fmt.Println(res)

  //res, err := clientLocal.Zscore(ctx, "page_rank", "bing.com")
  //fmt.Println(res)

  //res, err := clientLocal.Zrange(ctx, "page_rank", 0, -1, false)
  //fmt.Println(res)

  //res, err := clientLocal.Zrange(ctx, "page_rank", 0, -1, false)
  //fmt.Println(res)

  //res, err := clientLocal.Zrevrange(ctx, "page_rank", 0, -1, false)
  //fmt.Println(res)

  //res, err := clientLocal.ZINCRBY(ctx, "page_rank", 20, "bing.com")
  //fmt.Println(res)

  //res, err := clientLocal.Zadd(ctx, "page_rank", 20, "ssbing.com")
  //fmt.Println(res)

  //res, err := clientLocal.Zrem(ctx, "page_rank", "ooooewew", "ewewew")
  //fmt.Println(res)

  //1.set get
  //res, err := clientLocal.Set(ctx, "china666", 999)
  //
  //res2, err := clientLocal.Get(ctx, "china666")
  //
  //fmt.Println(res)
  //
  //fmt.Println(res2)
  //2.expire
  //res, err := clientLocal.Expire(ctx, "china666", 20000)
  //fmt.Println(res)

  //3.hset hget hlen hdel
  //res, err := clientLocal.Hset(ctx, "china_lining", "liuwei", 9999)
  //fmt.Println(res)
  //
  //res1, err := clientLocal.Hget(ctx, "china_lining", "liuwei")
  //fmt.Println(res1)

  //res1, err := clientLocal.Hlen(ctx, "china_lining")
  //fmt.Println(res1)

  //res1, err := clientLocal.Hdel(ctx, "china_lining", "china")
  //fmt.Println(res1)

  //res1, err := clientLocal.Hgetall(ctx, "china_lining")
  //fmt.Println(res1)

  //res1, err := clientLocal.Del(ctx, "china_lining")
  //fmt.Println(res1)

  //4.lpush rpush llen lrange lpop rpop

  //res, err := clientLocal.Rpush(ctx, "china_list", 23232343)
  //fmt.Println(res)
  //
  //res2, err := clientLocal.Llen(ctx, "china_list")
  //fmt.Println(res2)

  //res, err := clientLocal.Lrange(ctx, "china_list", 0, 10)
  //fmt.Println(res)

  //res2, err := clientLocal.Rpop(ctx, "china_list")
  //fmt.Println(res2)

  //5.sadd Scard Smembers Sismember
  //res, err := clientLocal.Sadd(ctx, "china_member", 1113402)
  //fmt.Println(res)

  //res, err := clientLocal.Scard(ctx, "china_member")
  //	//fmt.Println(res)

  //res, err := clientLocal.Srem(ctx, "china_member", 890)
  //fmt.Println(res)

  //res, err := clientLocal.Smembers(ctx, "china_member")
  //fmt.Println(res)

  //res, err := clientLocal.Sismember(ctx, "china_member", 456)
  //fmt.Println(res)

  //res, err := clientLocal.Exists(ctx, "china_member")
  //fmt.Println(res)

  //res, err := clientLocal.Keys(ctx, "*")
  //fmt.Println(res)

  //res, err := clientLocal.Ttl(ctx, "china_member")
  //fmt.Println(res)

  //res, err := clientLocal.Type(ctx, "china_member")
  //fmt.Println(res)

  //LpushCheck(label_bj, clientLocal, 0)
  //
  getSet(label_sh, clientSpider, 0)
  hetSet(label_sh, clientSpider, 0)

  var replyChan = make(chan int)
  var countChan = make(chan int)
  l := 1000 //暴力测试50000个消息，时间10秒，本本的并发每秒5000

  count := []int{}
  total := []int{}
  stime := time.Now()

  for i := 0; i < l; i++ {
    go func(i int) {
      countChan <- i //统计开出去的goroutine
      if i%2 == 0 {
        //ch := getSet(label_bj, clientLocal, i)
        ch := setSet(label_bj, clientLocal, i)
        //ch := hetSet(label_bj, clientLocal, i)
        reply := <-ch
        replyChan <- reply

      } else {
        //ch := getSet(label_sh, clientSpider, i)
        //ch := setSet(label_sh, clientSpider, i)
        ch := hetSet(label_sh, clientLocal, i)
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
      var timeLen = time.Since(stime)

      fmt.Printf("总消耗时间：%s, 成功：%d, 总共开出来的goroutine：%d\n", timeLen, len(count), len(total))
      break
    }

    select {
    case <-time.Tick(time.Duration(1000 * time.Millisecond)):
      fmt.Println("处理进度每1000毫秒", len(count))
    default:

    }
  }

}

func hetSet(label string, client *zgoredis, i int) chan int {
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

func setSet(label string, client *zgoredis, i int) chan int {
  //还需要一个上下文用来控制开出去的goroutine是否超时
  ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
  defer cancel()
  key := fmt.Sprintf("foo_%d", i)

  value := "wwwwwwwwwwwwwww"
  for i := 0; i < 3; i++ {
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

func getSet(label string, client *zgoredis, i int) chan int {
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

//func hgetSet(label string, client *zgoredis, i int) chan int {
//  //还需要一个上下文用来控制开出去的goroutine是否超时
//  ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
//  defer cancel()
//
//  key := "foo_china"
//
//  name := fmt.Sprintf("foo_%d", i)
//
//  result, err := client.Hget(ctx, key, name)
//  if err != nil {
//    panic(err)
//  }
//  out := make(chan int, 1)
//  select {
//  case <-ctx.Done():
//    fmt.Println("超时")
//    out <- 10001
//    return out
//  default:
//    fmt.Println(result)
//    out <- 1
//  }
//
//  return out
//}

func LpushCheck(label string, client *zgoredis, i int) chan int {
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

func TestConnPool_GetConnChan(t *testing.T) {
  c, err := (radix.PoolConfig{
    Dialer: radix.Dialer{
    },
  }).New(context.TODO(),"tcp","127.0.0.1:6380")
  if err != nil {
    fmt.Println("redis ", err)
  }
  fmt.Println(c)
}
