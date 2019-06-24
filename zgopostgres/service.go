// zgopostgres是对消息中间件Postgres的封装，提供新建连接，生产数据，消费数据接口
package zgopostgres

import (
	"git.zhugefang.com/gocore/zgo/comm"
	"git.zhugefang.com/gocore/zgo/config"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"sync"
)

var (
	currentLabels = make(map[string][]*config.ConnDetail) //用于存放label与具体Host:port的map
	muLabel       sync.RWMutex                            //用于并发读写上面的map
)

//Postgres 对外
type Postgreser interface {
	/*
	 label: 可选，如果使用者，用了2个或多个label时，需要调用这个函数，传入label
	*/
	// New 生产一条消息到Postgres
	New(label ...string) (*zgopostgres, error)

	/*
	 label: 可选，如果使用者，用了2个或多个label时，需要调用这个函数，传入label
	*/
	// GetConnChan 获取原生的生产者client，返回一个chan，使用者需要接收 <- chan
	GetConnChan(label ...string) (chan *pg.DB, error)

	Scan(values ...interface{}) orm.ColumnScanner
}

// Postgres用于对zgo.Postgres这个全局变量赋值
func Postgres(label string) Postgreser {
	return &zgopostgres{
		res: NewPostgresResourcer(label),
	}
}

// zgopostgres实现了Postgres的接口
type zgopostgres struct {
	res PostgresResourcer //使用resource另外的一个接口
}

// InitPostgres 初始化连接postgres，用于使用者zgo.engine时，zgo init
func InitPostgres(hsmIn map[string][]*config.ConnDetail, label ...string) chan *zgopostgres {
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

	InitPostgresResource(hsm)

	//自动为变量初始化对象
	initLabel := ""
	for k, _ := range hsm {
		if k != "" {
			initLabel = k
			break
		}
	}
	out := make(chan *zgopostgres)
	go func() {

		in, err := GetPostgres(initLabel)
		if err != nil {
			panic(err)
		}
		out <- in
		close(out)
	}()

	return out

}

// GetPostgres zgo内部获取一个连接postgres
func GetPostgres(label ...string) (*zgopostgres, error) {
	l, err := comm.GetCurrentLabel(label, muLabel, currentLabels)
	if err != nil {
		return nil, err
	}
	return &zgopostgres{
		res: NewPostgresResourcer(l),
	}, nil
}

// NewPostgres获取一个Postgres生产者的client，用于发送数据
func (n *zgopostgres) New(label ...string) (*zgopostgres, error) {
	return GetPostgres(label...)
}

//GetConnChan 供用户使用原生连接的chan
func (n *zgopostgres) GetConnChan(label ...string) (chan *pg.DB, error) {
	l, err := comm.GetCurrentLabel(label, muLabel, currentLabels)
	if err != nil {
		return nil, err
	}
	return n.res.GetConnChan(l), nil
}

func (n *zgopostgres) Scan(values ...interface{}) orm.ColumnScanner {
	return n.res.Scan(values...)
}
