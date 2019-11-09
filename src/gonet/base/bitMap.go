package base

import (
	"math"
	"unsafe"
)

const(
	size_int = int(unsafe.Sizeof(int(0))) * 8
)

type(
	BitMap struct {
		m_Bits []int
		m_Size int
	}

	IBitMap interface {
		Init(size int)
		Set(index int)//设置位
		Test(index int) bool//位是否被设置
		Clear(index int)//清楚位
		ClearAll()
	}
)

func (this *BitMap) Init(size int){
	this.m_Size = int(math.Ceil(float64(size) / float64(size_int) ))
	this.m_Bits = make([]int, this.m_Size)
}

func (this *BitMap) Set(index int){
	if index >= this.m_Size * size_int{
		return
	}

	this.m_Bits[index / size_int] |= 1 << uint(index % size_int)
}

func (this *BitMap) Test(index int) bool{
	if index >= this.m_Size * size_int{
		return false
	}

	return this.m_Bits[index / size_int] & (1 << uint(index % size_int)) != 0

}

func (this *BitMap) Clear(index int){
	if index >= this.m_Size * size_int{
		return
	}

	this.m_Bits[index / size_int] &= ^(1 << uint(index % size_int))
}

func (this *BitMap) ClearAll(){
	this.Init(this.m_Size * size_int)
}

func NewBitMap(size int) *BitMap{
	bitmap := &BitMap{}
	bitmap.Init(size)
	return bitmap
}