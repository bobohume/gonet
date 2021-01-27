package player

import (
	"gonet/server/world"
)

func (this *Player) SendToZone(funcName string, params  ...interface{}) {
	world.SendToZone(this.AccountId, this.GetZoneClusterId(), funcName, params...)
}

func (this *Player) AddMap() {
	this.SendToZone("LoginMap", 1, this.AccountId, this.GetGateClusterId(), world.SERVER.GetClusterMgr().Id())
}

func (this *Player) LeaveMap() {
	this.SendToZone("LogoutMap", 1, this.AccountId)
}

func (this *Player) ReloginMap() {
	this.SendToZone("ReloginMap", this.AccountId, this.GetGateClusterId())
}

//添加buff
func (this *Player) AddBuff(Orgint int, BuffId int) {
	if BuffId < 0{
		return
	}
	this.SendToZone("AddBuff", this.AccountId, Orgint, BuffId)
}

//删除buff
func (this *Player) RemoveBuff(BuffId int) {
	if BuffId < 0{
		return
	}
	this.SendToZone("RemoveBuff", this.AccountId, BuffId)
}

//批量添加buff
func (this *Player) AddBuffS(Orgint int, BuffId []int) {
	BuffIds :=  []int{}
	for i := 0; i < len(BuffId); i++{
		if BuffId[i] < 0{
			continue
		}
		BuffIds = append(BuffIds, int(BuffId[i]))
	}

	this.SendToZone("AddBuffS", this.AccountId, Orgint, BuffIds)
}

//批量删除buff
func (this *Player) RemoveBuffS(BuffId []int) {
	BuffIds :=  []int{}
	for i := 0; i < len(BuffId); i++{
		if BuffId[i] < 0{
			continue
		}
		BuffIds = append(BuffIds, int(BuffId[i]))
	}
	this.SendToZone("RemoveBuffS", this.AccountId, BuffIds)
}