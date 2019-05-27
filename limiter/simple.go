package limiter

/*
@Time : 2019-04-29 13:40
@Author : rubinus.chu
@File : simple
@project: zgo
*/

const (
	max = 9999
)

type SimpleBucketer interface {
	Get(num int32) int32
	Release(num int32) int32
	Len() int32
	Cap() int32
	BeLeft() int32
	Resize(size, offset int32)
	Clear() int32
}

type SimpleBucket struct {
	capacity int32
	offset   int32
	bucket   chan uint8
}

func NewBucket(cc int32) *SimpleBucket {
	if cc <= 0 || cc > max {
		cc = max
	}
	return &SimpleBucket{
		capacity: cc,
		bucket:   make(chan uint8, max),
		offset:   0,
	}
}

// Get 提取指定num个token，如果不够返回真实可用的token个数
func (cl *SimpleBucket) Get(num int32) int32 {
	if num == 0 || num > max {
		return -1
	}

	if cl.offset != 0 {
		//fmt.Println(cl.offset, "-------offset有值，返回不了")
		return 0
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

	if cl.offset != 0 {
		cl.offset -= num
		//fmt.Println(cl.offset, "=========cl.offset==========")
		return 0
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

// Clear 清空整个桶
func (cl *SimpleBucket) Clear() int32 {
	cl.offset = 0
	cLen := cl.Len()
	for i := 0; int32(i) < cLen; i++ {
		<-cl.bucket
	}
	return cLen
}

// Len 已经使用的长度
func (cl *SimpleBucket) Len() int32 {
	return int32(len(cl.bucket))
}

// Cap 整个桶的容量
func (cl *SimpleBucket) Cap() int32 {
	return cl.capacity
}

// BeLeft 剩余可以使用的
func (cl *SimpleBucket) BeLeft() int32 {
	return cl.Cap() - cl.Len()
}

func (cl *SimpleBucket) Resize(size int32, changeSize int32) {
	if cl.capacity == size || size == 0 {
		return //当传进来的size和原来长度一样时，不需要重置大小
	}

	oldCapacity := cl.capacity

	cl.capacity = size

	if oldCapacity > size { //如果减小
		if cl.Len() >= changeSize {
			//bucket中有足够的

			cl.offset = 0

			cl.Release(changeSize)

		} else {
			offset := changeSize - cl.Len()
			//释放掉部分
			cl.Release(cl.Len())

			cl.offset = offset

		}
	}
}

func (cl *SimpleBucket) ResizeBack(size int32) {
	if cl.capacity == size || size == 0 {
		return //当传进来的size和原来长度一样时，不需要重置大小
	}

	oldCapacity := cl.capacity

	cl.capacity = size

	if oldCapacity > size { //如果减小
		offset := oldCapacity - size

		cLen := cl.Len()

		if cLen-offset >= 0 { //桶中有足够的token可以释放
			//fmt.Println("--有足够的token可以释放--",cLen, offset,cl.capacity)

			cl.Release(offset)

			cl.offset = offset

		} else { //比较麻烦，表明token已全部或部分释放掉，正在回来的路上，属于中间态

			if cLen != 0 {

				nw := cl.Release(cLen) //释放掉部分

				cl.offset = offset - nw

				//fmt.Println("释放已有token:",nw, cl.offset,"----------------000000000===",offset, cLen,oldCapacity, cl.capacity, size)

			} else {
				cl.offset = offset

				//fmt.Println(cl.offset,"----------------111111111===",offset, cLen,oldCapacity, cl.capacity, size)

			}

		}

	}
}
