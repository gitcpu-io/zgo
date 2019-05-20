package limiter

/*
@Time : 2019-04-29 13:42
@Author : rubinus.chu
@File : limiter
@project: zgo
*/

type Bucketer interface {
	NewSimpleBucket(int32) SimpleBucketer
}

type bucket struct {
}

func (b *bucket) NewSimpleBucket(c int32) SimpleBucketer {
	return NewSimpleBucket(c)
}

func New() Bucketer {
	return &bucket{}
}
