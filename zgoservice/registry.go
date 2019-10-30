package zgoservice

import (
	"context"
	"errors"
	"fmt"
	"git.zhugefang.com/gocore/zgo/zgocrypto"
	"git.zhugefang.com/gocore/zgo/zgolb"
	"go.etcd.io/etcd/clientv3"
	"strings"
	"sync"
	"time"
)

/*
@Time : 2019-10-25 11:49
@Author : rubinus.chu
@File : registry
@project: zgo
*/

const (
	zgoServiceFrefix = "zgo/service"
)

var serviceChan = make(chan string, 1000)
var serviceList = make(map[string]zgolb.WR2er)
var mu = &sync.RWMutex{}

type LBResponse struct {
	SimpleHttpHost string
	SvcHost        string
	SvcHttpPort    string
	SvcGrpcPort    string
}

type RegistryAndDiscover interface {
	Registry(serviceName, svcHost, httpPort, grpcPort string) error
	UnRegistry() error
	//todo 发现当前服务使用的 其它服务名，用哪个就监听哪个
	Discovery(serviceNames []string) error
	LB(serviceName string) (lbRes *LBResponse, err error)
	Watch() chan string
}

//创建租约注册服务
type Service struct {
	LBResponse
	name          string //服务名
	client        *clientv3.Client
	lease         clientv3.Lease
	leaseResp     *clientv3.LeaseGrantResponse
	cancelfunc    func()
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	key           string
}

func NewService(ttl int64, addr []string) (RegistryAndDiscover, error) {
	conf := clientv3.Config{
		Endpoints:   addr,
		DialTimeout: 20 * time.Second,
	}

	var (
		client *clientv3.Client
	)

	if clientTem, err := clientv3.New(conf); err == nil {
		client = clientTem
	} else {
		return nil, err
	}

	service := &Service{
		client: client,
	}

	if err := service.setLease(ttl); err != nil {
		return nil, err
	}
	go service.ListenLeaseRespChan()

	return service, nil
}

func (service *Service) Watch() chan string {
	return serviceChan
}

//设置租约
func (service *Service) setLease(ttl int64) error {
	lease := clientv3.NewLease(service.client)

	//设置租约时间
	leaseResp, err := lease.Grant(context.TODO(), ttl)
	if err != nil {
		return err
	}

	//设置续租
	ctx, cancelFunc := context.WithCancel(context.TODO())
	leaseRespChan, err := lease.KeepAlive(ctx, leaseResp.ID)

	if err != nil {
		return err
	}

	service.lease = lease
	service.leaseResp = leaseResp
	service.cancelfunc = cancelFunc
	service.keepAliveChan = leaseRespChan
	return nil
}

//监听 续租情况
func (service *Service) ListenLeaseRespChan() {
	for {
		select {
		case leaseKeepResp := <-service.keepAliveChan:
			if leaseKeepResp == nil {
				//val := fmt.Sprintf("%s:%s:%s", service.SvcHost, service.SvcHttpPort, service.SvcGrpcPort)
				//ek := zgocrypto.New().Md5(val)
				//key := fmt.Sprintf("%s/%s/%s", zgoServiceFrefix, service.name, ek)
				//service.delServiceList(service.name, key, val)

				//fmt.Printf("\n%s，Service is Terminated\n", service.name)

				//return

			} else {
				//fmt.Printf("续租成功\n")
			}
		}
	}
}

//通过租约 注册服务
func (service *Service) Registry(serviceName, svcHost, httpPort, grpcPort string) error {
	kv := clientv3.NewKV(service.client)
	val := fmt.Sprintf("%s:%s:%s", svcHost, httpPort, grpcPort)
	key := zgocrypto.New().Md5(val)
	newKey := fmt.Sprintf("%s/%s/%s", zgoServiceFrefix, serviceName, key)
	_, err := kv.Put(context.TODO(), newKey, val, clientv3.WithLease(service.leaseResp.ID))
	service.name = serviceName
	service.SvcHost = svcHost
	service.SvcHttpPort = httpPort
	service.SvcGrpcPort = grpcPort
	return err
}

