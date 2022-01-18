package zgonsq

import (
  "context"
  "fmt"
  "github.com/gitcpu-io/zgo/config"
  "testing"
  "time"
)

const (
  label_bj = "nsq_label_bj"
  label_sh = "nsq_label_sh"
)

func TestProducer(t *testing.T) {
  hsm := make(map[string][]*config.ConnDetail)
  cd_bj := config.ConnDetail{
    C:        "北京主库-----nsq",
    Host:     "localhost",
    Port:     4150,
    ConnSize: 5,
    PoolSize: 246,
  }
  cd_bj2 := config.ConnDetail{
    C:        "北京从库2-----nsq",
    Host:     "localhost",
    Port:     4150,
    ConnSize: 5,
    PoolSize: 135,
  }
  cd_sh := config.ConnDetail{
    C:        "上海主库-----nsq",
    Host:     "localhost",
    Port:     4150,
    ConnSize: 5,
    PoolSize: 2000,
  }
  var s1 []*config.ConnDetail
  var s2 []*config.ConnDetail
  s1 = append(s1, &cd_bj, &cd_bj2)
  //s1 = append(s1, &cd_bj2)
  s2 = append(s2, &cd_sh)
  hsm = map[string][]*config.ConnDetail{
    label_bj: s1,
    label_sh: s2,
  }

  InitNsq(hsm) //测试时表示使用nsq，在origin中使用一次

  //测试读取nsq数据，wait for sdk init connection
  time.Sleep(4 * time.Second)

  clientBj, err := GetNsq(label_bj)
  clientSh, err := GetNsq(label_sh)
  if err != nil {
    panic(err)
  }

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
        ch := producer(label_bj, clientBj, i, true)
        reply := <-ch
        replyChan <- reply

      } else {
        ch := producer(label_sh, clientSh, i, false)
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
  time.Sleep(2 * time.Second)
}

func producer(label string, client *zgonsq, i int, b bool) chan int {
  var err error

  //还需要一个上下文用来控制开出去的goroutine是否超时
  ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
  defer cancel()
  //输入参数：上下文ctx，nsqClientChan里面是client的连接，args具体的查询操作参数

  body := []byte(fmt.Sprintf("msg--%s--%d", label, i))
  var rch chan uint8
  if b == true { //一次发送多条
    bodyMutil := [][]byte{
      body,
      body,
    }
    rch, err = client.ProducerMulti(ctx, label, bodyMutil)

  } else {
    rch, err = client.Producer(ctx, label, body)

  }
  if err != nil {
    panic(err)
  }

  out := make(chan int, 1)
  select {
  case <-ctx.Done():
    fmt.Println(label, "超时")
    out <- 10001
    return out
  case b := <-rch:
    if b == 1 {
      out <- 1

    } else {
      out <- 10001
    }
  }

  return out

}
