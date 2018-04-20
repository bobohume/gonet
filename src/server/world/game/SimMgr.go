package game

type(
	SimMgr struct {
		SimMap map [int] ISimObject
	}

	ISimMgr interface {
		GetSimObject(int)
	}
)

func (this *SimMgr) GetSimObject(simId int) ISimObject{
	pSimObj,  bExist := this.SimMap[simId]
	if bExist{
		return  pSimObj
	}
	return nil
}