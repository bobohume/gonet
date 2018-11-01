package account

import (
	"actor"
	"base"
	"message"
	"server/common"
	"sync"
)

type (
	HashSocketMap map[int] *common.ServerInfo

	ServerSocketManager struct{
		actor.Actor
		m_Locker	*sync.RWMutex
		m_SocketMap HashSocketMap
		m_GateMap 	HashSocketMap
	}

	IServerSocketManager interface {
		actor.IActor
		AddServerMap(*common.ServerInfo)
		ReleaseServerMap(int,bool)
		GetAllGate(base.IVector)
		GetSeverInfo(int, int) *common.ServerInfo
	}
)

func (this *ServerSocketManager) AddServerMap(pSeverInfo *common.ServerInfo){
	this.m_Locker.Lock()
	this.m_SocketMap[pSeverInfo.SocketId] = pSeverInfo
	this.m_Locker.Unlock()

	switch pSeverInfo.Type {
	case int(message.SERVICE_GATESERVER):
		this.m_Locker.Lock()
		this.m_GateMap[pSeverInfo.SocketId] = pSeverInfo
		this.m_Locker.Unlock()
		SERVER.GetLog().Printf("ADD GATE SERVER: [%d]-[%s:%d]", pSeverInfo.SocketId, pSeverInfo.Ip, pSeverInfo.Port)
	}
}

func (this *ServerSocketManager) ReleaseServerMap(socketid int, bClose bool){
	this.m_Locker.RLock()
	pServerInfo, exist := this.m_SocketMap[socketid]
	this.m_Locker.RUnlock()
	if !exist{
		return
	}

	SERVER.GetLog().Printf("服务器断开连接socketid[%d]",socketid)
	switch pServerInfo.Type {
	case int(message.SERVICE_GATESERVER):
		SERVER.GetLog().Printf("与Gate服务器断开连接,id[%d]",pServerInfo.SocketId)
		this.m_Locker.Lock()
		delete(this.m_GateMap, pServerInfo.SocketId)
		this.m_Locker.Unlock()
	}

	this.m_Locker.Lock()
	delete(this.m_SocketMap, pServerInfo.SocketId)
	this.m_Locker.Unlock()
	if bClose {
		SERVER.GetServer().StopClient(socketid)
	}
}

func (this *ServerSocketManager) GetSeverInfo(nType int, socketId int) *common.ServerInfo{
	switch nType {
	case int(message.SERVICE_GATESERVER):
		this.m_Locker.RLock()
		pServerInfo, _ := this.m_GateMap[socketId]
		this.m_Locker.RUnlock()
		return  pServerInfo
	}
	return  nil
}

func (this *ServerSocketManager) GetAllGate(vec base.IVector){
	this.m_Locker.RLock()
	for _, v := range this.m_GateMap{
		vec.Push_back(v)
	}
	this.m_Locker.RUnlock()
}

func (this *ServerSocketManager) Init(num int){
	this.Actor.Init(num)
	this.m_GateMap 		= make(HashSocketMap)
	this.m_SocketMap 	= make(HashSocketMap)
	this.m_Locker		= &sync.RWMutex{}
	this.RegisterCall("DISCONNECT", func(socketid int) {
		this.ReleaseServerMap(socketid, false)
	})

	this.RegisterCall("CONNECT", func(nType int, Ip string, Port int) {
		pServerInfo := &common.ServerInfo{}
		pServerInfo.SocketId = this.GetSocketId()
		pServerInfo.Type = nType
		pServerInfo.Ip = Ip
		pServerInfo.Port = Port

		this.AddServerMap(pServerInfo)
	})

	this.RegisterCall("G_ClientLost", func(accountId int, socketId int) {
		vec := base.NewVector()
		this.GetAllGate(vec)
		for _, v := range vec.Array(){
			pSeverInfo := v.(*common.ServerInfo)
			if pSeverInfo.SocketId != socketId{
				SERVER.GetServer().SendMsgByID(pSeverInfo.SocketId, "G_ClientLost", accountId)
			}
		}
	})
	this.Actor.Start()
}
