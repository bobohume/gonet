package common

type (
	BaseDataRes struct {
		DataMap map[interface{}]interface{}
	}

	IBaseDataRes interface {
		Close()
		Clear()
		Init()
		AddData(int, interface{})
		GetBaseData(int) interface{}
		Read() bool
	}
)

func (b *BaseDataRes) Close() {
	b.Clear()
}

func (b *BaseDataRes) Clear() {
	for i, _ := range b.DataMap {
		delete(b.DataMap, i)
	}
}

func (b *BaseDataRes) AddData(id int, pData interface{}) {
	b.DataMap[id] = pData
}

func (b *BaseDataRes) GetBaseData(id int) interface{} {
	pData, exist := b.DataMap[id]
	if exist {
		return pData
	}
	return nil
}

func (b *BaseDataRes) Init() {
	b.DataMap = make(map[interface{}]interface{})
}

func (b *BaseDataRes) Read() bool {
	return true
}
