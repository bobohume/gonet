package player

import (
	"gonet/actor"
)

func (this *Player) AddMap() {
	actor.MGR.SendMsg("mapmgr", "LoginMap", 1, this.AccountId, this.SocketId)
}

func (this *Player) LeaveMap() {
	actor.MGR.SendMsg("mapmgr", "LogoutMap", 1, this.AccountId)
}

func (this *Player) ReloginMap(SocketId int) {
	SendToMap(this.AccountId, "ReloginMap", this.AccountId, SocketId)
}

//--------------发送给客户端----------------------//
func SendToMap(Id int64, funcName string, params  ...interface{}){
	actor.MGR.SendMsgById("mapmgr", Id, funcName, params...)
}