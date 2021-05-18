package vector

import (
	"gonet/base/containers"
	"log"
)

func assert(x bool, y string) {
	if bool(x) == false {
		log.Printf("\nFatal :{%s}", y)
	}
}

const(
	VectorBlockSize = 16
)

type (
	Vector struct{
		mElementCount int
		mArraySize int
		mArray []interface{}
	}

	IVector interface {
		containers.Container
		insert(int)
		increment()
		decrement()

		Erase(int)
		PushFront(interface{})
		PushBack(interface{})
		PopFront()
		PopBack()
		Front() interface{}
		Back() interface{}
		Len() int
		Get(int) interface{}
		Swap(i, j int)
		Less(i, j int) bool
	}
)

func (this *Vector) insert(index int) {
	assert(index <= this.mElementCount, "Vector<T>::insert - out of bounds index.")

	if this.mElementCount == this.mArraySize{
		this.resize(this.mElementCount + 1)
	}else{
		this.mElementCount++
	}

	for i := this.mElementCount - 1; i > index; i--{
		this.mArray[i] = this.mArray[i-1]
	}
}

func (this *Vector) increment() {
	if this.mElementCount == this.mArraySize{
		this.resize(this.mElementCount + 1)
	}else{
		this.mElementCount++
	}
}

func (this *Vector) decrement() {
	assert(this.mElementCount != 0, "Vector<T>::decrement - cannot decrement zero-length vector.")
	this.mElementCount--
}

func (this *Vector) resize(newCount int) bool{
	if(newCount > 0){
		blocks := newCount / VectorBlockSize
		if newCount % VectorBlockSize != 0{
			blocks++
		}

		this.mElementCount = newCount
		this.mArraySize = blocks * VectorBlockSize
		newAarray := make([]interface{}, this.mArraySize + 1)
		copy(newAarray, this.mArray)
		this.mArray = newAarray
	}
	return  true
}

func (this *Vector) Erase(index int) {
	assert(index < this.mElementCount, "Vector<T>::erase - out of bounds index.")
	if index < this.mElementCount - 1 {
		copy(this.mArray[index:this.mElementCount], this.mArray[index+1:this.mElementCount])
	}

	this.mElementCount--
}

func (this *Vector) PushFront(value interface{}) {
	this.insert(0)
	this.mArray[0] = value
}

func (this *Vector) PushBack(value interface{}) {
	this.increment()
	this.mArray[this.mElementCount-1] = value
}

func (this *Vector) PopFront() {
	assert(this.mElementCount != 0, "Vector<T>::pop_front - cannot pop the front of a zero-length vector.")
	this.Erase(0)
}

func (this *Vector) PopBack() {
	assert(this.mElementCount != 0, "Vector<T>::pop_back - cannot pop the back of a zero-length vector.")
	this.decrement()
}

// Check that the index is within bounds of the list
func (this *Vector) withinRange(index int) bool {
	return index >= 0 && index < this.mElementCount
}

func (this *Vector) Front() interface{}{
	assert(this.mElementCount != 0, "Vector<T>::first - Error, no first element of a zero sized array! (const)")
	return this.mArray[0]
}

func (this *Vector) Back() interface{}{
	assert(this.mElementCount != 0, "Vector<T>::last - Error, no last element of a zero sized array! (const)")
	return this.mArray[this.mElementCount - 1]
}

func (this *Vector) Empty() bool{
	return  this.mElementCount == 0
}

func (this *Vector) Size() int{
	return  this.mArraySize
}

func (this *Vector) Clear() {
	this.mElementCount = 0
}

func (this *Vector) Len() int{
	return this.mElementCount
}

func (this *Vector) Get(index int) interface{}{
	assert(index < this.mElementCount, "Vector<T>::operator[] - out of bounds array access!")
	return this.mArray[index]
}

func (this *Vector) Values() []interface{}{
	return this.mArray[0:this.mElementCount]
}

func (this *Vector) Swap(i, j int){
	this.mArray[i], this.mArray[j] = this.mArray[j], this.mArray[i]
}

func (this *Vector) Less(i, j int) bool {
	return true
}

func NewVector() *Vector{
	return &Vector{}
}