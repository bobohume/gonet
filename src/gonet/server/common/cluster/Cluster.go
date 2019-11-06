package cluster

import (
	"gonet/actor"
	"gonet/base"
	"gonet/network"
	"gonet/server/common"
	"math"
	"reflect"
	"sync"
)

type(
	//集群包管理
	IClusterPacket interface {
		actor.IActor

		RegisterServer(int, string, int)
		SetSocketId(uint32)
	}

	//集群客户端
	Cluster struct{
		actor.Actor
		m_ClusterMap map[uint32] *network.ClientSocket
		m_ClusterList base.IVector
		m_ClusterLocker *sync.RWMutex
		m_Packet IClusterPacket
		m_PacketFunc network.HandleFunc
		m_Master  *Master
	}

	ICluster interface{
		//actor.IActor
		Init(num int, MasterType int, IP string, Port int, Endpoints []string)
		AddCluster(info *common.ClusterInfo)
		DelCluster(info *common.ClusterInfo)
		GetCluster(uint32) *network.ClientSocket

		BindPacket(IClusterPacket)
		BindPacketFunc(network.HandleFunc)
		SendMsg(uint32, string, ...interface{})//发送给集群特定服务器
		BalanceMsg(string, ...interface{})//负载给集群特定服务器
		BoardCastMsg(string, ...interface{})//给集群广播
		Send(uint32, []byte)
		BalanceSend([]byte)//负载给集群特定服务器

		BalanceCluster()uint32//负载均衡
		RandomCluster()uint32//随机分配
	}
)

func (this *Cluster) Init(num int, MasterType int, IP string, Port int, Endpoints []string) {
	this.Actor.Init(num)
	this.m_ClusterLocker = &sync.RWMutex{}
	this.m_ClusterMap = make(map[uint32] *network.ClientSocket)
	this.m_ClusterList = &base.Vector{}
	this.m_Master = NewMaster(MasterType, Endpoints, &this.Actor)

	//集群新加member
	this.RegisterCall("Cluster_Add", func(info *common.ClusterInfo){
		this.AddCluster(info)
	})

	//集群删除member
	this.RegisterCall("Cluster_Del", func(info *common.ClusterInfo){
		this.DelCluster(info)
	})

	this.Actor.Start()
}

func (this *Cluster) AddCluster(info *common.ClusterInfo){
	pClient := new(network.ClientSocket)
	pClient.Init(info.Ip, info.Port)
	packet :=  reflect.New(reflect.ValueOf(this.m_Packet).Elem().Type()).Interface().(IClusterPacket)
	packet.Init(1000)
	packet.SetSocketId(info.Id())
	pClient.BindPacketFunc(packet.PacketFunc)
	if this.m_PacketFunc != nil{
		pClient.BindPacketFunc(this.m_PacketFunc)
	}
	this.m_ClusterLocker.Lock()
	this.m_ClusterMap[info.Id()] = pClient
	this.m_ClusterList.Push_back(info)
	this.m_ClusterLocker.Unlock()
	pClient.Start()
}

func (this *Cluster) DelCluster(info *common.ClusterInfo){
	this.m_ClusterLocker.RLock()
	pClient, exist := this.m_ClusterMap[info.Id()]
	this.m_ClusterLocker.RUnlock()
	if exist{
		pClient.CallMsg("STOP_ACTOR")
		pClient.Stop()
	}

	this.m_ClusterLocker.Lock()
	delete(this.m_ClusterMap, info.Id())
	for i, v := range this.m_ClusterList.Array(){
		if v.(*common.ClusterInfo).Id() == info.Id(){
			this.m_ClusterList.Erase(i)
			break
		}
	}
	this.m_ClusterLocker.Unlock()
}

func (this *Cluster) GetCluster(socketId uint32) *network.ClientSocket{
	this.m_ClusterLocker.RLock()
	pClient, exist := this.m_ClusterMap[socketId]
	this.m_ClusterLocker.RUnlock()
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

func (this *Cluster) SendMsg(socketId uint32, funcName string, params  ...interface{}){
	pClient := this.GetCluster(socketId)
	if pClient != nil{
		pClient.SendMsg(funcName, params...)
	}
}

func (this *Cluster) BalanceMsg(funcName string, params  ...interface{}){
	pClient := this.GetCluster(this.RandomCluster())
	if pClient != nil{
		pClient.SendMsg(funcName, params...)
	}
}

func (this *Cluster) BoardCastMsg(funcName string, params  ...interface{}){
	clusterList := []*network.ClientSocket{}
	this.m_ClusterLocker.RLock()
	for _ ,v := range this.m_ClusterMap{
		clusterList = append(clusterList, v)
	}
	this.m_ClusterLocker.RUnlock()
	for _, v := range clusterList{
		v.SendMsg(funcName, params...)
	}
}

func (this *Cluster) Send(socketId uint32, buff []byte){
	pClient := this.GetCluster(socketId)
	if pClient != nil{
		pClient.Send(buff)
	}
}

func (this *Cluster) BalanceSend(buff []byte){
	pClient := this.GetCluster(this.RandomCluster())
	if pClient != nil{
		pClient.Send(buff)
	}
}

func (this *Cluster) BalanceCluster() uint32{
	nIndex := uint32(0)
	this.m_ClusterLocker.RLock()
	for _, v := range this.m_ClusterList.Array(){
		pClusterInfo := v.(*common.ClusterInfo)
		if pClusterInfo.Weight <= 10000{
			nIndex = pClusterInfo.Id()
			break
		}
	}
	this.m_ClusterLocker.RUnlock()
	if nIndex == 0{
		nIndex = this.RandomCluster()
	}
	return nIndex
}

func (this *Cluster) RandomCluster() uint32{
	nIndex := uint32(0)
	this.m_ClusterLocker.RLock()
	if this.m_ClusterList.Len() > 0{
		nLen := int(math.Max(float64(this.m_ClusterList.Len()-1), 0))
		nRand := base.RAND.RandI(0, nLen)
		nIndex = this.m_ClusterList.Get(nRand).(*common.ClusterInfo).Id()
	}
	this.m_ClusterLocker.RUnlock()
	return nIndex
}