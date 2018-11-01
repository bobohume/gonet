package common

type(
	BaseDataRes struct{
		DataMap map[interface{}] interface{}
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

func (this *BaseDataRes) Close(){
	this.Clear()
}

func (this *BaseDataRes) Clear(){
	for i,_ := range this.DataMap{
		delete(this.DataMap, i)
	}
}

func (this *BaseDataRes) AddData(id int, pData interface{}){
	this.DataMap[id] = pData
}

func (this *BaseDataRes) GetBaseData(id int) interface{}{
	pData, exist := this.DataMap[id]
	if exist{
		return pData
	}
	return nil
}

func (this *BaseDataRes) Init(){
	this.DataMap = make(map[interface{}] interface{})
}

func (this *BaseDataRes) Read() bool{
	return true
}