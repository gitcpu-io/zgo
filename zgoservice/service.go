package zgoservice

import (
	"github.com/gitcpu-io/zgo/config"
	"strings"
)

/*
@Time : 2019-10-25 11:43
@Author : rubinus.chu
@File : service
@project: zgo
*/

type Servicer interface {
	//默认使用 zgo engine的etcd，heartbeat为心跳时间间隔，单位是秒
	New(ttl int64, addr string) (RegistryAndDiscover, error)
	LB(serviceName string) (lbRes *LBResponse, err error)
	Watch() chan string
}

type service struct {
	res RegistryAndDiscover
}

func GetService(ttl int64, addr []string) (RegistryAndDiscover, error) {
	newService, err := NewService(ttl, addr)
	if err != nil {
		return nil, err
	}
	return newService, nil
}

func (s *service) New(ttl int64, addr string) (RegistryAndDiscover, error) {
	var addrs []string
	if addr == "" {
		addrs = config.Conf.EtcdHosts
	} else {
		addrs = strings.Split(addr, ",")
	}
	res, err := GetService(ttl, addrs)
	s.res = res
	return res, err
}

func (s *service) LB(serviceName string) (lbRes *LBResponse, err error) {
	return s.res.LB(serviceName)
}

func (s *service) Watch() chan string {
	return s.res.Watch()
}

func New() *service {
	return &service{}
}
