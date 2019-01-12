package common

import (
	"actor"
	"network"
	"reflect"
	"sync"
)

type(
	IDispatchPacket interface {
		actor.IActor

		RegisterServer(int, string, int)
		SetSocketId(int)
	}

	DispatchMgr struct{
		actor.Actor
		m_pClients map[int] *network.ClientSocket
		m_pLocker  *sync.RWMutex
		m_Packet IDispatchPacket
		m_PacketFunc network.HandleFunc
	}

	IDispatchMgr interface{
		//actor.IActor
		InitSocket([]string, []int, []int)
		AddSocket(string, int, int)
		DelSocket(int)
		GetSocket(int) *network.ClientSocket
		BalanceSocket()int

		BindPacket(IDispatchPacket)
		BindPacketFunc(network.HandleFunc)
		SendMsg(int, string, ...interface{})
		Send(int, []byte)
	}
)

func (this *DispatchMgr) Init(num int) {
	this.Actor.Init(num)
	this.m_pLocker = &sync.RWMutex{}
	this.m_pClients = make(map[int] *network.ClientSocket)

	this.RegisterCall("Dispatch_Socket_Init", func(aIp []string, aPort []int, aSocket []int){
		this.InitSocket(aIp, aPort, aSocket)
	})

	this.RegisterCall("Dispatch_Socket_Add", func(Ip string, Port int, Socket int){
		this.AddSocket(Ip, Port, Socket)
	})

	this.RegisterCall("Dispatch_Socket_Del", func(Ip string, Port int, Socket int){
		this.DelSocket(Socket)
	})

	this.Actor.Start()
}

func (this *DispatchMgr) InitSocket(aIp []string, aPort []int, aSocket []int){
	this.m_pLocker.Lock()
	for _, v := range this.m_pClients{
		v.CallMsg("STOP_ACTOR")
		v.Stop()
	}
	this.m_pClients = make(map[int] *network.ClientSocket)
	this.m_pLocker.Unlock()
	for i, v := range aIp{
		this.AddSocket(v, aPort[i], aSocket[i])
	}
}

func (this *DispatchMgr) AddSocket(ip string, port, socketId int){
	pClient := new(network.ClientSocket)
	pClient.Init(ip, port)
	packet :=  reflect.New(reflect.ValueOf(this.m_Packet).Elem().Type()).Interface().(IDispatchPacket)
	packet.Init(1000)
	packet.SetSocketId(socketId)
	pClient.BindPacketFunc(packet.PacketFunc)
	if this.m_PacketFunc != nil{
		pClient.BindPacketFunc(this.m_PacketFunc)
	}
	this.m_pLocker.Lock()
	this.m_pClients[socketId] = pClient
	this.m_pLocker.Unlock()
	pClient.Start()
}

func (this *DispatchMgr) DelSocket(socketId int){
	this.m_pLocker.RLock()
	pClient, exist := this.m_pClients[socketId]
	this.m_pLocker.RUnlock()
	if exist{
		pClient.CallMsg("STOP_ACTOR")
		pClient.Stop()
	}

	this.m_pLocker.Lock()
	delete(this.m_pClients, socketId)
	this.m_pLocker.Unlock()
}

func (this *DispatchMgr) GetSocket(socketId int) *network.ClientSocket{
	this.m_pLocker.RLock()
	pClient, exist := this.m_pClients[socketId]
	this.m_pLocker.RUnlock()
	if exist{
		return pClient
	}
	return nil
}

func (this *DispatchMgr) BalanceSocket()int{
	nIndex := 0
	this.m_pLocker.RLock()
	for i, _ := range this.m_pClients{
		nIndex =  i
		break
	}
	this.m_pLocker.RUnlock()
	return nIndex
}

func (this *DispatchMgr) BindPacket(packet IDispatchPacket){
	this.m_Packet = packet
}

func (this *DispatchMgr) BindPacketFunc(packetFunc network.HandleFunc){
	this.m_PacketFunc = packetFunc
}

func (this *DispatchMgr) SendMsg(socketId int, funcName string, params  ...interface{}){
	pClient := this.GetSocket(socketId)
	if pClient != nil{
		pClient.SendMsg(funcName, params...)
	}
}

func (this *DispatchMgr) Send(socketId int, buff []byte){
	pClient := this.GetSocket(socketId)
	if pClient != nil{
		pClient.Send(buff)
	}
}


