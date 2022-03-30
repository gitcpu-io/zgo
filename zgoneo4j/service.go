// zgoneo4j是对中间件Neo4j的封装，提供新建连接
package zgoneo4j

import (
  "github.com/gitcpu-io/zgo/comm"
  "github.com/gitcpu-io/zgo/config"
  "github.com/neo4j/neo4j-go-driver/neo4j"
  "sync"
)

var (
  currentLabels = make(map[string][]*config.ConnDetail) //用于存放label与具体Host:port的map
  muLabel       *sync.RWMutex                            //用于并发读写上面的map
)

//Neo4j 对外
type Neo4jer interface {
  /*
   label: 可选，如果使用者，用了2个或多个label时，需要调用这个函数，传入label
  */
  // New 生产一条消息到Neo4j
  New(label ...string) (*zgoneo4j, error)

  /*
   label: 可选，如果使用者，用了2个或多个label时，需要调用这个函数，传入label
  */
  // GetConnChan 获取原生的生产者client，返回一个chan，使用者需要接收 <- chan
  GetConnChan(label ...string) (chan neo4j.Session, error)
}

// Neo4j用于对zgo.Neo4j这个全局变量赋值
func Neo4j(label string) Neo4jer {
  return &zgoneo4j{
    res: NewNeo4jResourcer(label),
  }
}

// zgoneo4j实现了Neo4j的接口
type zgoneo4j struct {
  res Neo4jResourcer //使用resource另外的一个接口
}

// InitNeo4j 初始化连接，用于使用者zgo.engine时，zgo init
func InitNeo4j(hsmIn map[string][]*config.ConnDetail, label ...string) chan *zgoneo4j {
  muLabel.Lock()
  defer muLabel.Unlock()

  var hsm map[string][]*config.ConnDetail

  if len(label) > 0 && len(currentLabels) > 0 { //此时是destory操作,传入的hsm是nil
    //fmt.Println("--destory--前",currentLabels)
    for _, v := range label {
      delete(currentLabels, v)
    }
    hsm = currentLabels
    //fmt.Println("--destory--后",currentLabels)

  } else { //这是第一次创建操作或etcd中变更时init again操作
    hsm = hsmIn
    //currentLabels = hsm	//this operation is error
    for k, v := range hsm { //so big bug can't set hsm to currentLabels，must be for, may be have old label
      currentLabels[k] = v
    }
  }

  if len(hsm) == 0 {
    return nil
  }

  InitNeo4jResource(hsm)

  //自动为变量初始化对象
  initLabel := ""
  for k := range hsm {
    if k != "" {
      initLabel = k
      break
    }
  }
  out := make(chan *zgoneo4j)
  go func() {

    in, err := GetNeo4j(initLabel)
    if err != nil {
      panic(err)
    }
    out <- in
    close(out)
  }()

  return out

}

// GetNeo4j zgo内部获取一个连接
func GetNeo4j(label ...string) (*zgoneo4j, error) {
  l, err := comm.GetCurrentLabel(label, muLabel, currentLabels)
  if err != nil {
    return nil, err
  }
  return &zgoneo4j{
    res: NewNeo4jResourcer(l),
  }, nil
}

// NewNeo4j获取一个Neo4j生产者的client，用于发送数据
func (n *zgoneo4j) New(label ...string) (*zgoneo4j, error) {
  return GetNeo4j(label...)
}

//GetConnChan 供用户使用原生连接的chan
func (n *zgoneo4j) GetConnChan(label ...string) (chan neo4j.Session, error) {
  l, err := comm.GetCurrentLabel(label, muLabel, currentLabels)
  if err != nil {
    return nil, err
  }
  return n.res.GetConnChan(l), nil
}
