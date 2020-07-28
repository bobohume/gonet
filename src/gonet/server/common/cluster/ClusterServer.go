package cluster

import (
	"context"
	"gonet/actor"
	"gonet/base"
	"gonet/message"
	"gonet/network"
	"gonet/rpc"
	"gonet/server/common"
	"math"
	"sync"
)

const(
	MAX_CLUSTER_TYPE = int(message.SERVICE_WORLDDBSERVER)
	DISPATCH_CLUSTER = message.SERVICE_GATESERVER//转发服
)

type(
	HashClusterMap map[uint32] *common.ClusterInfo
	HashClusterSocketMap map[uint32] *common.ClusterInfo

	//集群服务器
	ClusterServer struct{
		actor.Actor
		*Service //集群注册
		m_ClusterMap [MAX_CLUSTER_TYPE]HashClusterMap
		m_ClusterSocketMap [MAX_CLUSTER_TYPE]HashClusterSocketMap
		m_ClusterList [MAX_CLUSTER_TYPE]base.IVector
		m_ClusterLocker *sync.RWMutex
		m_pService *network.ServerSocket//socket管理
	}

	IClusterServer interface{
		InitService(Type message.SERVICE, IP string, Port int, Endpoints []string)
		RegisterClusterCall()//注册集群通用回调
		AddCluster(info *common.ClusterInfo)
		DelCluster(info *common.ClusterInfo)
		GetCluster(rpc.RpcHead) *common.ClusterInfo
		GetClusterBySocket(uint32) *common.ClusterInfo

		BindServer(*network.ServerSocket)
		SendMsg(rpc.RpcHead, string, ...interface{})//发送给集群特定服务器
		Send(rpc.RpcHead, []byte)//发送给集群特定服务器

		BalanceCluster(message.SERVICE)	rpc.RpcHead//负载均衡
		RandomCluster(message.SERVICE)	rpc.RpcHead//随机分配

		sendPoint(rpc.RpcHead, []byte)//发送给集群特定服务器
		balanceSend(rpc.RpcHead, []byte)//负载给集群特定服务器
		boardCastSend(head rpc.RpcHead, buff []byte)//给集群广播
	}
)

func (this *ClusterServer) InitService(Type message.SERVICE, IP string, Port int, Endpoints []string) {
	this.m_ClusterLocker = &sync.RWMutex{}
	//注册服务器
	this.Service = NewService(Type, IP, Port, Endpoints)
	for i := 0; i < MAX_CLUSTER_TYPE; i++{
		this.m_ClusterMap[i] = make(HashClusterMap)
		this.m_ClusterSocketMap[i] = make(HashClusterSocketMap)
		this.m_ClusterList[i]  = &base.Vector{}
	}
}

func (this *ClusterServer) RegisterClusterCall(){
	this.RegisterCall("COMMON_RegisterRequest", func(ctx context.Context, info *message.ClusterInfo) {
		pServerInfo := &common.ClusterInfo{ClusterInfo:*info}
		pServerInfo.SocketId = this.GetRpcHead(ctx).SocketId

		this.AddCluster(pServerInfo)
	})

	//链接断开
	this.RegisterCall("DISCONNECT", func(ctx context.Context, socketId uint32) {
		pCluster := this.GetClusterBySocket(socketId)
		if pCluster != nil{
			this.DelCluster(pCluster)
		}
	})
}

func (this *ClusterServer) AddCluster(info *common.ClusterInfo){
	this.m_ClusterLocker.Lock()
	this.m_ClusterMap[info.Type][info.Id()] = info
	this.m_ClusterSocketMap[info.Type][info.SocketId] = info
	this.m_ClusterList[info.Type].Push_back(info)
	this.m_ClusterLocker.Unlock()
	this.m_pService.SendMsg(rpc.RpcHead{SocketId:info.SocketId}, "COMMON_RegisterResponse")
	switch info.Type {
	case message.SERVICE_GATESERVER:
		base.GLOG.Printf("ADD GATE SERVER: [%d]-[%s:%d]", info.SocketId, info.Ip, info.Port)
	case message.SERVICE_WORLDSERVER:
		base.GLOG.Printf("ADD WORLD SERVER: [%d]-[%s:%d]", info.SocketId, info.Ip, info.Port)
	case message.SERVICE_ZONESERVER:
		base.GLOG.Printf("ADD ZONE SERVER: [%d]-[%s:%d]", info.SocketId, info.Ip, info.Port)
	}
}

