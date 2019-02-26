package zgoredis

import (
	"fmt"
	"github.com/mediocregopher/radix"
	"sync"
	"time"
)

const (
	limitConn = 50   //如果是连接集群就是每台数据库长连接50个，单机连也是50个
	mchSize   = 1000 //mchSize越大，越用不完，会休眠越久，不用长时间塞连接进channel
	sleepTime = 100  //goroutine休眠时间为1000毫秒
)

type connPool struct {
	label        string
	hosts        []string
	connChan     chan *radix.Pool
	clients      []*radix.Pool
	connChanChan chan chan *radix.Pool
}

func NewConnPool(label string) *connPool {
	return &connPool{
		label: label,
	}
}

var (
	connChanMap map[string]chan *radix.Pool
	mu          sync.RWMutex
	hsmu        sync.RWMutex
)

//连接对外的接口
type ConnPooler interface {
	GetConnChan(label string) chan *radix.Pool
}

//InitConnPool 对外暴露
func InitConnPool(hsm map[string][]string) {
	initConnPool(hsm)
}

func initConnPool(hsm map[string][]string) { //仅跑一次
	hsmu.RLock()
	defer hsmu.RUnlock()
	connChanMap = make(map[string]chan *radix.Pool)
	for label, hosts := range hsm {
		//fmt.Println(label)
		//fmt.Println(hosts)
		go func(label string, v []string, mu sync.RWMutex) { //并发连接多个不同的地址和端口
			mu.Lock()
			defer mu.RLock()
			c := &connPool{
				label:        label,
				hosts:        hosts,
				connChan:     make(chan *radix.Pool, mchSize),
				connChanChan: make(chan chan *radix.Pool, limitConn*len(hosts)),
			}
			connChanMap[label] = c.connChan
			c.setConnPoolToChan(label) //call 创建连接到chan中
			//fmt.Println(label, hosts, "hsm=====",len(hsm), connChanMap)

		}(label, hosts, mu)

	}

}

//GetConnChan 通过label并发安全读map
func (cp *connPool) GetConnChan(label string) chan *radix.Pool {
	mu.RLock()
	defer mu.RUnlock()
	return connChanMap[label]
}

func (cp *connPool) setConnPoolToChan(label string) {
	//每个host:port连接创建50个连接，放入slice中
	go func() {
		for sessionCh := range cp.connChanChan {
			if session, ok := <-sessionCh; ok {
				//保存channel中的连接到数组中
				cp.clients = append(cp.clients, session)
			}
		}
	}()

	for i := 0; i < limitConn; i++ {
		for _, host := range cp.hosts {
			//把并发创建的数据库的连接channel，放进channel中
			cp.connChanChan <- cp.createClient(host)
		}
	}

	go func() {
		for {
			//如果连接全部创建完成，且channel中有了足够的mchSize个连接；循环确保channel中有连接
			//mchSize越大，越用不完，会休眠越久，不用长时间塞连接进channel
			if len(cp.connChan) < mchSize && len(cp.clients) == limitConn {
				for _, s := range cp.clients {
					if s != nil {
						cp.connChan <- s
					}
				}

			} else {
				//大多时间是在执行下面一行sleep
				time.Sleep(sleepTime * time.Millisecond)
				//fmt.Println(len(connChanMap), "--connChanMap--", label)
			}
		}

	}()

	go func() {
		time.Sleep(2000 * time.Millisecond) //仅仅为了查看创建的连接数，创建数据库连接时间：90ms
		fmt.Println("init redis connection to connChan ...", len(cp.connChan), label)
	}()
}

//createClient 创建客户端连接
func (cp *connPool) createClient(address string) chan *radix.Pool {
	out := make(chan *radix.Pool)
	go func() {
		customConnFunc := func(network, addr string) (radix.Conn, error) {
			return radix.Dial(network, addr,
				radix.DialTimeout(10*time.Second), radix.DialSelectDB(0), radix.DialAuthPass(""),
			)
		}
		c, err := radix.NewPool("tcp", address, 100, radix.PoolConnFunc(customConnFunc))
		if err != nil {
			fmt.Println("redis ", err)
		}
		out <- c
		//fmt.Println(time.Now().Sub(stime))	//创建数据库连接时间：90ms
	}()
	return out
}
