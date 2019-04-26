package cluster

import (
	"fmt"
	"gonet/actor"
	"gonet/network"
	"reflect"
	"sync"
)

type(
	//集群包管理
	IClusterPacket interface {
		actor.IActor

		RegisterServer(int, string, int)
		SetSocketId(int)
	}

	//集群客户端
	Cluster struct{
		actor.Actor
		m_pClients map[int] *network.ClientSocket
		m_pLocker  *sync.RWMutex
		m_Packet IClusterPacket
		m_PacketFunc network.HandleFunc
	}

	ICluster interface{
		//actor.IActor
		Init(num int, nType int)
		InitSocket([]string, []int, []int)
		AddSocket(string, int, int)
		DelSocket(string, int, int)
		GetSocket(int) *network.ClientSocket

		BindPacket(IClusterPacket)
		BindPacketFunc(network.HandleFunc)
		SendMsg(int, string, ...interface{})
		Send(int, []byte)
	}
)

func GetClusterMsg(msg string, nType int) string{
	return fmt.Sprintf("%s_%d", msg, nType)
}

func (this *Cluster) Init(num int, nType int) {
	this.Actor.Init(num)
	this.m_pLocker = &sync.RWMutex{}
	this.m_pClients = make(map[int] *network.ClientSocket)

	//集群初始化
	this.RegisterCall(GetClusterMsg("Cluster_Socket_Init", nType), func(aIp []string, aPort []int, aSocketId []int){
		this.InitSocket(aIp, aPort, aSocketId)
	})

	//集群新加member
	this.RegisterCall(GetClusterMsg("Cluster_Socket_Add", nType), func(Ip string, Port int, Socket int){
		this.AddSocket(Ip, Port, Socket)
	})

	//集群删除member
	this.RegisterCall(GetClusterMsg("Cluster_Socket_Del", nType), func(Ip string, Port int, Socket int){
		this.DelSocket(Ip, Port, Socket)
	})

	this.Actor.Start()
}

func (this *Cluster) InitSocket(aIp []string, aPort []int, aSocketId []int){
	this.m_pLocker.Lock()
	for _, v := range this.m_pClients{
		v.CallMsg("STOP_ACTOR")
		v.Stop()
	}
	this.m_pClients = make(map[int] *network.ClientSocket)
	this.m_pLocker.Unlock()
	for i, v := range aIp{
		this.AddSocket(v, aPort[i], aSocketId[i])
	}
}

func (this *Cluster) AddSocket(Ip string, Port, SocketId int){
	pClient := new(network.ClientSocket)
	pClient.Init(Ip, Port)
	packet :=  reflect.New(reflect.ValueOf(this.m_Packet).Elem().Type()).Interface().(IClusterPacket)
	packet.Init(1000)
	packet.SetSocketId(SocketId)
	pClient.BindPacketFunc(packet.PacketFunc)
	if this.m_PacketFunc != nil{
		pClient.BindPacketFunc(this.m_PacketFunc)
	}
	this.m_pLocker.Lock()
	this.m_pClients[SocketId] = pClient
	this.m_pLocker.Unlock()
	pClient.Start()
}

func (this *Cluster) DelSocket(Ip string, Port, SocketId int){
	this.m_pLocker.RLock()
	pClient, exist := this.m_pClients[SocketId]
	this.m_pLocker.RUnlock()
	if exist{
		pClient.CallMsg("STOP_ACTOR")
		pClient.Stop()
	}

	this.m_pLocker.Lock()
	delete(this.m_pClients, SocketId)
	this.m_pLocker.Unlock()
}

func (this *Cluster) GetSocket(socketId int) *network.ClientSocket{
	this.m_pLocker.RLock()
	pClient, exist := this.m_pClients[socketId]
	this.m_pLocker.RUnlock()
	if exist{
		return pClient
	}
	return nil
}

func (this *Cluster) BindPacket(packet IClusterPacket){
	this.m_Packet = packet
}

func (this *Cluster) BindPacketFunc(packetFunc network.HandleFunc){
	this.m_PacketFunc = packetFunc
}

func (this *Cluster) SendMsg(socketId int, funcName string, params  ...interface{}){
	pClient := this.GetSocket(socketId)
	if pClient != nil{
		pClient.SendMsg(funcName, params...)
	}
}

func (this *Cluster) Send(socketId int, buff []byte){
	pClient := this.GetSocket(socketId)
	if pClient != nil{
		pClient.Send(buff)
	}
}


