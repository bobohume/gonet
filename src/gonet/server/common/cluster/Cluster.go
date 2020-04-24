package cluster

import (
	"gonet/actor"
	"gonet/base"
	"gonet/message"
	"gonet/network"
	"gonet/rpc"
	"gonet/server/common"
	"gonet/server/common/cluster/etv3"
	"math"
	"reflect"
	"sync"
)

type(
	Service etv3.Service
	Master etv3.Master
	Snowflake etv3.Snowflake
	//集群包管理
	IClusterPacket interface {
		actor.IActor
		SetClusterId(uint32)
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
		Init(num int, MasterType message.SERVICE, IP string, Port int, Endpoints []string)
		AddCluster(info *common.ClusterInfo)
		DelCluster(info *common.ClusterInfo)
		GetCluster(rpc.RpcHead) *network.ClientSocket

		BindPacket(IClusterPacket)
		BindPacketFunc(network.HandleFunc)
		SendMsg(rpc.RpcHead, string, ...interface{})//发送给集群特定服务器
		Send(rpc.RpcHead, []byte)//发送给集群特定服务器

		BalanceCluster()rpc.RpcHead//负载均衡
		RandomCluster()rpc.RpcHead///随机分配

		sendPoint(rpc.RpcHead, []byte)//发送给集群特定服务器
		balanceSend(rpc.RpcHead, []byte)//负载给集群特定服务器
		boardCastSend(rpc.RpcHead, []byte)//给集群广播
	}
)

//注册服务器
func NewService(Type message.SERVICE, IP string, Port int, Endpoints []string) *Service{
	service := &etv3.Service{}
	service.Init(Type, IP, Port, Endpoints)
	return (*Service)(service)
}

//监控服务器
func NewMaster(Type message.SERVICE, Endpoints []string, pActor actor.IActor) *Master {
	master := &etv3.Master{}
	master.Init(Type, Endpoints, pActor)
	return (*Master)(master)
}

//uuid生成器
func NewSnowflake(Endpoints []string) *Snowflake{
	uuid := &etv3.Snowflake{}
	uuid.Init(Endpoints)
	return (*Snowflake)(uuid)
}

func (this *Cluster) Init(num int, MasterType message.SERVICE, IP string, Port int, Endpoints []string) {
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
	pClient.Init(info.Ip, int(info.Port))
	packet :=  reflect.New(reflect.ValueOf(this.m_Packet).Elem().Type()).Interface().(IClusterPacket)
	packet.Init(1000)
	packet.SetClusterId(info.Id())
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
	pClient, bEx := this.m_ClusterMap[info.Id()]
	this.m_ClusterLocker.RUnlock()
	if bEx{
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

func (this *Cluster) GetCluster(head rpc.RpcHead) *network.ClientSocket{
	this.m_ClusterLocker.RLock()
	pClient, bEx := this.m_ClusterMap[head.ClusterId]
	this.m_ClusterLocker.RUnlock()
	if bEx{
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

func (this *Cluster) sendPoint(head rpc.RpcHead, buff []byte){
	pClient := this.GetCluster(head)
	if pClient != nil{
		pClient.Send(head, buff)
	}
}

func (this *Cluster) balanceSend(head rpc.RpcHead, buff []byte){
	head = this.RandomCluster()
	pClient := this.GetCluster(head)
	if pClient != nil{
		pClient.Send(head, buff)
	}
}

func (this *Cluster) boardCastSend(head rpc.RpcHead, buff []byte){
	clusterList := []*network.ClientSocket{}
	this.m_ClusterLocker.RLock()
	for _ ,v := range this.m_ClusterMap{
		clusterList = append(clusterList, v)
	}
	this.m_ClusterLocker.RUnlock()
	for _, v := range clusterList{
		v.Send(head, buff)
	}
}

func (this *Cluster) SendMsg(head rpc.RpcHead, funcName string, params  ...interface{}){
	buff := base.SetTcpEnd(rpc.Marshal(head, funcName, params...))
	this.Send(head, buff)
}

func (this *Cluster) Send(head rpc.RpcHead, buff []byte){
	switch head.SendType{
	case message.SEND_BALANCE:
		this.balanceSend(head, buff)
	case message.SEND_POINT:
		this.sendPoint(head, buff)
	default:
		this.boardCastSend(head, buff)
	}
}

func (this *Cluster) BalanceCluster() rpc.RpcHead{
	head := rpc.RpcHead{}
	this.m_ClusterLocker.RLock()
	for _, v := range this.m_ClusterList.Array(){
		pClusterInfo := v.(*common.ClusterInfo)
		if pClusterInfo.Weight <= 10000{
			head.ClusterId = pClusterInfo.Id()
			head.SocketId = pClusterInfo.SocketId
			break
		}
	}
	this.m_ClusterLocker.RUnlock()
	if head.ClusterId == 0{
		return this.RandomCluster()
	}
	return head
}

func (this *Cluster) RandomCluster() rpc.RpcHead{
	head := rpc.RpcHead{}
	this.m_ClusterLocker.RLock()
	if this.m_ClusterList.Len() > 0{
		nLen := int(math.Max(float64(this.m_ClusterList.Len()-1), 0))
		nRand := base.RAND.RandI(0, nLen)
		info := this.m_ClusterList.Get(nRand).(*common.ClusterInfo)
		head.ClusterId = info.Id()
		head.SocketId = info.SocketId
	}
	this.m_ClusterLocker.RUnlock()
	return head
}