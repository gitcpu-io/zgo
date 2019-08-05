package zgolimiter

/*
@Time : 2019-04-29 13:42
@Author : rubinus.chu
@File : limiter
@project: zgo
*/

type Bucketer interface {
	NewBucket(int32) SimpleBucketer
}

type bucket struct {
}

func (b *bucket) NewBucket(c int32) SimpleBucketer {
	return NewBucket(c)
}

func New() Bucketer {
	return &bucket{}
}
