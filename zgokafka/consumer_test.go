package zgokafka

import (
  "fmt"
  "github.com/gitcpu-io/zgo/config"
  "testing"
  "time"
)

func TestConsumer(t *testing.T) {
  //hsm := make(map[string][]*config.ConnDetail)
  cd_bj := config.ConnDetail{
    C:        "北京主库-----kafka",
    Host:     "localhost",
    Port:     9092,
    ConnSize: 2,
    PoolSize: 100,
  }
  cd_bj2 := config.ConnDetail{
    C:        "北京从库2-----kafka",
    Host:     "localhost",
    Port:     9092,
    ConnSize: 2,
    PoolSize: 100,
  }
  cd_sh := config.ConnDetail{
    C:        "上海主库-----kafka",
    Host:     "localhost",
    Port:     9092,
    ConnSize: 2,
    PoolSize: 100,
  }
  var s1 []*config.ConnDetail
  var s2 []*config.ConnDetail
  s1 = append(s1, &cd_bj, &cd_bj2)
  s2 = append(s2, &cd_sh)
  hsm := map[string][]*config.ConnDetail{
    label_bj: s1,
    label_sh: s2,
  }
  InitKafka(hsm) //测试时表示使用kafka，在origin中使用一次

  labelBj, err := GetKafka(label_bj)
  if err != nil {
    panic(err)
  }
  labelSh, err := GetKafka(label_sh)
  if err != nil {
    panic(err)
  }
  c := chat{
    Topic:   label_bj,
    GroupId: label_bj,
    Kafka:   labelBj,
  }
  go c.Consumer()

  c2 := chat{
    Topic:   label_sh,
    GroupId: label_sh,
    Kafka:   labelSh,
  }
  go c2.Consumer()
  for val := range time.Tick(time.Duration(3 * time.Second)) {
    fmt.Println("一直在消费着",val)
  }
  //time.Sleep(3 * time.Second)
}
