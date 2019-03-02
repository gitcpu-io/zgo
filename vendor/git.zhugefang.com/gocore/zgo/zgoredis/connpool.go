package zgoredis

import (
	"fmt"
	"git.zhugefang.com/gocore/zgo/config"
	"github.com/mediocregopher/radix"
	"math/rand"
	"sync"
	"time"
)

const (
	//limitConn = 50   //如果是连接集群就是每台数据库长连接50个，单机连也是50个
	//mchSize   = 1000 //mchSize越大，越用不完，会休眠越久，不用长时间塞连接进channel
	sleepTime = 100 //goroutine休眠时间为1000毫秒
)

var (
	connChanMap map[string]chan *radix.Pool
	hsmu        sync.RWMutex
)

type connPool struct {
	label string
	m     sync.RWMutex
}

func NewConnPool(label string) *connPool {
	return &connPool{
		label: label,
	}
}

//连接对外的接口
type ConnPooler interface {
	GetConnChan(label string) chan *radix.Pool
}

//InitConnPool 对外暴露
func InitConnPool(hsm map[string][]*config.ConnDetail) {
	initConnPool(hsm)
}

func initConnPool(hsm map[string][]*config.ConnDetail) { //仅跑一次
	hsmu.RLock()
	defer hsmu.RUnlock()
	connChanMap = make(map[string]chan *radix.Pool)
	ch := make(chan *config.Labelconns)
	go func() {
		for lahosts := range ch {
			//fmt.Println("-----------------", lahosts)
			label := lahosts.Label
			hosts := lahosts.Hosts
			for k, v := range hosts {
				index := fmt.Sprintf("%s:%d", label, k)
				c := &connPool{
					label: label,
				}
				connChanMap[index] = c.createClient(fmt.Sprintf("redis://%s:%s@%s:%d", v.Username, v.Password, v.Host, v.Port), v.Db, v.PoolSize)
			}
			//fmt.Println(label, hosts, "hsm==map=====", len(connChanMap), connChanMap)
		}
	}()

	for label, val := range hsm {
		lcs := &config.Labelconns{
			Label: label,
			Hosts: val,
		}
		ch <- lcs
	}
	close(ch)

}

//GetConnChan 通过label并发安全读map
func (cp *connPool) GetConnChan(label string) chan *radix.Pool {
	cp.m.RLock()
	defer cp.m.RUnlock()

	labLen := 0
	if v, ok := currentLabels[label]; ok {
		labLen = len(v)
	}
	index := rand.Intn(labLen) //随机取一个相同label下的连接
	return connChanMap[fmt.Sprintf("%s:%d", label, index)]
}

//createClient 创建客户端连接
func (cp *connPool) createClient(address string, db int, poolsize int) chan *radix.Pool {
	out := make(chan *radix.Pool)
	go func() {
		customConnFunc := func(network, addr string) (radix.Conn, error) {
			return radix.Dial(network, addr,
				radix.DialTimeout(10*time.Second), radix.DialSelectDB(db), radix.DialAuthPass(""),
			)
		}
		c, err := radix.NewPool("tcp", address, poolsize, radix.PoolConnFunc(customConnFunc))
		if err != nil {
			fmt.Println("redis ", err)
		}
		out <- c
		//fmt.Println(time.Now().Sub(stime))	//创建数据库连接时间：90ms
	}()
	return out
}
