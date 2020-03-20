package player

import (
	"gonet/actor"
	"gonet/message"
)

func (this *Player) AddMap() {
	actor.MGR.SendMsg("mapmgr", "LoginMap", 1, this.AccountId, this.GetGateSocketId())
}

func (this *Player) LeaveMap() {
	actor.MGR.SendMsg("mapmgr", "LogoutMap", 1, this.AccountId)
}

func (this *Player) ReloginMap(SocketId int) {
	SendToMap(this.AccountId, "ReloginMap", this.AccountId, SocketId)
}

//--------------发送给客户端----------------------//
func SendToMap(Id int64, funcName string, params  ...interface{}){
	params1 := make([]interface{}, len(params)+1)
	params1[0] = &message.RpcHead{Id:Id}
	copy(params1[1:], params)
	actor.MGR.SendMsg("mapmgr", funcName, params1...)
}