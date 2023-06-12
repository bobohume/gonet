package cm

//********************************************************
//  DataMgr
//********************************************************
type (
	DataMgr[T any] struct {
		DataMap map[int64]*T
		Nil     *T
	}
)

func (a *DataMgr[T]) Init() {
	a.DataMap = make(map[int64]*T)
}

func (a *DataMgr[T]) AddData(Id int64, ac *T) {
	a.DataMap[Id] = ac
}

func (a *DataMgr[T]) DelData(Id int64) {
	delete(a.DataMap, Id)
}

func (a *DataMgr[T]) GetData(Id int64) *T {
	ac, bEx := a.DataMap[Id]
	if bEx {
		return ac
	}
	return a.Nil
}
