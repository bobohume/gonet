package netgate

import (
	"actor"
	"network"
	"sync"
)

type(
	DispatchMgr struct{
		actor.Actor
		m_pWorldClients map[int] *network.ClientSocket
		m_pWorldLocker *sync.RWMutex
	}

	IDispatchMgr interface{
		//actor.IActor

		InitWorldSocket(aIp []string, aPort []int, aSocket []int)
		AddWorldSocket(string, int, int)
		DelWorldSocket(int)
		GetWorldSocket(int) *network.ClientSocket
		BalanceWorldSocket()int
		SendMsg(int, string, ...interface{})
		Send(int, []byte)
	}
)

func (this *DispatchMgr) Init(num int) {
	this.Actor.Init(num)
	this.m_pWorldLocker = &sync.RWMutex{}
	this.m_pWorldClients = make(map[int] *network.ClientSocket)
	this.RegisterCall("Dispatch_Socket_Init", func(aIp []string, aPort []int, aSocket []int){
		this.InitWorldSocket(aIp, aPort, aSocket)
	})

	this.RegisterCall("Dispatch_Socket_Add", func(Ip string, Port int, Socket int){
		this.AddWorldSocket(Ip, Port, Socket)
	})

	this.RegisterCall("Dispatch_Socket_Del", func(Ip string, Port int, Socket int){
		this.DelWorldSocket(Socket)
	})

	this.Actor.Start()
}

func (this *DispatchMgr) InitWorldSocket(aIp []string, aPort []int, aSocket []int){
	this.m_pWorldLocker.Lock()
	for _, v := range this.m_pWorldClients{
		v.CallMsg("STOP_ACTOR")
		v.Stop()
	}
	this.m_pWorldClients = make(map[int] *network.ClientSocket)
	this.m_pWorldLocker.Unlock()
	for i, v := range aIp{
		this.AddWorldSocket(v, aPort[i], aSocket[i])
	}
}

func (this *DispatchMgr) AddWorldSocket(ip string, port, socketId int){
	pClient := new(network.ClientSocket)
	pClient.Init(ip, port)
	packet := new(WorldProcess)
	packet.Init(1000)
	packet.SetSocketId(socketId)
	pClient.BindPacketFunc(packet.PacketFunc)
	pClient.BindPacketFunc(DispatchPacketToClient)
	this.m_pWorldLocker.Lock()
	this.m_pWorldClients[socketId] = pClient
	this.m_pWorldLocker.Unlock()
	pClient.Start()
}

func (this *DispatchMgr) DelWorldSocket(socketId int){
	this.m_pWorldLocker.RLock()
	pClient, exist := this.m_pWorldClients[socketId]
	this.m_pWorldLocker.RUnlock()
	if exist{
		pClient.CallMsg("STOP_ACTOR")
		pClient.Stop()
	}

	this.m_pWorldLocker.Lock()
	delete(this.m_pWorldClients, socketId)
	this.m_pWorldLocker.Unlock()
}

func (this *DispatchMgr) GetWorldSocket(socketId int) *network.ClientSocket{
	this.m_pWorldLocker.RLock()
	pClient, exist := this.m_pWorldClients[socketId]
	this.m_pWorldLocker.RUnlock()
	if exist{
		return pClient
	}
	return nil
}

func (this *DispatchMgr) BalanceWorldSocket()int{
	nIndex := 0
	this.m_pWorldLocker.RLock()
	for i, _ := range this.m_pWorldClients{
		nIndex =  i
		break
	}
	this.m_pWorldLocker.RUnlock()
	return nIndex
}

func (this *DispatchMgr) SendMsg(socketId int, funcName string, params  ...interface{}){
	pClient := this.GetWorldSocket(socketId)
	if pClient != nil{
		pClient.SendMsg(funcName, params...)
	}
}

func (this *DispatchMgr) Send(socketId int, buff []byte){
	pClient := this.GetWorldSocket(socketId)
	if pClient != nil{
		pClient.Send(buff)
	}
}
