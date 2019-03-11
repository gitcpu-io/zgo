package zgokafka

import (
	"fmt"
	"git.zhugefang.com/gocore/zgo/config"
	"github.com/Shopify/sarama"
	"math/rand"
	"sync"
	"time"
)

const (
	//limitConn = 50    //如果是连接集群就是每台数据库长连接50个，单机连也是50个
	//mchSize   = 20000 //mchSize越大，越用不完，会休眠越久，不用长时间塞连接进channel
	sleepTime = 1000 //goroutine休眠时间为1000毫秒
)

var (
	connChanMap map[string]chan *sarama.AsyncProducer
	mu          sync.RWMutex //用于锁定connChanMap
	hsmu        sync.RWMutex
)

//连接对外的接口
type ConnPooler interface {
	GetConnChan(label string) chan *sarama.AsyncProducer
}

type connPool struct {
	label        string
	m            sync.RWMutex
	connChan     chan *sarama.AsyncProducer
	clients      []*sarama.AsyncProducer
	connChanChan chan chan *sarama.AsyncProducer
}

func NewConnPool(label string) *connPool {
	return &connPool{
		label: label,
	}
}

//InitConnPool 对外暴露
func InitConnPool(hsm map[string][]*config.ConnDetail) {
	initConnPool(hsm)
}

func initConnPool(hsm map[string][]*config.ConnDetail) { //仅跑一次
	hsmu.RLock()
	defer hsmu.RUnlock()

	connChanMap = make(map[string]chan *sarama.AsyncProducer)

	ch := make(chan *config.Labelconns)
	go func() {
		for lahosts := range ch {
			label := lahosts.Label
			hosts := lahosts.Hosts

			for k, v := range hosts {
				index := fmt.Sprintf("%s:%d", label, k)
				c := &connPool{
					label:        label,
					connChan:     make(chan *sarama.AsyncProducer, v.PoolSize),
					connChanChan: make(chan chan *sarama.AsyncProducer, v.ConnSize),
				}
				mu.Lock()
				connChanMap[index] = c.connChan
				mu.Unlock()
				go c.setConnPoolToChan(index, v) //call 创建连接到chan中
			}

			//fmt.Println(label, hosts, "hsm=====",len(hsm), connChanMap)
		}
	}()

	for label, val := range hsm {
		lcs := &config.Labelconns{
			Label: label,
			Hosts: val,
		}
		//fmt.Println(lcs,"-----")
		ch <- lcs
	}
	close(ch)

}

//GetConnChan 通过label并发安全读map
func (cp *connPool) GetConnChan(label string) chan *sarama.AsyncProducer {
	cp.m.RLock()
	defer cp.m.RUnlock()

	labLen := 0
	if v, ok := currentLabels[label]; ok {
		labLen = len(v)
	}
	index := rand.Intn(labLen) //随机取一个相同label下的连接
	return connChanMap[fmt.Sprintf("%s:%d", label, index)]
}

func (cp *connPool) setConnPoolToChan(label string, hosts *config.ConnDetail) {
	//fmt.Sprintf(label, "--", hosts)
	//每个host:port连接创建50个连接，放入slice中
	go func() {
		for sessionCh := range cp.connChanChan {
			if session, ok := <-sessionCh; ok {
				//保存channel中的连接到数组中
				cp.clients = append(cp.clients, session)
			}
		}
	}()

	for i := 0; i < hosts.ConnSize; i++ {
		//把并发创建的数据库的连接channel，放进channel中
		hh := []string{
			fmt.Sprintf("%s:%d", hosts.Host, hosts.Port),
		}
		cp.connChanChan <- cp.createClient(hh)
	}

	go func() {
		for {
			//如果连接全部创建完成，且channel中有了足够的mchSize个连接；循环确保channel中有连接
			//mchSize越大，越用不完，会休眠越久，不用长时间塞连接进channel
			if len(cp.connChan) < hosts.PoolSize && len(cp.clients) >= hosts.ConnSize/2 {
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
		fmt.Printf("init Kafka to Channel [%d] ... [%s] Host:%s, Port:%d, Conn:%d, Pool:%d, %s\n",
			len(cp.connChan), label, hosts.Host, hosts.Port, hosts.ConnSize, hosts.PoolSize, hosts.C)
	}()
}

//createClient 创建客户端连接
func (cp *connPool) createClient(address []string) chan *sarama.AsyncProducer {
	out := make(chan *sarama.AsyncProducer)
	go func() {
		c := sarama.NewConfig()
		c.Producer.Return.Successes = true
		p, err := sarama.NewAsyncProducer(address, c)
		if err != nil {
			fmt.Printf("sarama.NewSyncProducer err:%s \n", err)
			out <- nil
			return
		}
		out <- &p
	}()
	return out
}

//func CreateClient(address []string) chan *sarama.AsyncProducer {
//	out := make(chan *sarama.AsyncProducer)
//	go func() {
//		c := sarama.NewConfig()
//		c.Producer.Return.Successes = true
//		p, err := sarama.NewAsyncProducer(address, c)
//		if err != nil {
//			fmt.Printf("sarama.NewSyncProducer err:%s \n", err)
//			out <- nil
//			return
//		}
//		out <- &p
//	}()
//	return out
//}