func (this *ClusterServer) DelCluster(info *common.ClusterInfo){
	this.m_ClusterLocker.RLock()
	_, bEx := this.m_ClusterMap[info.Type][info.Id()]
	this.m_ClusterLocker.RUnlock()
	if bEx{
		this.m_ClusterLocker.Lock()
		delete(this.m_ClusterMap[info.Type], info.Id())
		delete(this.m_ClusterSocketMap[info.Type], info.SocketId)
		for i, v := range this.m_ClusterList[info.Type].Array(){
			if v.(*common.ClusterInfo).SocketId == info.SocketId{
				this.m_ClusterList[info.Type].Erase(i)
				break
			}
		}
		this.m_ClusterLocker.Unlock()
	}

	base.GLOG.Printf("服务器断开连接socketid[%d]",info.SocketId)
	switch info.Type {
	case message.SERVICE_GATESERVER:
		base.GLOG.Printf("与Gate服务器断开连接,id[%d]",info.SocketId)
	case message.SERVICE_WORLDSERVER:
		base.GLOG.Printf("与World服务器断开连接,id[%d]",info.SocketId)
	case message.SERVICE_ZONESERVER:
		base.GLOG.Printf("与Zone服务器断开连接,id[%d]",info.SocketId)
	}
}

func (this *ClusterServer) GetCluster(head rpc.RpcHead) *common.ClusterInfo{
	this.m_ClusterLocker.RLock()
	defer this.m_ClusterLocker.RUnlock()
	pClient, bEx := this.m_ClusterMap[DISPATCH_CLUSTER][head.ClusterId]
	if bEx{
		return pClient
	}
	return nil
}

func (this *ClusterServer) GetClusterBySocket(socketId uint32) *common.ClusterInfo{
	this.m_ClusterLocker.RLock()
	defer this.m_ClusterLocker.RUnlock()
	for i := 0; i < MAX_CLUSTER_TYPE; i++ {
		pClient, bEx := this.m_ClusterSocketMap[i][socketId]
		if bEx{
			return pClient
		}
	}
	return nil
}

func (this *ClusterServer) BindServer(pService *network.ServerSocket){
	this.m_pService = pService
}

func (this *ClusterServer) sendPoint(head rpc.RpcHead, buff []byte){
	pCluster:= this.GetCluster(head)
	if pCluster != nil {
		head.SocketId = pCluster.SocketId
		this.m_pService.Send(head, buff)
	}
}

func (this *ClusterServer)balanceSend(head rpc.RpcHead, buff []byte){
	head = this.RandomCluster(DISPATCH_CLUSTER)
	pCluster := this.GetCluster(head)
	if pCluster != nil{
		this.m_pService.Send(head, buff)
	}
}

func (this *ClusterServer) boardCastSend(head rpc.RpcHead, buff []byte){
	clusterList := []*common.ClusterInfo{}
	this.m_ClusterLocker.RLock()
	for _ ,v := range this.m_ClusterMap[DISPATCH_CLUSTER]{
		clusterList = append(clusterList, v)
	}
	this.m_ClusterLocker.RUnlock()
	for _, v := range clusterList{
		head.SocketId = v.SocketId
		this.m_pService.Send(head, buff)
	}
}

func (this *ClusterServer) SendMsg(head rpc.RpcHead, funcName string, params  ...interface{}){
	buff := base.SetTcpEnd(rpc.Marshal(head, funcName, params...))
	this.Send(head, buff)
}

func (this *ClusterServer) Send(head rpc.RpcHead, buff []byte){
	if head.DestServerType != DISPATCH_CLUSTER{
		this.balanceSend(head, buff)
	}else{
		switch head.SendType{
		case message.SEND_BALANCE:
			this.balanceSend(head, buff)
		case message.SEND_POINT:
			this.sendPoint(head, buff)
		default:
			this.boardCastSend(head, buff)
		}
	}
}

func (this *ClusterServer) BalanceCluster(nType message.SERVICE) rpc.RpcHead{
	head := rpc.RpcHead{}
	this.m_ClusterLocker.RLock()
	for _, v := range this.m_ClusterList[nType].Array(){
		pClusterInfo := v.(*common.ClusterInfo)
		if pClusterInfo.Weight <= 10000{
			head.ClusterId = pClusterInfo.Id()
			head.SocketId = pClusterInfo.SocketId
			break
		}
	}
	this.m_ClusterLocker.RUnlock()
	if head.ClusterId == 0{
		return this.RandomCluster(nType)
	}
	return head
}

func (this *ClusterServer) RandomCluster(nType message.SERVICE) rpc.RpcHead{
	head := rpc.RpcHead{}
	this.m_ClusterLocker.RLock()
	if this.m_ClusterList[nType].Len() > 0{
		nLen := int(math.Max(float64(this.m_ClusterList[nType].Len()-1), 0))
		nRand := base.RAND.RandI(0, nLen)
		info := this.m_ClusterList[nType].Get(nRand).(*common.ClusterInfo)
		head.ClusterId = info.Id()
		head.SocketId = info.SocketId
	}
	this.m_ClusterLocker.RUnlock()
	return head
}