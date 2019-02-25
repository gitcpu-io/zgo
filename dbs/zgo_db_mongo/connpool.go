package zgo_db_mongo

import (
	"fmt"
	"github.com/globalsign/mgo"

	"time"
)

func MongoClientChan() chan *mgo.Session {
	return MongoChan
}

const (
	limitConn = 50    //如果是连接集群就是每台数据库长连接50个，单机连也是50个
	mchSize   = 20000 //mchSize越大，越用不完，会休眠越久，不用长时间塞连接进channel
	sleepTime = 1000  //goroutine休眠时间为1000毫秒
)

var (
	MongoChan   chan *mgo.Session
	mgoSessions []*mgo.Session
)

func init() {
	Init()
}

func Init() {

	MongoChan = make(chan *mgo.Session, mchSize)
	mgoSessions = []*mgo.Session{}

	addr := []string{
		"127.0.0.1:27017",
	}

	//每个host:port连接创建50个连接，放入slice中
	ssChanChan := make(chan chan *mgo.Session, limitConn*len(addr))
	go func() {
		for sessionCh := range ssChanChan {
			if session, ok := <-sessionCh; ok {
				//保存channel中的连接到数组中
				mgoSessions = append(mgoSessions, session)
			}
		}
	}()

	for i := 0; i < limitConn; i++ {
		for _, host := range addr {
			//把并发创建的数据库的连接channel，放进channel中
			ssChanChan <- createConnection(host)
		}
	}

	go func() {
		for {
			//如果连接全部创建完成，且channel中有了足够的mchSize个连接；循环确保channel中有连接
			//mchSize越大，越用不完，会休眠越久，不用长时间塞连接进channel
			if len(MongoChan) < mchSize && len(mgoSessions) == limitConn {
				for _, s := range mgoSessions {
					if s != nil {
						MongoChan <- s
					}
				}

			} else {
				//大多时间是在执行下面一行sleep
				time.Sleep(sleepTime * time.Millisecond)
				//fmt.Println(len(MongoChan), "--MongoChan--")
			}
		}

	}()

	go func() {
		time.Sleep(2000 * time.Millisecond) //仅仅为了查看创建的连接数，创建数据库连接时间：90ms
		fmt.Println("init mongo connection to MongoChan ...", len(MongoChan))
	}()
}

//createConnection 并发创建数据库连接
func createConnection(host string) chan *mgo.Session {
	out := make(chan *mgo.Session)
	go func() {
		//stime := time.Now()

		dialInfo := mgo.DialInfo{
			Addrs: []string{host},
			//Database: "local",
			//Username: username,
			//Password: password,
			//PoolLimit: 50000,
			Timeout: time.Duration(60 * time.Second),
		}

		session, err := mgo.DialWithInfo(&dialInfo)

		if err != nil || session == nil {
			fmt.Println(session, err)
			out <- nil
			return
		}
		session.SetMode(mgo.Monotonic, true)
		session.SetSafe(&mgo.Safe{
			WMode: "majority",
		})
		out <- session
		//fmt.Println(time.Now().Sub(stime))	//创建数据库连接时间：90ms
	}()
	return out
}
