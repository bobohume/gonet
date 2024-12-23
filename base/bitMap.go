package base

type (
	BitMap[V ~int | ~uint] struct {
		bits map[int]V
	}

	IBitMap[V ~int | ~uint] interface {
		Init()
		Set(index int, flag bool) //设置位
		Get(index int) bool       //位是否被设置
		ClearAll()
	}
)

func (b *BitMap[V]) Init() {
	b.bits = make(map[int]V)
}

func (b *BitMap[V]) Set(index int, flag bool) {
	if flag {
		b.bits[index/size_int] |= 1 << V(index%size_int)
	} else {
		b.bits[index/size_int] &= ^(1 << V(index%size_int))
	}
}

func (b *BitMap[V]) Get(index int) bool {
	return b.bits[index/size_int]&(1<<V(index%size_int)) != 0
}

func (b *BitMap[V]) ClearAll() {
	b.Init()
}
