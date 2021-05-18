package cluster

import (
	"context"
	"github.com/nats-io/nats.go"
	"gonet/actor"
	"gonet/base"
	"gonet/base/vector"
	"gonet/common"
	"gonet/network"
	"gonet/rpc"
	"log"
	"sync"
)

const(
	MAX_CLUSTER_NUM = int(rpc.SERVICE_ZONESERVER) + 1
)

type(
	HashClusterMap map[uint32] *common.ClusterInfo
	HashClusterSocketMap map[uint32] *common.ClusterInfo

	//集群服务器
	Cluster struct{
		actor.Actor
		*Service //集群注册
		m_ClusterMap [MAX_CLUSTER_NUM]HashClusterMap
		m_ClusterLocker [MAX_CLUSTER_NUM]*sync.RWMutex
		m_HashRing	[MAX_CLUSTER_NUM]*base.HashRing//hash一致性
		m_Conn      *nats.Conn
		m_DieChan	chan bool
		m_Master    *Master
		m_ClusterInfoMap map[uint32] *common.ClusterInfo
		m_PacketFuncList	*vector.Vector//call back
	}

	ICluster interface{
		Init(num int, info *common.ClusterInfo, Endpoints []string, natsUrl string)
		RegisterClusterCall()//注册集群通用回调
		AddCluster(info *common.ClusterInfo)
		DelCluster(info *common.ClusterInfo)
		GetCluster(rpc.RpcHead) *common.ClusterInfo

		BindPacketFunc(network.HandleFunc)
		SendMsg(rpc.RpcHead, string, ...interface{})//发送给集群特定服务器
		Send(rpc.RpcHead, []byte)//发送给集群特定服务器

		RandomCluster(head rpc.RpcHead)	rpc.RpcHead//随机分配
	}

	EmptyClusterInfo struct {
		common.ClusterInfo
	}
)

func (this *EmptyClusterInfo) String() string{
	return ""
}

func (this *Cluster) Init(num int, info *common.ClusterInfo, Endpoints []string, natsUrl string) {
	this.Actor.Init(num)
	this.RegisterClusterCall()
	for  i := 0; i < MAX_CLUSTER_NUM; i++{
		this.m_ClusterLocker[i] = &sync.RWMutex{}
		this.m_ClusterMap[i] = make(HashClusterMap)
		this.m_HashRing[i] = base.NewHashRing()
	}
	//注册服务器
	this.Service = NewService(info, Endpoints)
	this.m_Master = NewMaster(&EmptyClusterInfo{}, Endpoints, &this.Actor)
	this.m_ClusterInfoMap = make(map[uint32]*common.ClusterInfo)
	this.m_PacketFuncList = vector.NewVector()

	conn, err := setupNatsConn(
		natsUrl,
		this.m_DieChan,
	)
	if err != nil {
		log.Fatal("nats connect error!!!!")
	}
	this.m_Conn = conn

	this.m_Conn.Subscribe(getChannel(*info), func(msg *nats.Msg) {
		this.HandlePacket(0, msg.Data)
	})

	this.m_Conn.Subscribe(getTopicChannel(*info), func(msg *nats.Msg) {
		this.HandlePacket(0, msg.Data)
	})

	this.Actor.Start()
}

func (this *Cluster) RegisterClusterCall(){
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
}

func (this *Cluster) AddCluster(info *common.ClusterInfo){
	this.m_ClusterLocker[info.Type].Lock()
	this.m_ClusterMap[info.Type][info.Id()] = info
	this.m_ClusterLocker[info.Type].Unlock()
	this.m_HashRing[info.Type].Add(info.IpString())
	base.GLOG.Printf("服务器[%s:%s:%d]建立连接", info.String(), info.Ip, info.Port)
}

func (this *Cluster) DelCluster(info *common.ClusterInfo){
	this.m_ClusterLocker[info.Type].RLock()
	_, bEx := this.m_ClusterMap[info.Type][info.Id()]
	this.m_ClusterLocker[info.Type].RUnlock()
	if bEx{
		this.m_ClusterLocker[info.Type].Lock()
		delete(this.m_ClusterMap[info.Type], info.Id())
		this.m_ClusterLocker[info.Type].Unlock()
	}

	this.m_HashRing[info.Type].Remove(info.IpString())
	base.GLOG.Printf("服务器[%s:%s:%d]断开连接", info.String(), info.Ip, info.Port)
}

func (this *Cluster) GetCluster(head rpc.RpcHead) *common.ClusterInfo {
	this.m_ClusterLocker[head.DestServerType].RLock()
	defer this.m_ClusterLocker[head.DestServerType].RUnlock()
	pClient, bEx := this.m_ClusterMap[head.DestServerType][head.ClusterId]
	if bEx{
		return pClient
	}
	return nil
}


func (this *Cluster) BindPacketFunc(callfunc network.HandleFunc){
	this.m_PacketFuncList.PushBack(callfunc)
}

func (this *Cluster) HandlePacket(Id uint32, buff []byte){
	for _,v := range this.m_PacketFuncList.Values() {
		if (v.(network.HandleFunc)(Id, buff)){
			break
		}
	}
}

func (this *Cluster) SendMsg(head rpc.RpcHead, funcName string, params  ...interface{}){
	head.SrcClusterId = this.Id()
	buff := rpc.Marshal(head, funcName, params...)
	this.Send(head, buff)
}

func (this *Cluster) Send(head rpc.RpcHead, buff []byte){
	switch head.SendType{
	case rpc.SEND_BALANCE:
		_, head.ClusterId = this.m_HashRing[head.DestServerType].Get64(head.Id)
		this.m_Conn.Publish(getRpcChannel(head) ,buff)
	case rpc.SEND_POINT:
		this.m_Conn.Publish(getRpcChannel(head) ,buff)
	default:
		this.m_Conn.Publish(getRpcTopicChannel(head), buff)
	}
}

func (this *Cluster) RandomCluster(head rpc.RpcHead) rpc.RpcHead{
	if head.Id == 0{
		head.Id = int64(uint32(base.RAND.RandI(1, 0xFFFFFFFF)))
	}
	_, head.ClusterId = this.m_HashRing[head.DestServerType].Get64(head.Id)
	pCluster := this.GetCluster(head)
	if pCluster != nil{
		head.SocketId = pCluster.SocketId
	}
	return head
}
