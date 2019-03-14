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
	sleepTime = 1000 //goroutine休眠时间为1000毫秒
)

var (
	connChanMap = make(map[string]chan *radix.Pool)
	mu          sync.RWMutex //用于锁定connChanMap
	hsmu        sync.RWMutex
)

type connPool struct {
	label        string
	m            sync.RWMutex
	connChan     chan *radix.Pool
	clients      []*radix.Pool
	connChanChan chan chan *radix.Pool
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
	//connChanMap = make(map[string]chan *radix.Pool)
	ch := make(chan *config.Labelconns)
	go func() {
		for lahosts := range ch {
			//fmt.Println("-----------------", lahosts)
			label := lahosts.Label
			hosts := lahosts.Hosts
			for k, v := range hosts {
				index := fmt.Sprintf("%s:%d", label, k)
				c := &connPool{
					label:        label,
					connChan:     make(chan *radix.Pool, v.PoolSize),
					connChanChan: make(chan chan *radix.Pool, v.ConnSize),
				}
				mu.Lock()
				connChanMap[index] = c.connChan
				mu.Unlock()
				go c.setConnPoolToChan(index, v)
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

func (cp *connPool) setConnPoolToChan(label string, hosts *config.ConnDetail) {
	//每个host:port连接创建50个连接，放入slice中
	go func() {
		for sessionCh := range cp.connChanChan {
			if session, ok := <-sessionCh; ok {
				//保存channel中的连接到数组中
				cp.clients = append(cp.clients, session)
			}
		}
	}()

	for i := 0; i < 10; i++ {
		//把并发创建的数据库的连接channel，放进channel中
		cp.connChanChan <- cp.createClient(fmt.Sprintf("%s:%d", hosts.Host, hosts.Port), hosts.Db, hosts.PoolSize, hosts.Password)
	}

	go func() {
		for {
			//如果连接全部创建完成，且channel中有了足够的mchSize个连接；循环确保channel中有连接
			//mchSize越大，越用不完，会休眠越久，不用长时间塞连接进channel
			if len(cp.connChan) < hosts.PoolSize && len(cp.clients) >= 1 {
				for _, s := range cp.clients {
					if s != nil {
						cp.connChan <- s
					}
				}

			} else {
				//大多时间是在执行下面一行sleep
				time.Sleep(sleepTime * time.Millisecond)
				//fmt.Println(len(cp.connChan), "--connChan--", label, hosts.Host, hosts.Port)
				//fmt.Println(len(connChanMap), "--connChanMap--", label, hosts.Host, hosts.Port)
			}
		}

	}()

	go func() {
		time.Sleep(2000 * time.Millisecond) //仅仅为了查看创建的连接数，创建数据库连接时间：90ms
		fmt.Printf("init Redis to Channel [%d] ... [%s] Host:%s, Port:%d, Conn:%d, Pool:%d, %s\n",
			len(cp.connChan), label, hosts.Host, hosts.Port, hosts.ConnSize, hosts.PoolSize, hosts.C)
	}()
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
func (cp *connPool) createClient(address string, db int, poolsize int, password string) chan *radix.Pool {
	out := make(chan *radix.Pool)
	go func() {
		//stime := time.Now()
		//fmt.Println(address,db,poolsize,password)
		customConnFunc := func(network, addr string) (radix.Conn, error) {
			return radix.Dial(network, addr,
				radix.DialTimeout(30*time.Second), radix.DialSelectDB(db), radix.DialAuthPass(password),
			)
		}
		c, err := radix.NewPool("tcp", address, 10, radix.PoolConnFunc(customConnFunc))
		if err != nil {
			fmt.Println("redis ", err)
			out <- nil
			return
		}
		out <- c
		//fmt.Println(time.Now().Sub(stime))	//创建数据库连接时间：90ms
	}()
	return out
}
