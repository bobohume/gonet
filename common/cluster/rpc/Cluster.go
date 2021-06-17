package rpc

import (
	"context"
	"gonet/actor"
	"gonet/base"
	"gonet/base/vector"
	"gonet/common"
	"gonet/common/cluster/etv3"
	"gonet/network"
	"gonet/rpc"
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

	ClusterNode struct {
		*network.ClientSocket
		*common.ClusterInfo
	}

	//集群客户端
	Cluster struct{
		actor.Actor
		m_ClusterMap map[uint32] *ClusterNode
		m_ClusterLocker *sync.RWMutex
		m_Packet IClusterPacket
		m_Master  *Master
		m_HashRing	*base.HashRing//hash一致性
		m_ClusterInfoMap map[uint32] *common.ClusterInfo
		m_PacketFuncList	*vector.Vector//call back
	}

	ICluster interface{
		//actor.IActor
		Init(info *common.ClusterInfo, Endpoints []string)
		AddCluster(info *common.ClusterInfo)
		DelCluster(info *common.ClusterInfo)
		GetCluster(rpc.RpcHead) *ClusterNode

		BindPacket(IClusterPacket)
		BindPacketFunc(network.PacketFunc)
		SendMsg(rpc.RpcHead, string, ...interface{})//发送给集群特定服务器
		Send(rpc.RpcHead, []byte)//发送给集群特定服务器

		RandomCluster(head rpc.RpcHead) rpc.RpcHead///随机分配

		sendPoint(rpc.RpcHead, []byte)//发送给集群特定服务器
		balanceSend(rpc.RpcHead, []byte)//负载给集群特定服务器
		boardCastSend(rpc.RpcHead, []byte)//给集群广播
	}
)

//注册服务器
func NewService(info *common.ClusterInfo, Endpoints []string) *Service{
	service := &etv3.Service{}
	service.Init(info, Endpoints)
	return (*Service)(service)
}

//监控服务器
func NewMaster(info *common.ClusterInfo, Endpoints []string, pActor actor.IActor) *Master {
	master := &etv3.Master{}
	master.Init(info, Endpoints, pActor)
	return (*Master)(master)
}

//uuid生成器
func NewSnowflake(Endpoints []string) *Snowflake{
	uuid := &etv3.Snowflake{}
	uuid.Init(Endpoints)
	return (*Snowflake)(uuid)
}

func (this *Cluster) Init(info *common.ClusterInfo, Endpoints []string) {
	this.Actor.Init()
	this.m_ClusterLocker = &sync.RWMutex{}
	this.m_ClusterMap = make(map[uint32] *ClusterNode)
	this.m_Master = NewMaster(info, Endpoints, &this.Actor)
	this.m_HashRing = base.NewHashRing()
	this.m_ClusterInfoMap = make(map[uint32]*common.ClusterInfo)
	this.m_PacketFuncList = vector.NewVector()

	//集群新加member
	this.RegisterCall("Cluster_Add", func(ctx context.Context, info *common.ClusterInfo){
		_, bEx := this.m_ClusterInfoMap[info.Id()]
		if !bEx {
			this.AddCluster(info)
			this.m_ClusterInfoMap[info.Id()] = info
		}
	})

	//集群删除member
	this.RegisterCall("Cluster_Del", func(ctx context.Context, info *common.ClusterInfo){
		delete(this.m_ClusterInfoMap, info.Id())
		this.DelCluster(info)
	})

	//链接断开
	this.RegisterCall("DISCONNECT", func(ctx context.Context, ClusterId uint32) {
		pInfo, bEx := this.m_ClusterInfoMap[ClusterId]
		if bEx {
			this.DelCluster(pInfo)
		}
		delete(this.m_ClusterInfoMap, ClusterId)
	})

	this.Actor.Start()
}

func (this *Cluster) AddCluster(info *common.ClusterInfo){
	pClient := new(network.ClientSocket)
	pClient.Init(info.Ip, int(info.Port))
	packet :=  reflect.New(reflect.ValueOf(this.m_Packet).Elem().Type()).Interface().(IClusterPacket)
	packet.Init()
	packet.SetClusterId(info.Id())
	pClient.BindPacketFunc(packet.PacketFunc)
	for _, v := range this.m_PacketFuncList.Values(){
		pClient.BindPacketFunc(v.(network.PacketFunc))
	}
	this.m_ClusterLocker.Lock()
	this.m_ClusterMap[info.Id()] = &ClusterNode{ClientSocket:pClient, ClusterInfo:info}
	this.m_ClusterLocker.Unlock()
	this.m_HashRing.Add(info.IpString())
	pClient.Start()
}

func (this *Cluster) DelCluster(info *common.ClusterInfo){
	this.m_ClusterLocker.RLock()
	pCluster, bEx := this.m_ClusterMap[info.Id()]
	this.m_ClusterLocker.RUnlock()
	if bEx{
		pCluster.CallMsg("STOP_ACTOR")
		pCluster.Stop()
	}

	this.m_ClusterLocker.Lock()
	delete(this.m_ClusterMap, info.Id())
	this.m_ClusterLocker.Unlock()
	this.m_HashRing.Remove(info.IpString())
}

func (this *Cluster) GetCluster(head rpc.RpcHead) *ClusterNode{
	this.m_ClusterLocker.RLock()
	pCluster, bEx := this.m_ClusterMap[head.ClusterId]
	this.m_ClusterLocker.RUnlock()
	if bEx{
		return pCluster
	}
	return nil
}

func (this *Cluster) BindPacketFunc(callfunc network.PacketFunc){
	this.m_PacketFuncList.PushBack(callfunc)
}

func (this *Cluster) BindPacket(packet IClusterPacket){
	this.m_Packet = packet
}


func (this *Cluster) sendPoint(head rpc.RpcHead, buff []byte){
	pCluster := this.GetCluster(head)
	if pCluster != nil{
		pCluster.Send(head, buff)
	}
}

func (this *Cluster) balanceSend(head rpc.RpcHead, buff []byte){
	_, head.ClusterId = this.m_HashRing.Get64(head.Id)
	pClient := this.GetCluster(head)
	if pClient != nil{
		pClient.Send(head, buff)
	}
}

func (this *Cluster) boardCastSend(head rpc.RpcHead, buff []byte){
	clusterList := []*ClusterNode{}
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
	buff := rpc.Marshal(head, funcName, params...)
	this.Send(head, buff)
}

func (this *Cluster) Send(head rpc.RpcHead, buff []byte){
	switch head.SendType{
	case rpc.SEND_BALANCE:
		this.balanceSend(head, buff)
	case rpc.SEND_POINT:
		this.sendPoint(head, buff)
	default:
		this.boardCastSend(head, buff)
	}
}

func (this *Cluster) RandomCluster(head rpc.RpcHead) rpc.RpcHead{
	if head.Id == 0{
		head.Id = int64(uint32(base.RAND.RandI(1, 0xFFFFFFFF)))
	}
	_, head.ClusterId = this.m_HashRing.Get64(head.Id)
	pCluster := this.GetCluster(head)
	if pCluster != nil{
		head.SocketId = pCluster.SocketId
	}
	return head
}