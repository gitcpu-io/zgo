package zgoes

import (
	"net/http"
	"sync"
)

const (
	limitConn = 50    //如果是连接集群就是每台数据库长连接50个，单机连也是50个
	mchSize   = 20000 //mchSize越大，越用不完，会休眠越久，不用长时间塞连接进channel
	sleepTime = 1000  //goroutine休眠时间为1000毫秒
)

type EsResource struct {
	EsClient *http.Client
	UIR      string
}

var (
	connChanMap map[string]chan *EsResource
	mu          sync.RWMutex
	hsmu        sync.RWMutex
)
//连接对外的接口
type ConnPooler interface {
	GetConnChan(label string) chan *EsResource
}
type connPool struct {
	label        string
	hosts        []string
	connChan     chan *EsResource
	clients      []*EsResource
	connChanChan chan chan *EsResource
}
