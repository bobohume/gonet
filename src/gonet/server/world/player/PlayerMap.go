package player

func (this *Player) AddMap() {
	//actor.SendMsg("mapmgr", "LoginMap", 1, this.AccountId, this.SocketId)
}

func (this *Player) LeaveMap() {
	//actor.SendMsg("mapmgr", "LogoutMap", 1, this.AccountId)
}

func (this *Player) ReloginMap(SocketId int) {
	//SendToMap(1, "ReloginMap", this.AccountId, SocketId)
}

//--------------发送给客户端----------------------//
/*func SendToMap(mapId int64, funcName string, params  ...interface{}){
	pMap := game.MAPMGR.GetActor(mapId)
	if pMap != nil{
		pMap.SendMsg(funcName, params...)
	}
}*/