//撤销租约
func (service *Service) UnRegistry() error {
	if service.leaseResp == nil || service.leaseResp.Error != "" {
		return nil
	}
	service.cancelfunc()
	_, err := service.lease.Revoke(context.TODO(), service.leaseResp.ID)
	return err
}

//内部负载均衡
func (service *Service) LB(serviceName string) (lbRes *LBResponse, err error) {
	mu.RLock()
	defer mu.RUnlock()
	key := fmt.Sprintf("%s/%s", zgoServiceFrefix, serviceName)
	lb := serviceList[key]
	if lb == nil {
		return nil, errors.New("服务map创建失败 " + serviceName)
	}
	res, err := lb.Balance()
	if err != nil {
		return nil, err
	}
	split := strings.Split(res, ":")
	lbRes = &LBResponse{}
	lbRes.SvcHost = split[0]
	lbRes.SimpleHttpHost = fmt.Sprintf("http://%s", lbRes.SvcHost)
	if len(split) > 0 {
		lbRes.SvcHttpPort = split[1]
		lbRes.SimpleHttpHost = fmt.Sprintf("%s:%s", lbRes.SimpleHttpHost, lbRes.SvcHttpPort)
	}
	if len(split) > 1 {
		lbRes.SvcGrpcPort = split[2]
	}

	return lbRes, nil

}

func (service *Service) Discovery(serviceNames []string) error {
	wg := &sync.WaitGroup{}
	for _, serviceName := range serviceNames {
		wg.Add(1)
		key := fmt.Sprintf("%s/%s", zgoServiceFrefix, serviceName)
		go service.getService(key, wg)
	}
	wg.Wait()
	return nil
}

func (service *Service) getService(prefix string, wg *sync.WaitGroup) (addrs []string, err error) {
	defer wg.Done()

	resp, err := service.client.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	addrs = service.getAddrs(resp)

	mu.Lock()
	defer mu.Unlock()

	serviceList[prefix] = zgolb.NewWR2ByArr(addrs) //一个frefix一个lb

	//fmt.Println("最新map：",serviceList[prefix], prefix)

	go service.watcher(prefix)

	return addrs, nil
}

func (service *Service) watcher(prefix string) {
	rch := service.client.Watch(context.Background(), prefix, clientv3.WithPrefix(), clientv3.WithPrevKV())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			key := strings.Split(string(ev.Kv.Key), "/")
			serviceName := key[2]
			newKey := fmt.Sprintf("%s/%s", zgoServiceFrefix, serviceName)
			switch ev.Type {
			case clientv3.EventTypePut:
				service.setServiceList(serviceName, newKey, string(ev.Kv.Value))
			case clientv3.EventTypeDelete:
				service.delServiceList(serviceName, newKey, string(ev.PrevKv.Value))
			}
		}
	}
}

func (service *Service) getAddrs(resp *clientv3.GetResponse) []string {
	addrs := make([]string, 0)
	if resp == nil || resp.Kvs == nil {
		return addrs
	}
	for i := range resp.Kvs {
		if v := resp.Kvs[i].Value; v != nil {
			addrs = append(addrs, string(v))
		}
	}
	return addrs
}

//setServiceList 添加到map 从LB中
func (service *Service) setServiceList(serviceName, key, val string) {
	if serviceList[key] == nil {
		serviceList[key] = zgolb.NewWR2ByArr([]string{})
	}
	mu.Lock()
	defer mu.Unlock()
	serviceList[key].Add(string(val))

	//发送到channel 通知
	serviceChan <- serviceName
}

//delServiceList 从map中找到key，并从LB中删除
func (service *Service) delServiceList(serviceName, key, val string) {
	if serviceList[key] == nil {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	serviceList[key].Remove(val)

	//发送到channel 通知
	serviceChan <- serviceName
}
