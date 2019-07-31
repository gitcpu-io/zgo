package zgolb

/*
@Time : 2019-07-30 19:40
@Author : rubinus.chu
@File : lb
@project: zgo
*/

type Lber interface {
	WR2(childs ...string) WR2er
}

type lb struct {
}

func (b *lb) WR2(childs ...string) WR2er {
	return NewWR2(childs...)
}

func NewLB() Lber {
	return &lb{}
}
