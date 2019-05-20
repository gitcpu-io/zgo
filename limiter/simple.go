package limiter

/*
@Time : 2019-04-29 13:40
@Author : rubinus.chu
@File : simple
@project: zgo
*/

type SimpleBucketer interface {
	GetToken() int32
	ReleaseToken()
	Len() int32
	Reset(size int32)
}

type SimpleBucket struct {
	concurrent int32
	bucket     chan uint8
}

func NewSimpleBucket(cc int32) *SimpleBucket {
	return &SimpleBucket{
		concurrent: cc,
		bucket:     make(chan uint8, cc),
	}
}

func (cl *SimpleBucket) GetToken() int32 {
	if cl.Len() >= cl.concurrent {
		return 0
	}
	cl.bucket <- 1
	return cl.Len()
}

func (cl *SimpleBucket) ReleaseToken() {
	<-cl.bucket
}

func (cl *SimpleBucket) Len() int32 {
	return int32(len(cl.bucket))
}

func (cl *SimpleBucket) Reset(size int32) {
	if cl.Len() == size {
		return //当传进来的size和bucket长度一样时，不需要重置大小
	}
	if cl.Len() > size { //如果减小，就从循环从chan中读取offset
		for i := 0; i < int(cl.Len()-size); i++ {
			cl.ReleaseToken()
		}
	}
	cl.concurrent = size //变大或变小时对并发量进行重置
}
