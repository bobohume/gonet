package account

import (
	"gonet/actor"
	"gonet/base"
	"gonet/message"
	"gonet/server/common"
	"sync"
)

type (
	HashSocketMap map[int] *common.ServerInfo

	ServerSocketManager struct{
		actor.Actor
		m_SocketMap 	HashSocketMap
		m_SocketLocker	*sync.RWMutex
		m_GateMap 		HashSocketMap
		m_GateLocker	*sync.RWMutex
		m_WorldMap		HashSocketMap
		m_WorldLocker	*sync.RWMutex
	}

	IServerSocketManager interface {
		actor.IActor
		AddServerMap(*common.ServerInfo)
		ReleaseServerMap(int, bool)
		GetAllGate(base.IVector)
		GetAllWorld(base.IVector)
		GetSeverInfo(int, int) *common.ServerInfo
		KickWorldPlayer(accountId int64)
	}
)

func (this *ServerSocketManager) AddServerMap(pServerInfo *common.ServerInfo){
	this.m_SocketLocker.Lock()
	this.m_SocketMap[pServerInfo.SocketId] = pServerInfo
	this.m_SocketLocker.Unlock()

	switch pServerInfo.Type {
	case int(message.SERVICE_GATESERVER):
		this.m_GateLocker.Lock()
		this.m_GateMap[pServerInfo.SocketId] = pServerInfo
		this.m_GateLocker.Unlock()
		SERVER.GetLog().Printf("ADD GATE SERVER: [%d]-[%s:%d]", pServerInfo.SocketId, pServerInfo.Ip, pServerInfo.Port)
		aIp := []string{}
		aPort := []int{}
		aSocket := []int{}
		vec := base.NewVector()
		this.GetAllWorld(vec)
		for _, v := range vec.Array(){
			pInfo, bOk := v.(*common.ServerInfo)
			if bOk && pInfo != nil{
				aIp = append(aIp, pInfo.Ip)
				aPort = append(aPort, pInfo.Port)
				aSocket = append(aSocket, pInfo.SocketId)
			}
		}
		SERVER.GetServer().SendMsgByID(pServerInfo.SocketId, "Dispatch_Socket_Init", aIp, aPort, aSocket)
	case int(message.SERVICE_WORLDSERVER):
		this.m_WorldLocker.Lock()
		this.m_WorldMap[pServerInfo.SocketId] = pServerInfo
		this.m_WorldLocker.Unlock()
		SERVER.GetLog().Printf("ADD WORLD SERVER: [%d]-[%s:%d]", pServerInfo.SocketId, pServerInfo.Ip, pServerInfo.Port)
		vec := base.NewVector()
		this.GetAllGate(vec)
		for _, v := range vec.Array(){
			pInfo, bOk := v.(*common.ServerInfo)
			if bOk && pInfo != nil{
				SERVER.GetServer().SendMsgByID(pInfo.SocketId, "Dispatch_Socket_Add", pServerInfo.Ip, pServerInfo.Port, pServerInfo.SocketId)
			}
		}
	}
}

func (this *ServerSocketManager) ReleaseServerMap(socketid int, bClose bool){
	this.m_SocketLocker.RLock()
	pServerInfo, exist := this.m_SocketMap[socketid]
	this.m_SocketLocker.RUnlock()
	if !exist{
		return
	}

	SERVER.GetLog().Printf("服务器断开连接socketid[%d]",socketid)
	switch pServerInfo.Type {
	case int(message.SERVICE_GATESERVER):
		SERVER.GetLog().Printf("与Gate服务器断开连接,id[%d]",pServerInfo.SocketId)
		this.m_GateLocker.Lock()
		delete(this.m_GateMap, pServerInfo.SocketId)
		this.m_GateLocker.Unlock()
	case int(message.SERVICE_WORLDSERVER):
		SERVER.GetLog().Printf("与World服务器断开连接,id[%d]",pServerInfo.SocketId)
		this.m_WorldLocker.Lock()
		delete(this.m_WorldMap, pServerInfo.SocketId)
		this.m_WorldLocker.Unlock()
		vec := base.NewVector()
		this.GetAllGate(vec)
		for _, v := range vec.Array(){
			pInfo, bOk := v.(*common.ServerInfo)
			if bOk && pInfo != nil {
				SERVER.GetServer().SendMsgByID(pInfo.SocketId, "Dispatch_Socket_Del", pServerInfo.Ip, pServerInfo.Port, pServerInfo.SocketId)
			}
		}
	}

	this.m_SocketLocker.Lock()
	delete(this.m_SocketMap, pServerInfo.SocketId)
	this.m_SocketLocker.Unlock()
	if bClose {
		SERVER.GetServer().StopClient(socketid)
	}
}

func (this *ServerSocketManager) GetSeverInfo(nType int, socketId int) *common.ServerInfo{
	switch nType {
	case int(message.SERVICE_GATESERVER):
		this.m_GateLocker.RLock()
		pServerInfo, _ := this.m_GateMap[socketId]
		this.m_GateLocker.RUnlock()
		return  pServerInfo
	case int(message.SERVICE_WORLDSERVER):
		this.m_WorldLocker.RLock()
		pServerInfo, _ := this.m_WorldMap[socketId]
		this.m_WorldLocker.RUnlock()
		return  pServerInfo
	}
	return  nil
}

func (this *ServerSocketManager) GetAllGate(vec base.IVector){
	this.m_GateLocker.RLock()
	for _, v := range this.m_GateMap{
		vec.Push_back(v)
	}
	this.m_GateLocker.RUnlock()
}

func (this *ServerSocketManager) GetAllWorld(vec base.IVector){
	this.m_WorldLocker.RLock()
	for _, v := range this.m_WorldMap{
		vec.Push_back(v)
	}
	this.m_WorldLocker.RUnlock()
}

func (this *ServerSocketManager) KickWorldPlayer(accountId int64){
	vec := base.NewVector()
	this.GetAllWorld(vec)
	for _, v := range vec.Array(){
		pSeverInfo, bOk := v.(*common.ServerInfo)
		if bOk && pSeverInfo!= nil{
			SERVER.GetServer().SendMsgByID(pSeverInfo.SocketId, "G_ClientLost", accountId)
		}
	}
}

func (this *ServerSocketManager) Init(num int){
	this.Actor.Init(num)
	this.m_GateMap 		= make(HashSocketMap)
	this.m_SocketMap 	= make(HashSocketMap)
	this.m_WorldMap		= make(HashSocketMap)
	this.m_SocketLocker	= &sync.RWMutex{}
	this.m_GateLocker	= &sync.RWMutex{}
	this.m_WorldLocker	= &sync.RWMutex{}
	this.RegisterCall("DISCONNECT", func(socketid int) {
		this.ReleaseServerMap(socketid, false)
	})

	this.RegisterCall("CONNECT", func(pServerInfo *common.ServerInfo) {
		this.AddServerMap(pServerInfo)
	})

	this.Actor.Start()
}