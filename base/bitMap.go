package base

import (
	"math"
	"unsafe"
)

const (
	size_int = int(unsafe.Sizeof(int(0))) * 8
)

type (
	BitMap struct {
		bits []int
		size int
	}

	IBitMap interface {
		Init(size int)
		Set(index int)       //设置位
		Test(index int) bool //位是否被设置
		Clear(index int)     //清楚位
		ClearAll()
	}
)

func (b *BitMap) Init(size int) {
	b.size = int(math.Ceil(float64(size) / float64(size_int)))
	b.bits = make([]int, b.size)
}

func (b *BitMap) Set(index int) {
	if index >= b.size*size_int {
		return
	}

	b.bits[index/size_int] |= 1 << uint(index%size_int)
}

func (b *BitMap) Test(index int) bool {
	if index >= b.size*size_int {
		return false
	}

	return b.bits[index/size_int]&(1<<uint(index%size_int)) != 0

}

func (b *BitMap) Clear(index int) {
	if index >= b.size*size_int {
		return
	}

	b.bits[index/size_int] &= ^(1 << uint(index%size_int))
}

func (b *BitMap) ClearAll() {
	b.Init(b.size * size_int)
}

func NewBitMap(size int) *BitMap {
	bitmap := &BitMap{}
	bitmap.Init(size)
	return bitmap
}
