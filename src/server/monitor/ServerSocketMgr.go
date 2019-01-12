package monitor

import (
	"actor"
	"base"
	"fmt"
	"message"
	"server/common"
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
		m_AccountMap	HashSocketMap
		m_AccountLocker	*sync.RWMutex
		m_WorkIdQue		base.IWorkIdQue
	}

	IServerSocketManager interface {
		actor.IActor
		AddServerMap(*common.ServerInfo)
		ReleaseServerMap(int, bool)
		GetAllGate(base.IVector)
		GetAllWorld(base.IVector)
		GetAllAccount(base.IVector)
		GetSeverInfo(int, int) *common.ServerInfo
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
	case int(message.SERVICE_WORLDSERVER):
		workId := this.m_WorkIdQue.Add(fmt.Sprintf("%s:%d", pServerInfo.Ip, pServerInfo.Port))
		SERVER.GetServer().SendMsgByID(pServerInfo.SocketId, "SnowFlake_WorkId", workId)
		this.m_WorldLocker.Lock()
		this.m_WorldMap[pServerInfo.SocketId] = pServerInfo
		this.m_WorldLocker.Unlock()
		SERVER.GetLog().Printf("ADD WORLD SERVER: [%d]-[%s:%d]", pServerInfo.SocketId, pServerInfo.Ip, pServerInfo.Port)
	case int(message.SERVICE_ACCOUNTSERVER):
		workId := this.m_WorkIdQue.Add(fmt.Sprintf("%s:%d", pServerInfo.Ip, pServerInfo.Port))
		SERVER.GetServer().SendMsgByID(pServerInfo.SocketId, "SnowFlake_WorkId", workId)
		this.m_AccountLocker.Lock()
		this.m_AccountMap[pServerInfo.SocketId] = pServerInfo
		this.m_AccountLocker.Unlock()
		SERVER.GetLog().Printf("ADD ACCOUNT SERVER: [%d]-[%s:%d]", pServerInfo.SocketId, pServerInfo.Ip, pServerInfo.Port)
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
		this.m_WorkIdQue.Del(fmt.Sprintf("%s:%d", pServerInfo.Ip, pServerInfo.Port))
		SERVER.GetLog().Printf("与World服务器断开连接,id[%d]",pServerInfo.SocketId)
		this.m_WorldLocker.Lock()
		delete(this.m_WorldMap, pServerInfo.SocketId)
		this.m_WorldLocker.Unlock()
	case int(message.SERVICE_ACCOUNTSERVER):
		this.m_WorkIdQue.Del(fmt.Sprintf("%s:%d", pServerInfo.Ip, pServerInfo.Port))
		SERVER.GetLog().Printf("与Account服务器断开连接,id[%d]",pServerInfo.SocketId)
		this.m_AccountLocker.Lock()
		delete(this.m_AccountMap, pServerInfo.SocketId)
		this.m_AccountLocker.Unlock()
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
	case int(message.SERVICE_ACCOUNTSERVER):
		this.m_AccountLocker.RLock()
		pServerInfo, _ := this.m_AccountMap[socketId]
		this.m_AccountLocker.RUnlock()
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

func (this *ServerSocketManager) GetAllAccount(vec base.IVector){
	this.m_AccountLocker.RLock()
	for _, v := range this.m_AccountMap{
		vec.Push_back(v)
	}
	this.m_AccountLocker.RUnlock()
}

func (this *ServerSocketManager) Init(num int){
	this.Actor.Init(num)
	this.m_GateMap 		= make(HashSocketMap)
	this.m_SocketMap 	= make(HashSocketMap)
	this.m_WorldMap		= make(HashSocketMap)
	this.m_AccountMap   = make(HashSocketMap)
	this.m_SocketLocker	= &sync.RWMutex{}
	this.m_GateLocker	= &sync.RWMutex{}
	this.m_WorldLocker	= &sync.RWMutex{}
	this.m_AccountLocker = &sync.RWMutex{}
	this.m_WorkIdQue = &base.WorkIdQue{}
	this.m_WorkIdQue.Init(0)
	this.m_WorkIdQue.Add(UserNetIP + ":" + UserNetPort)
	this.RegisterCall("DISCONNECT", func(socketid int) {
		this.ReleaseServerMap(socketid, false)
	})

	this.RegisterCall("CONNECT", func(pServerInfo *common.ServerInfo) {
		this.AddServerMap(pServerInfo)
	})

	this.Actor.Start()
}
