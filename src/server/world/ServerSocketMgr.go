package world

import (
	"actor"
	"container/list"
	"message"
)

type (
	HashSocketMap map[int] *stServerInfo

	CServerSocketManager struct{
		actor.Actor
		m_SocketMap HashSocketMap
		m_GateMap 	HashSocketMap
	}

	IServerSocketManager interface {
		actor.IActor
		AddServerMap(stServerInfo)
		ReleaseServerMap(int,bool)
		GetAllGate(list.List)
		GetGateId(int)int
		GetSeverInfo(int, int) *stServerInfo
	}

	stServerInfo struct {
		Type int//服务类型编号
		GateId int//服务网关编号
		Ip string//服务IP
		Port int//服务端口
		SocketId int//连接句柄
	}
)

func (this *CServerSocketManager) AddServerMap(info stServerInfo){
	delete(this.m_SocketMap, info.SocketId)

	pServerInfo := new(stServerInfo)
	*pServerInfo = info
	this.m_SocketMap[info.SocketId] = pServerInfo

	switch info.Type {
	case int(message.SERVICE_GATESERVER):
		this.m_GateMap[info.GateId] = pServerInfo
		SERVER.GetLog().Printf("ADD GATE SERVER: [%d]-[%s:%d]", pServerInfo.GateId, pServerInfo.Ip, pServerInfo.Port)
	}
}

func (this *CServerSocketManager) ReleaseServerMap(socketid int, bClose bool){
	pServerInfo, exist := this.m_SocketMap[socketid]
	if !exist{
		return
	}

	SERVER.GetLog().Printf("服务器断开连接socketid[%d]",socketid)
	switch pServerInfo.Type {
	case int(message.SERVICE_GATESERVER):
		SERVER.GetLog().Printf("与Gate服务器断开连接,id[%d]",pServerInfo.GateId)
		delete(this.m_SocketMap, pServerInfo.GateId)
	}

	delete(this.m_SocketMap, pServerInfo.SocketId)
	if bClose {
		SERVER.GetServer().StopClient(socketid)
	}
}

func (this *CServerSocketManager) GetSeverInfo(nType int, gateid int) *stServerInfo{
	switch nType {
	case int(message.SERVICE_GATESERVER):
		pServerInfo, _ := this.m_GateMap[gateid]
		return  pServerInfo
	}
	return  nil
}

func (this *CServerSocketManager) GetGateId(socketid int)int{
	pServerInfo, exist := this.m_SocketMap[socketid]
	if !exist{
		return -1
	}
	return  pServerInfo.GateId
}

func (this *CServerSocketManager) GetAllGate(gates list.List){
	for _, v := range this.m_GateMap{
		gates.PushBack(v)
	}
}


func (this *CServerSocketManager) Init(num int){
	this.Actor.Init(num)
	this.m_GateMap 		= make(HashSocketMap)
	this.m_SocketMap 	= make(HashSocketMap)
	this.RegisterCall("DISCONNECT", func(caller *actor.Caller, socketid int) {
		this.ReleaseServerMap(socketid, false)
	})

	this.RegisterCall("CONNECT", func(caller *actor.Caller, nType int, GateId int, Ip string, Port int) {
		var ServerInfo stServerInfo
		ServerInfo.SocketId = caller.SocketId
		ServerInfo.Type = nType
		ServerInfo.GateId = GateId
		ServerInfo.Ip = Ip
		ServerInfo.Port = Port

		this.AddServerMap(ServerInfo)
	})
	this.Actor.Start()
}
