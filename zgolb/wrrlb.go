package zgolb

import (
	"errors"
	"sync"
)

/*
@Time : 2019-07-30 21:32
@Author : rubinus.chu
@File : wrrlb
@project: zgo
*/

type WR2er interface {
	Add(child string)
	AddWeight(child string, weight int)
	Exist(child string) bool
	Remove(child string)
	Balance() (string, error)
}

type WR2 struct {
	weight int
	childs []string
	sync.Mutex
}

func NewWR2(childs ...string) *WR2 {
	return &WR2{weight: 1, childs: childs}
}

func (wr2 *WR2) Add(child string) {
	wr2.Lock()
	defer wr2.Unlock()

	for _, h := range wr2.childs {
		if h == child {
			return
		}
	}
	wr2.childs = append(wr2.childs, child)
}

func (wr2 *WR2) AddWeight(child string, weight int) {
	if weight < 0 || weight > 10 {
		weight = 1
	}
	wr2.Lock()
	defer wr2.Unlock()

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
	wr2.Lock()
	defer wr2.Unlock()

	for _, h := range wr2.childs {
		if h == child {
			return true
		}
	}

	return false
}

func (wr2 *WR2) Remove(child string) {
	wr2.Lock()
	defer wr2.Unlock()

	for i, h := range wr2.childs {
		if child == h {
			wr2.childs = append(wr2.childs[:i], wr2.childs[i+1:]...)
		}
	}
}

func (wr2 *WR2) Balance() (string, error) {
	wr2.Lock()
	defer wr2.Unlock()

	if len(wr2.childs) == 0 {
		return "", errors.New("no child")
	}

	child := wr2.childs[wr2.weight%len(wr2.childs)]
	wr2.weight++

	return child, nil
}
