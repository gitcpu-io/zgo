package limiter

/*
@Time : 2019-04-29 13:40
@Author : rubinus.chu
@File : simple
@project: zgo
*/

const (
	max = 999999999
)

type SimpleBucketer interface {
	Get(num int32) int32
	Release(num int32) int32
	Len() int32
	Cap() int32
	Resize(size int32)
}

type SimpleBucket struct {
	capacity int32
	offset   int32
	bucket   chan uint8
}

func NewBucket(cc int32) *SimpleBucket {
	if cc <= 0 {
		cc = max
	}
	return &SimpleBucket{
		capacity: cc,
		bucket:   make(chan uint8, max),
	}
}

// Get 提取指定num个token，如果不够返回真实可用的token个数
func (cl *SimpleBucket) Get(num int32) int32 {
	if num == 0 || num > max {
		return -1
	}

	cLen := cl.Len()

	if cLen >= cl.capacity {
		return 0
	}

	beef := cl.capacity - cLen
	if beef < num { //没有足够的token可以提取
		num = beef
	}
	for i := 0; int32(i) < num; i++ {
		cl.bucket <- 1
	}
	return num
}

// Release 释放指定个num的token，如果释放太多，按真实的可释放数
func (cl *SimpleBucket) Release(num int32) int32 {
	if num == 0 || num > max {
		return -1
	}

	if cl.offset > 0 {
		cl.offset--
		return -1 //resize时过滤掉偏移量
	}

	beef := cl.Len()

	if beef == 0 {
		return 0
	}

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

func (cl *SimpleBucket) Cap() int32 {
	return cl.capacity
}

func (cl *SimpleBucket) Resize(size int32) {
	if cl.capacity == size {
		return //当传进来的size和原来长度一样时，不需要重置大小
	}

	if cl.capacity > size { //如果减小
		cLen := cl.Len()
		beef := cLen - size
		if beef >= 0 { //有足够的token可以释放
			cl.Release(cl.capacity - size)
		} else if beef < 0 { //比较麻烦，表明token已全部或部分释放掉，正在回来的路上，属于中间态
			if cLen != 0 {
				cl.offset = size - cLen //还有中间态未释放掉
				cl.Release(cLen)        //释放掉部分
			} else {
				cl.offset = size
			}
		}
	}
	cl.capacity = size //变大或变小时对并发量进行重置
}
