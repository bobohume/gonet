package rpc

import (
	"context"
	"gonet/actor"
	"gonet/base"
	"gonet/base/cluster/etv3"
	"gonet/base/vector"
	"gonet/network"
	"gonet/rpc"
	"reflect"
	"sync"
)

type (
	Service   etv3.Service
	Master    etv3.Master
	Snowflake etv3.Snowflake
	MailBox   etv3.MailBox
	//集群包管理
	IClusterPacket interface {
		actor.IActor
		SetClusterId(uint32)
	}

	ClusterNode struct {
		*network.ClientSocket
		*rpc.ClusterInfo
	}

	//集群客户端
	Cluster struct {
		actor.Actor
		clusterMap      map[uint32]*ClusterNode
		clusterLocker   *sync.RWMutex
		packet          IClusterPacket
		master          *Master
		hashRing        *base.HashRing //hash一致性
		clusterInfoMapg map[uint32]*rpc.ClusterInfo
		packetFuncList  *vector.Vector[network.PacketFunc] //call back
	}

	ICluster interface {
		actor.IActor
		InitCluster(info *rpc.ClusterInfo, Endpoints []string)
		AddCluster(info *rpc.ClusterInfo)
		DelCluster(info *rpc.ClusterInfo)
		GetCluster(rpc.RpcHead) *ClusterNode

		BindPacket(IClusterPacket)
		BindPacketFunc(network.PacketFunc)

		RandomCluster(head rpc.RpcHead) rpc.RpcHead ///随机分配

		sendPoint(rpc.RpcHead, rpc.Packet)     //发送给集群特定服务器
		balanceSend(rpc.RpcHead, rpc.Packet)   //负载给集群特定服务器
		boardCastSend(rpc.RpcHead, rpc.Packet) //给集群广播
	}
)

// 注册服务器
func NewService(info *rpc.ClusterInfo, Endpoints []string) *Service {
	service := &etv3.Service{}
	service.Init(info, Endpoints)
	return (*Service)(service)
}

// 监控服务器
func NewMaster(info *rpc.ClusterInfo, Endpoints []string) *Master {
	master := &etv3.Master{}
	master.Init(info, Endpoints)
	return (*Master)(master)
}

// uuid生成器
func NewSnowflake(Endpoints []string) *Snowflake {
	uuid := &etv3.Snowflake{}
	uuid.Init(Endpoints)
	return (*Snowflake)(uuid)
}

func (c *Cluster) InitCluster(info *rpc.ClusterInfo, Endpoints []string) {
	c.Actor.Init()
	c.clusterLocker = &sync.RWMutex{}
	c.clusterMap = make(map[uint32]*ClusterNode)
	c.master = NewMaster(info, Endpoints)
	c.hashRing = base.NewHashRing()
	c.clusterInfoMapg = make(map[uint32]*rpc.ClusterInfo)
	c.packetFuncList = &vector.Vector[network.PacketFunc]{}
	actor.MGR.RegisterActor(c)
	c.Actor.Start()
}

// 集群新加member
func (c *Cluster) Cluster_Add(ctx context.Context, info *rpc.ClusterInfo) {
	_, bEx := c.clusterInfoMapg[info.Id()]
	if !bEx {
		c.AddCluster(info)
		c.clusterInfoMapg[info.Id()] = info
	}
}

// 集群删除member
func (c *Cluster) Cluster_Del(ctx context.Context, info *rpc.ClusterInfo) {
	delete(c.clusterInfoMapg, info.Id())
	c.DelCluster(info)
}

// 链接断开
func (c *Cluster) DISCONNECT(ctx context.Context, ClusterId uint32) {
	info, bEx := c.clusterInfoMapg[ClusterId]
	if bEx {
		c.DelCluster(info)
	}
	delete(c.clusterInfoMapg, ClusterId)
}

func (c *Cluster) AddCluster(info *rpc.ClusterInfo) {
	client := new(network.ClientSocket)
	client.Init(info.Ip, int(info.Port))
	packet := reflect.New(reflect.ValueOf(c.packet).Elem().Type()).Interface().(IClusterPacket)
	packet.Init()
	packet.SetClusterId(info.Id())
	client.BindPacketFunc(actor.MGR.PacketFunc)
	for _, v := range c.packetFuncList.Values() {
		client.BindPacketFunc(v)
	}
	c.clusterLocker.Lock()
	c.clusterMap[info.Id()] = &ClusterNode{ClientSocket: client, ClusterInfo: info}
	c.clusterLocker.Unlock()
	c.hashRing.Add(info.IpString())
	client.Start()
}

func (c *Cluster) DelCluster(info *rpc.ClusterInfo) {
	c.clusterLocker.RLock()
	cluster, bEx := c.clusterMap[info.Id()]
	c.clusterLocker.RUnlock()
	if bEx {
		cluster.CallMsg(rpc.RpcHead{}, "STOP_ACTOR")
		cluster.Stop()
	}

	c.clusterLocker.Lock()
	delete(c.clusterMap, info.Id())
	c.clusterLocker.Unlock()
	c.hashRing.Remove(info.IpString())
}

func (c *Cluster) GetCluster(head rpc.RpcHead) *ClusterNode {
	c.clusterLocker.RLock()
	cluster, bEx := c.clusterMap[head.ClusterId]
	c.clusterLocker.RUnlock()
	if bEx {
		return cluster
	}
	return nil
}

func (c *Cluster) BindPacketFunc(callfunc network.PacketFunc) {
	c.packetFuncList.PushBack(callfunc)
}

func (c *Cluster) BindPacket(packet IClusterPacket) {
	c.packet = packet
}

func (c *Cluster) sendPoint(head rpc.RpcHead, packet rpc.Packet) {
	pCluster := c.GetCluster(head)
	if pCluster != nil {
		pCluster.Send(head, packet)
	}
}

func (c *Cluster) balanceSend(head rpc.RpcHead, packet rpc.Packet) {
	_, head.ClusterId = c.hashRing.Get64(head.Id)
	pClient := c.GetCluster(head)
	if pClient != nil {
		pClient.Send(head, packet)
	}
}

func (c *Cluster) boardCastSend(head rpc.RpcHead, packet rpc.Packet) {
	clusterList := []*ClusterNode{}
	c.clusterLocker.RLock()
	for _, v := range c.clusterMap {
		clusterList = append(clusterList, v)
	}
	c.clusterLocker.RUnlock()
	for _, v := range clusterList {
		v.Send(head, packet)
	}
}

func (c *Cluster) SendMsg(head rpc.RpcHead, funcName string, params ...interface{}) {
	c.Send(head, rpc.Marshal(&head, &funcName, params...))
}

func (c *Cluster) Send(head rpc.RpcHead, packet rpc.Packet) {
	switch head.SendType {
	//case rpc.SEND_BALANCE:
	//	c.balanceSend(head, packet)
	case rpc.SEND_POINT:
		c.sendPoint(head, packet)
	default:
		c.boardCastSend(head, packet)
	}
}

func (c *Cluster) RandomCluster(head rpc.RpcHead) rpc.RpcHead {
	if head.Id == 0 {
		head.Id = int64(uint32(base.RandI(1, 0xFFFFFFFF)))
	}
	_, head.ClusterId = c.hashRing.Get64(head.Id)
	pCluster := c.GetCluster(head)
	if pCluster != nil {
		head.SocketId = pCluster.SocketId
	}
	return head
}
