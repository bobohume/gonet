package common

type(
	CBaseDataRes struct{
		DataMap map[interface{}] interface{}
	}

	IBaseDataRes interface {
		Close()
		Clear()
		Init()
		AddData(int, interface{})
		GetData(int) interface{}
		Read() bool
	}
)

func (this *CBaseDataRes) Close(){
	this.Clear()
}

func (this *CBaseDataRes) Clear(){
	for i,_ := range this.DataMap{
		delete(this.DataMap, i)
	}
}

func (this *CBaseDataRes) AddData(id interface{}, pData interface{}){
	this.DataMap[id] = pData
}

func (this *CBaseDataRes) GetData(id interface{}) interface{}{
	pData, exist := this.DataMap[id]
	if exist{
		return pData
	}
	return nil
}

func (this *CBaseDataRes) Init(){
	this.DataMap = make(map[interface{}] interface{})
}

func (this *CBaseDataRes) Read() bool{
	return true
}