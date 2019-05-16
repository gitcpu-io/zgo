package limiter

/*
@Time : 2019-04-29 13:40
@Author : rubinus.chu
@File : simple
@project: zgo
*/

type SimpleBucketer interface {
	GetToken() int
	ReleaseToken()
}

type SimpleBucket struct {
	concurrent int
	bucket     chan int
}

func NewSimpleBucket(cc int) *SimpleBucket {
	return &SimpleBucket{
		concurrent: cc,
		bucket:     make(chan int, cc),
	}
}

func (cl *SimpleBucket) GetToken() int {
	if len(cl.bucket) >= cl.concurrent {
		return 0
	}
	cl.bucket <- 1
	return len(cl.bucket)
}

func (cl *SimpleBucket) ReleaseToken() {
	<-cl.bucket
}
