package player

import (
	"gonet/common/cluster"
	"gonet/server/game"
)

func (p *Player) SendToZone(funcName string, params ...interface{}) {
	game.SendToZone(p.PlayerId, p.ZoneClusterId, funcName, params...)
}

func (p *Player) AddMap() {
	p.SendToZone("MapMgr.LoginMap", 200000, p.PlayerId, p.GetGateClusterId(), cluster.MGR.Id())
}

func (p *Player) LeaveMap() {
	p.SendToZone("MapMgr.LogoutMap", 200000, p.PlayerId)
}

func (p *Player) ReloginMap() {
	p.SendToZone("MapMgr.ReloginMap", p.PlayerId, p.GetGateClusterId())
}

//添加buff
func (p *Player) AddBuff(Orgint int, BuffId int) {
	if BuffId < 0 {
		return
	}
	p.SendToZone("Map.AddBuff", p.PlayerId, Orgint, BuffId)
}

//删除buff
func (p *Player) RemoveBuff(BuffId int) {
	if BuffId < 0 {
		return
	}
	p.SendToZone("Map.RemoveBuff", p.PlayerId, BuffId)
}

//批量添加buff
func (p *Player) AddBuffS(Orgint int, BuffId []int) {
	BuffIds := []int{}
	for i := 0; i < len(BuffId); i++ {
		if BuffId[i] < 0 {
			continue
		}
		BuffIds = append(BuffIds, int(BuffId[i]))
	}

	p.SendToZone("Map.AddBuffS", p.PlayerId, Orgint, BuffIds)
}

//批量删除buff
func (p *Player) RemoveBuffS(BuffId []int) {
	BuffIds := []int{}
	for i := 0; i < len(BuffId); i++ {
		if BuffId[i] < 0 {
			continue
		}
		BuffIds = append(BuffIds, int(BuffId[i]))
	}
	p.SendToZone("Map.RemoveBuffS", p.PlayerId, BuffIds)
}
