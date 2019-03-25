package base

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
		insert(int)
		increment()
		decrement()

		Erase(int)
		Push_front(interface{})
		Push_back(interface{})
		Pop_front()
		Pop_back()
		Front() interface{}
		Back() interface{}
		Begin() *interface{}
		End() *interface{}
		Next(*int) *interface{}
		First() interface{}
		Last() interface{}
		Empty() bool
		Size() int
		Clear()
		Len() int
		Get(int) interface{}
		Array() []interface{}
		Swap(i, j int)
		Less(i, j int) bool
	}
)

func (this *Vector) insert(index int) {
	Assert(index <= this.mElementCount, "Vector<T>::insert - out of bounds index.")

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
	Assert(this.mElementCount != 0, "Vector<T>::decrement - cannot decrement zero-length vector.");
	this.mElementCount--
}

func (this *Vector) resize(newCount int) bool{
	if(newCount > 0){
		blocks := newCount / VectorBlockSize;
		if newCount % VectorBlockSize != 0{
			blocks++;
		}

		this.mElementCount = newCount
		this.mArraySize = blocks * VectorBlockSize;
		this.mArray = append(this.mArray, make([]interface{}, VectorBlockSize + 1)...)
	}
	return  true;
}

func (this *Vector) Erase(index int) {
	Assert(index < this.mElementCount, "Vector<T>::erase - out of bounds index.")

	if index < this.mElementCount - 1 {
		for i := index; i < this.mElementCount - 1; i++{
			this.mArray[i] = this.mArray[i+1]
		}
	}

	this.mElementCount--
}

func (this *Vector) Push_front(value interface{}) {
	this.insert(0)
	this.mArray[0] = value
}

func (this *Vector) Push_back(value interface{}) {
	this.increment()
	this.mArray[this.mElementCount-1] = value
}

func (this *Vector) Pop_front() {
	Assert(this.mElementCount != 0, "Vector<T>::pop_front - cannot pop the front of a zero-length vector.")
	this.Erase(0);
}

func (this *Vector) Pop_back() {
	Assert(this.mElementCount != 0, "Vector<T>::pop_back - cannot pop the back of a zero-length vector.")
	this.decrement()
}

func (this *Vector) Front() interface{}{
	return this.First()
}

func (this *Vector) Back() interface{}{
	return this.Last()
}

func (this *Vector) Begin() *interface{}{
	if this.Empty(){
		return nil
	}else{
		return &this.mArray[0]
	}
}

func (this *Vector) End() *interface{}{
	if this.Empty() {
		return nil
	}else{
		return &this.mArray[this.mElementCount]
	}
}

func (this *Vector) Next(index *int) *interface{}{
	if *index < this.mElementCount - 1 {
		*index++
		return &this.mArray[*index]
	}

	return this.End()
}

func (this *Vector) First() interface{}{
	Assert(this.mElementCount != 0, "Vector<T>::first - Error, no first element of a zero sized array! (const)")
	return this.mArray[0]
}

func (this *Vector) Last() interface{}{
	Assert(this.mElementCount != 0, "Vector<T>::last - Error, no last element of a zero sized array! (const)")
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
	Assert(index < this.mElementCount, "Vector<T>::operator[] - out of bounds array access!");
	return this.mArray[index];
}

func (this *Vector) Array() []interface{}{
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