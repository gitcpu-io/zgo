package limiter

/*
@Time : 2019-04-29 13:40
@Author : rubinus.chu
@File : simple
@project: zgo
*/

type SimpleBucketer interface {
	GetToken(num int32) int32
	ReleaseToken(num int32) int32
	Len() int32
	Capacity() int32
	Reset(size int32)
}

type SimpleBucket struct {
	concurrent int32
	bucket     chan uint8
}

func NewSimpleBucket(cc int32) *SimpleBucket {
	return &SimpleBucket{
		concurrent: cc,
		bucket:     make(chan uint8, 999999999),
	}
}

// GetToken 提取指定num个token，如果不够返回真实可用的token个数
func (cl *SimpleBucket) GetToken(num int32) int32 {

	cLen := cl.Len()

	if cLen >= cl.concurrent {
		return 0
	}

	beef := cl.concurrent - cLen
	if beef < num { //没有足够的token可以提取
		num = beef
	}
	for i := 0; int32(i) < num; i++ {
		cl.bucket <- 1
	}
	return num
}

// ReleaseToken 释放指定个num的token，如果释放太多，按真实的可释放数
func (cl *SimpleBucket) ReleaseToken(num int32) int32 {
	beef := cl.Len()
	if beef < num {
		num = beef
	}
	for i := 0; int32(i) < num; i++ {
		<-cl.bucket
	}
	return num
}

func (cl *SimpleBucket) Len() int32 {
	return int32(len(cl.bucket))
}

func (cl *SimpleBucket) Capacity() int32 {
	return cl.concurrent
}

func (cl *SimpleBucket) Reset(size int32) {
	if cl.concurrent == size {
		return //当传进来的size和原来长度一样时，不需要重置大小
	}

	if cl.concurrent > size { //如果减小，就从循环从chan中读取offset

		beef := cl.Len() - size
		if beef > 0 {
			cl.ReleaseToken(beef)
		}
	}
	cl.concurrent = size //变大或变小时对并发量进行重置
}
