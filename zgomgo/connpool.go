package zgomgo

import (
	"context"
	"fmt"
	"git.zhugefang.com/gocore/zgo/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
	"sync"
	"time"
)

const (
	sleepTime = 1000 //goroutine休眠时间为1000毫秒
)

var (
	connChanMap = make(map[string]chan *mongo.Client)
	mu          sync.RWMutex //用于锁定connChanMap
	hsmu        sync.RWMutex
)

//连接对外的接口
type ConnPooler interface {
	GetConnChan(label string) chan *mongo.Client
}

type connPool struct {
	label        string
	m            sync.RWMutex
	connChan     chan *mongo.Client
	clients      []*mongo.Client
	connChanChan chan chan *mongo.Client
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

	ch := make(chan *config.Labelconns)
	go func() {
		for lahosts := range ch {
			label := lahosts.Label
			hosts := lahosts.Hosts

			for k, v := range hosts {
				index := fmt.Sprintf("%s:%d", label, k)
				c := &connPool{
					label:        label,
					connChan:     make(chan *mongo.Client, v.PoolSize),
					connChanChan: make(chan chan *mongo.Client, v.ConnSize),
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
func (cp *connPool) GetConnChan(label string) chan *mongo.Client {
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
		var address string
		if hosts.Username != "" && hosts.Password != "" {
			address = fmt.Sprintf("mongodb://%s:%s@%s:%d", hosts.Username, hosts.Password, hosts.Host, hosts.Port)

		} else {
			address = fmt.Sprintf("mongodb://%s:%d", hosts.Host, hosts.Port)

		}
		cp.connChanChan <- cp.createClient(address, hosts.Username, hosts.Password, hosts.DbName, hosts.PoolSize)
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
		fmt.Printf("init Mongo to Channel [%d] ... [%s] Host:%s, Port:%d, Conn:%d, Pool:%d, %s\n",
			len(cp.connChan), label, hosts.Host, hosts.Port, hosts.ConnSize, hosts.PoolSize, hosts.C)
	}()
}

//createClient 创建客户端连接
func (cp *connPool) createClient(address, username, password, dbname string, poolSize int) chan *mongo.Client {
	out := make(chan *mongo.Client)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		opts := &options.ClientOptions{
			//Hosts: []string{
			//	"47.95.20.12:27018",
			//	"47.95.20.12:27019",
			//},
			Auth: &options.Credential{
				AuthSource: dbname,
				Username:   username,
				Password:   password,
			},
		}
		opts.SetMaxPoolSize(uint64(poolSize))
		c, err := mongo.Connect(ctx, options.Client().ApplyURI(address), opts)
		//c, err := mongo.Connect(ctx, opts)
		if err != nil {
			fmt.Println("connection Mongo failed :", address, err)
			out <- nil
			return
		}
		out <- c
	}()
	return out
}
