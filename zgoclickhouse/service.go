// zgoclickhouse是对clickhouse的封装，提供新建连接，生产数据，消费数据接口
package zgoclickhouse

import (
	"database/sql"
	"github.com/rubinus/zgo/comm"
	"github.com/rubinus/zgo/config"
	"sync"
)

var (
	currentLabels = make(map[string][]*config.ConnDetail) //用于存放label与具体Host:port的map
	muLabel       sync.RWMutex                            //用于并发读写上面的map
)

//ClickHouse 对外
type ClickHouseer interface {
	/*
	 label: 可选，如果使用者，用了2个或多个label时，需要调用这个函数，传入label
	*/
	// New 生产一条消息到ClickHouse
	New(label ...string) (*zgoclickhouse, error)

	/*
	 label: 可选，如果使用者，用了2个或多个label时，需要调用这个函数，传入label
	*/
	// GetConnChan 获取原生的生产者client，返回一个chan，使用者需要接收 <- chan
	GetConnChan(label ...string) (chan *sql.DB, error)
}

// ClickHouse用于对zgo.ClickHouse这个全局变量赋值
func ClickHouse(label string) ClickHouseer {
	return &zgoclickhouse{
		res: NewClickHouseResourcer(label),
	}
}

// zgoclickhouse实现了ClickHouse的接口
type zgoclickhouse struct {
	res ClickHouseResourcer //使用resource另外的一个接口
}

// InitClickHouse 初始化连接postgres，用于使用者zgo.engine时，zgo init
func InitClickHouse(hsmIn map[string][]*config.ConnDetail, label ...string) chan *zgoclickhouse {
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

	InitClickHouseResource(hsm)

	//自动为变量初始化对象
	initLabel := ""
	for k, _ := range hsm {
		if k != "" {
			initLabel = k
			break
		}
	}
	out := make(chan *zgoclickhouse)
	go func() {

		in, err := GetClickHouse(initLabel)
		if err != nil {
			panic(err)
		}
		out <- in
		close(out)
	}()

	return out

}

// GetClickHouse zgo内部获取一个连接postgres
func GetClickHouse(label ...string) (*zgoclickhouse, error) {
	l, err := comm.GetCurrentLabel(label, muLabel, currentLabels)
	if err != nil {
		return nil, err
	}
	return &zgoclickhouse{
		res: NewClickHouseResourcer(l),
	}, nil
}

// NewClickHouse获取一个ClickHouse生产者的client，用于发送数据
func (n *zgoclickhouse) New(label ...string) (*zgoclickhouse, error) {
	return GetClickHouse(label...)
}

//GetConnChan 供用户使用原生连接的chan
func (n *zgoclickhouse) GetConnChan(label ...string) (chan *sql.DB, error) {
	l, err := comm.GetCurrentLabel(label, muLabel, currentLabels)
	if err != nil {
		return nil, err
	}
	return n.res.GetConnChan(l), nil
}
