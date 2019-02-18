package zgo_db_mongo

import (
	"fmt"
	"github.com/globalsign/mgo"
	"time"
)

func MongoClientChan() chan *mgo.Session {
	return MongoChan
}

const limitConn = 50

var (
	MongoChan   chan *mgo.Session
	mgoSessions []*mgo.Session
)

func init() {
	MongoChan = make(chan *mgo.Session, 10000)
	mgoSessions = []*mgo.Session{}

	addr := []string{
		"127.0.0.1:27017",
	}

	//每个host:port连接创建50个连接，放入slice中
	ssChanChan := make(chan chan *mgo.Session, limitConn*len(addr))
	go func() {
		for sessionCh := range ssChanChan {
			if session, ok := <-sessionCh; ok {
				mgoSessions = append(mgoSessions, session)
			}
		}
	}()

	for i := 0; i < limitConn; i++ {
		for _, host := range addr {
			ssChanChan <- createConnection(host)
		}
	}

	go func() {
		for {
			if len(MongoChan) < 10000 {
				for _, s := range mgoSessions {
					if s != nil {
						MongoChan <- s
					}
				}
			}
			time.Sleep(limitConn * time.Millisecond)
			//fmt.Println(len(MongoChan), "--MongoChan--")
		}

	}()
	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println("init mongo connection to MongoChan ...", len(MongoChan))
	}()

}

func createConnection(host string) chan *mgo.Session {
	out := make(chan *mgo.Session)
	go func() {
		dialInfo := mgo.DialInfo{
			Addrs: []string{host},
			//Database: "local",
			//Username: username,
			//Password: password,
			//PoolLimit: 50000,
			Timeout: time.Duration(60 * time.Second),
		}
		//fmt.Println("ing...",host)
		session, err := mgo.DialWithInfo(&dialInfo)
		//fmt.Println("done...",host)
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
	}()
	return out

}
