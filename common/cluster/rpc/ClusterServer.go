package rpc

import (
	"gonet/actor"
	"gonet/base"
	"gonet/common"
	"gonet/network"
	"gonet/rpc"
	"sync"

	"golang.org/x/net/context"
)

type (
	HashClusterMap       map[uint32]*common.ClusterInfo
	HashClusterSocketMap map[uint32]*common.ClusterInfo

	//集群服务器
	ClusterServer struct {
		actor.Actor
		*Service         //集群注册
		clusterMap       HashClusterMap
		clusterSocketMap HashClusterSocketMap
		clusterLocker    *sync.RWMutex
		service          *network.ServerSocket //socket管理
		hashRing         *base.HashRing        //hash一致性
	}

	IClusterServer interface {
		InitService(info *common.ClusterInfo, Endpoints []string)
		RegisterClusterCall() //注册集群通用回调
		AddCluster(info *common.ClusterInfo)
		DelCluster(info *common.ClusterInfo)
		GetCluster(rpc.RpcHead) *common.ClusterInfo
		GetClusterBySocket(uint32) *common.ClusterInfo

		BindServer(*network.ServerSocket)
		SendMsg(rpc.RpcHead, string, ...interface{}) //发送给集群特定服务器
		Send(rpc.RpcHead, []byte)                    //发送给集群特定服务器

		RandomCluster(head rpc.RpcHead) rpc.RpcHead //随机分配

		sendPoint(rpc.RpcHead, []byte)               //发送给集群特定服务器
		balanceSend(rpc.RpcHead, []byte)             //负载给集群特定服务器
		boardCastSend(head rpc.RpcHead, buff []byte) //给集群广播
	}
)

func (c *ClusterServer) InitService(info *common.ClusterInfo, Endpoints []string) {
	c.Actor.Init()
	c.clusterLocker = &sync.RWMutex{}
	//注册服务器
	c.Service = NewService(info, Endpoints)
	c.clusterMap = make(HashClusterMap)
	c.clusterSocketMap = make(HashClusterSocketMap)
	c.hashRing = base.NewHashRing()
	actor.MGR.RegisterActor(c)
}

func (c *ClusterServer) RegisterClusterCall() {
}

func (c *ClusterServer) AddCluster(info *common.ClusterInfo) {
	c.clusterLocker.Lock()
	c.clusterMap[info.Id()] = info
	c.clusterSocketMap[info.SocketId] = info
	c.clusterLocker.Unlock()
	c.hashRing.Add(info.IpString())
	c.service.SendMsg(rpc.RpcHead{SocketId: info.SocketId}, "COMMON_RegisterResponse")
	switch info.Type {
	case rpc.SERVICE_GATE:
		base.LOG.Printf("ADD Gate SERVER: [%d]-[%s:%d]", info.SocketId, info.Ip, info.Port)
	}
}

func (c *ClusterServer) DelCluster(info *common.ClusterInfo) {
	c.clusterLocker.RLock()
	_, bEx := c.clusterMap[info.Id()]
	c.clusterLocker.RUnlock()
	if bEx {
		c.clusterLocker.Lock()
		delete(c.clusterMap, info.Id())
		delete(c.clusterSocketMap, info.SocketId)
		c.clusterLocker.Unlock()
	}

	c.hashRing.Remove(info.IpString())
	base.LOG.Printf("服务器断开连接socketid[%d]", info.SocketId)
	switch info.Type {
	case rpc.SERVICE_GATE:
		base.LOG.Printf("与Gate服务器断开连接,id[%d]", info.SocketId)
	}
}

func (c *ClusterServer) GetCluster(head rpc.RpcHead) *common.ClusterInfo {
	c.clusterLocker.RLock()
	defer c.clusterLocker.RUnlock()
	client, bEx := c.clusterMap[head.ClusterId]
	if bEx {
		return client
	}
	return nil
}

func (c *ClusterServer) GetClusterBySocket(socketId uint32) *common.ClusterInfo {
	c.clusterLocker.RLock()
	defer c.clusterLocker.RUnlock()
	client, bEx := c.clusterSocketMap[socketId]
	if bEx {
		return client
	}
	return nil
}

func (c *ClusterServer) BindServer(pService *network.ServerSocket) {
	c.service = pService
}

func (c *ClusterServer) sendPoint(head rpc.RpcHead, packet rpc.Packet) {
	pCluster := c.GetCluster(head)
	if pCluster != nil {
		head.SocketId = pCluster.SocketId
		c.service.Send(head, packet)
	}
}

func (c *ClusterServer) balanceSend(head rpc.RpcHead, packet rpc.Packet) {
	_, head.ClusterId = c.hashRing.Get64(head.Id)
	pCluster := c.GetCluster(head)
	if pCluster != nil {
		head.SocketId = pCluster.SocketId
		c.service.Send(head, packet)
	}
}

func (c *ClusterServer) boardCastSend(head rpc.RpcHead, packet rpc.Packet) {
	clusterList := []*common.ClusterInfo{}
	c.clusterLocker.RLock()
	for _, v := range c.clusterMap {
		clusterList = append(clusterList, v)
	}
	c.clusterLocker.RUnlock()
	for _, v := range clusterList {
		head.SocketId = v.SocketId
		c.service.Send(head, packet)
	}
}

func (c *ClusterServer) SendMsg(head rpc.RpcHead, funcName string, params ...interface{}) {
	c.Send(head, rpc.Marshal(&head, &funcName, params...))
}

func (c *ClusterServer) Send(head rpc.RpcHead, packet rpc.Packet) {
	if head.DestServerType != rpc.SERVICE_GATE {
		c.balanceSend(head, packet)
	} else {
		switch head.SendType {
		//case rpc.SEND_BALANCE:
		//	c.balanceSend(head, packet)
		case rpc.SEND_POINT:
			c.sendPoint(head, packet)
		default:
			c.boardCastSend(head, packet)
		}
	}
}

func (c *ClusterServer) RandomCluster(head rpc.RpcHead) rpc.RpcHead {
	if head.Id == 0 {
		head.Id = int64(uint32(base.RAND.RandI(1, 0xFFFFFFFF)))
	}
	_, head.ClusterId = c.hashRing.Get64(head.Id)
	pCluster := c.GetCluster(head)
	if pCluster != nil {
		head.SocketId = pCluster.SocketId
	}
	return head
}

func (c *ClusterServer) COMMON_RegisterRequest(ctx context.Context, info *common.ClusterInfo) {
	pServerInfo := info
	pServerInfo.SocketId = c.GetRpcHead(ctx).SocketId

	c.AddCluster(pServerInfo)
}

//链接断开
func (c *ClusterServer) DISCONNECT(ctx context.Context, socketId uint32) {
	pCluster := c.GetClusterBySocket(socketId)
	if pCluster != nil {
		c.DelCluster(pCluster)
	}
}
