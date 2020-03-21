package player

import (
	"gonet/actor"
	"gonet/message"
	"gonet/server/world"
)

func (this *Player) AddMap() {
	actor.MGR.SendMsg("mapmgr", "LoginMap", 1, this.AccountId, this.GetGateClusterId())
}

func (this *Player) LeaveMap() {
	actor.MGR.SendMsg("mapmgr", "LogoutMap", 1, this.AccountId)
}

func (this *Player) ReloginMap() {
	//SendToMap(this.AccountId, "ReloginMap", this.AccountId, this.GetGateClusterId())
}

//--------------发送给客户端----------------------//
func SendToZone(Id int64, ClusterId int, funcName string, params  ...interface{}){
	params1 := make([]interface{}, len(params)+1)
	params1[0] = &message.RpcHead{Id:Id}
	copy(params1[1:], params)
	world.SERVER.GetClusterMgr().SendMsg(ClusterId, funcName, params1...)
}