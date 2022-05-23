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

const (
	VectorBlockSize = 16
)

type (
	Vector struct {
		elementCount int
		arraySize    int
		array        []interface{}
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

func (v *Vector) insert(index int) {
	assert(index <= v.elementCount, "Vector<T>::insert - out of bounds index.")

	if v.elementCount == v.arraySize {
		v.resize(v.elementCount + 1)
	} else {
		v.elementCount++
	}

	for i := v.elementCount - 1; i > index; i-- {
		v.array[i] = v.array[i-1]
	}
}

func (v *Vector) increment() {
	if v.elementCount == v.arraySize {
		v.resize(v.elementCount + 1)
	} else {
		v.elementCount++
	}
}

func (v *Vector) decrement() {
	assert(v.elementCount != 0, "Vector<T>::decrement - cannot decrement zero-length vector.")
	v.elementCount--
}

func (v *Vector) resize(newCount int) bool {
	if newCount > 0 {
		blocks := newCount / VectorBlockSize
		if newCount%VectorBlockSize != 0 {
			blocks++
		}

		v.elementCount = newCount
		v.arraySize = blocks * VectorBlockSize
		newAarray := make([]interface{}, v.arraySize+1)
		copy(newAarray, v.array)
		v.array = newAarray
	}
	return true
}

func (v *Vector) Erase(index int) {
	assert(index < v.elementCount, "Vector<T>::erase - out of bounds index.")
	if index < v.elementCount-1 {
		copy(v.array[index:v.elementCount], v.array[index+1:v.elementCount])
	}

	v.elementCount--
}

func (v *Vector) PushFront(value interface{}) {
	v.insert(0)
	v.array[0] = value
}

func (v *Vector) PushBack(value interface{}) {
	v.increment()
	v.array[v.elementCount-1] = value
}

func (v *Vector) PopFront() {
	assert(v.elementCount != 0, "Vector<T>::pop_front - cannot pop the front of a zero-length vector.")
	v.Erase(0)
}

func (v *Vector) PopBack() {
	assert(v.elementCount != 0, "Vector<T>::pop_back - cannot pop the back of a zero-length vector.")
	v.decrement()
}

// Check that the index is within bounds of the list
func (v *Vector) withinRange(index int) bool {
	return index >= 0 && index < v.elementCount
}

func (v *Vector) Front() interface{} {
	assert(v.elementCount != 0, "Vector<T>::first - Error, no first element of a zero sized array! (const)")
	return v.array[0]
}

func (v *Vector) Back() interface{} {
	assert(v.elementCount != 0, "Vector<T>::last - Error, no last element of a zero sized array! (const)")
	return v.array[v.elementCount-1]
}

func (v *Vector) Empty() bool {
	return v.elementCount == 0
}

func (v *Vector) Size() int {
	return v.arraySize
}

func (v *Vector) Clear() {
	v.elementCount = 0
}

func (v *Vector) Len() int {
	return v.elementCount
}

func (v *Vector) Get(index int) interface{} {
	assert(index < v.elementCount, "Vector<T>::operator[] - out of bounds array access!")
	return v.array[index]
}

func (v *Vector) Values() []interface{} {
	return v.array[0:v.elementCount]
}

func (v *Vector) Swap(i, j int) {
	v.array[i], v.array[j] = v.array[j], v.array[i]
}

func (v *Vector) Less(i, j int) bool {
	return true
}

func NewVector() *Vector {
	return &Vector{}
}
