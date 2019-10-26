package zgolb

import (
	"errors"
	"fmt"
	"sync"
)

/*
@Time : 2019-07-30 21:32
@Author : rubinus.chu
@File : wrrlb
@project: zgo
*/

type WR2er interface {
	Len() int
	Add(child string)
	AddWeight(child string, weight int)
	Exist(child string) bool
	Remove(child string)
	Balance() (string, error)
}

type WR2 struct {
	weight int
	childs []string
	lock   *sync.RWMutex
}

func NewWR2(childs ...string) *WR2 {
	return &WR2{weight: 1, childs: childs, lock: &sync.RWMutex{}}
}

func NewWR2ByArr(childs []string) *WR2 {
	return &WR2{weight: 1, childs: childs, lock: &sync.RWMutex{}}
}

func (wr2 *WR2) Len() int {
	wr2.lock.RLock()
	defer wr2.lock.RUnlock()

	return len(wr2.childs)
}

func (wr2 *WR2) Add(child string) {
	wr2.lock.Lock()
	defer wr2.lock.Unlock()

	for _, h := range wr2.childs {
		if h == child {
			return
		}
	}
	wr2.childs = append(wr2.childs, child)
	fmt.Printf("######注册 %s 服务之后最新结果：%s\n", child, wr2.childs)
}

func (wr2 *WR2) AddWeight(child string, weight int) {
	if weight < 0 || weight > 10 {
		weight = 1
	}
	wr2.lock.Lock()
	defer wr2.lock.Unlock()

	for _, h := range wr2.childs {
		if h == child {
			return
		}
	}

	for i := 0; i < weight; i++ {
		wr2.childs = append(wr2.childs, child)
	}
}

func (wr2 *WR2) Exist(child string) bool {
	wr2.lock.Lock()
	defer wr2.lock.Unlock()

	for _, h := range wr2.childs {
		if h == child {
			return true
		}
	}

	return false
}

func (wr2 *WR2) Remove(child string) {
	wr2.lock.Lock()
	defer wr2.lock.Unlock()

	for i, h := range wr2.childs {
		if child == h {
			wr2.childs = append(wr2.childs[:i], wr2.childs[i+1:]...)
		}
	}
	fmt.Printf("#####注销 %s 服务之后最新结果：%s\n", child, wr2.childs)
}

func (wr2 *WR2) Balance() (string, error) {
	if wr2 == nil {
		return "", nil
	}
	wr2.lock.Lock()
	defer wr2.lock.Unlock()

	if len(wr2.childs) == 0 {
		return "", errors.New("######没有负载节点######")
	}

	child := wr2.childs[wr2.weight%len(wr2.childs)]
	wr2.weight++

	return child, nil
}
