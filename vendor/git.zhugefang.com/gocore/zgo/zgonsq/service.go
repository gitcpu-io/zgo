package zgonsq

import (
	"context"
	"git.zhugefang.com/gocore/zgo/comm"
	"git.zhugefang.com/gocore/zgo/config"
	"github.com/nsqio/go-nsq"
	"sync"
)

var (
	currentLabels = make(map[string][]*config.ConnDetail)
	muLabel       sync.RWMutex
)

//Nsq 对外
type Nsqer interface {
	NewNsq(label ...string) (*zgonsq, error)
	GetConnChan(label ...string) (chan *nsq.Producer, error)
	Producer(ctx context.Context, topic string, body []byte) (chan uint8, error)
	ProducerMulti(ctx context.Context, topic string, body [][]byte) (chan uint8, error)
	Consumer(topic, channel string, mode int, fn NsqHandlerFunc)
}

func Nsq(label string) Nsqer {
	return &zgonsq{
		res: NewNsqResourcer(label),
	}
}

//zgonsq实现了Nsq的接口
type zgonsq struct {
	res NsqResourcer //使用resource另外的一个接口
}

//InitNsq 初始化连接nsq
func InitNsq(hsm map[string][]*config.ConnDetail) chan *zgonsq {
	muLabel.Lock()
	defer muLabel.Unlock()

	currentLabels = hsm
	InitNsqResource(hsm)

	//自动为变量初始化对象
	initLabel := ""
	for k, _ := range hsm {
		if k != "" {
			initLabel = k
			break
		}
	}
	out := make(chan *zgonsq)
	go func() {

		in, err := GetNsq(initLabel)
		if err != nil {
			panic(err)
		}
		out <- in
		close(out)
	}()

	return out

}

//GetNsq zgo内部获取一个连接nsq
func GetNsq(label ...string) (*zgonsq, error) {
	l, err := comm.GetCurrentLabel(label, muLabel, currentLabels)
	if err != nil {
		return nil, err
	}
	return &zgonsq{
		res: NewNsqResourcer(l),
	}, nil
}

//getCurrentLabel 着重判断输入的label与zgo engine中用户方的label
//func getCurrentLabel(label ...string) (string, error) {
//	muLabel.RLock()
//	defer muLabel.RUnlock()
//
//	lcl := len(currentLabels)
//	if lcl == 0 {
//		return "", errors.New("invalid label in zgo engine or engine not start.")
//	}
//	if len(label) == 0 { //用户没有选择
//		if lcl >= 1 {
//			//自动返回默认的第一个
//			l := ""
//			for k, _ := range currentLabels {
//				l = k
//				break
//			}
//			return l, nil
//		} else {
//			return "", errors.New("invalid label in zgo engine.")
//		}
//	} else if len(label) > 1 {
//		return "", errors.New("you are choose must be one label or defalut zero.")
//	} else {
//		if _, ok := currentLabels[label[0]]; ok {
//			return label[0], nil
//		} else {
//			return "", errors.New("invalid label for u input.")
//		}
//	}
//}

func (n *zgonsq) NewNsq(label ...string) (*zgonsq, error) {
	return GetNsq(label...)
}

//GetConnChan 供用户使用原生连接的chan
func (n *zgonsq) GetConnChan(label ...string) (chan *nsq.Producer, error) {
	l, err := comm.GetCurrentLabel(label, muLabel, currentLabels)
	if err != nil {
		return nil, err
	}
	return n.res.GetConnChan(l), nil
}

func (n *zgonsq) Producer(ctx context.Context, topic string, body []byte) (chan uint8, error) {
	return n.res.Producer(ctx, topic, body)
}

func (n *zgonsq) ProducerMulti(ctx context.Context, topic string, body [][]byte) (chan uint8, error) {
	return n.res.ProducerMulti(ctx, topic, body)
}

func (n *zgonsq) Consumer(topic, channel string, mode int, fn NsqHandlerFunc) {
	go n.res.Consumer(topic, channel, mode, fn)
}
