package cluster

import (
	"gonet/actor"
	"gonet/base"
	"gonet/message"
	"gonet/network"
	"gonet/server/common"
	"math"
	"sync"
)

const(
	MAX_CLUSTER_TYPE = int(message.SERVICE_WORLDDBSERVER)
)

type(
	HashClusterMap map[int] *common.ServerInfo

	//集群服务器
	ClusterServer struct{
		actor.Actor

		m_ClusterMap [MAX_CLUSTER_TYPE]HashClusterMap
		m_ClusterList [MAX_CLUSTER_TYPE]base.IVector
		m_ClusterLocker *sync.RWMutex
		m_Service  *Service//集群注册
		m_pService *network.ServerSocket//socket管理
	}

	IClusterServer interface{
		InitService(Type int, IP string, Port int, Endpoints []string)
		RegisterClusterCall()//注册集群通用回调
		AddCluster(info *common.ServerInfo)
		DelCluster(info *common.ServerInfo)
		GetCluster(int) *common.ServerInfo

		BindServer(*network.ServerSocket)
		SendMsg(int, string, ...interface{})//发送给集群特定服务器
		BalanceMsg(int, string, ...interface{})//负载给集群特定服务器
		BoardCastMsg(int, string, ...interface{})//给集群广播
		Send(int, []byte)
		BalanceSend(int, []byte)//负载给集群特定服务器

		BalanceCluster(int)int//负载均衡
		RandomCluster(int)int//随机分配
	}
)

func (this *ClusterServer) InitService(Type int, IP string, Port int, Endpoints []string) {
	this.m_ClusterLocker = &sync.RWMutex{}
	//注册服务器
	this.m_Service = NewService(Type, IP, Port, Endpoints)
	for i := 0; i < MAX_CLUSTER_TYPE; i++{
		this.m_ClusterMap[i] = make(HashClusterMap)
		this.m_ClusterList[i]  = &base.Vector{}
	}
}

func (this *ClusterServer) RegisterClusterCall(){
	this.RegisterCall("COMMON_RegisterRequest", func(nType int, Ip string, Port int) {
		pServerInfo := new(common.ServerInfo)
		pServerInfo.SocketId = this.GetSocketId()
		pServerInfo.Type = nType
		pServerInfo.Ip = Ip
		pServerInfo.Port = Port

		this.AddCluster(pServerInfo)
	})

	//链接断开
	this.RegisterCall("DISCONNECT", func(socketId int) {
		pCluster := this.GetCluster(socketId)
		if pCluster != nil{
			this.DelCluster(pCluster)
		}
	})
}

func (this *ClusterServer) AddCluster(info *common.ServerInfo){
	this.m_ClusterLocker.Lock()
	this.m_ClusterMap[info.Type][info.SocketId] = info
	this.m_ClusterList[info.Type].Push_back(info)
	this.m_ClusterLocker.Unlock()
	this.m_pService.SendMsgById(info.SocketId, "COMMON_RegisterResponse")
	switch info.Type {
	case int(message.SERVICE_GATESERVER):
		base.GLOG.Printf("ADD GATE SERVER: [%d]-[%s:%d]", info.SocketId, info.Ip, info.Port)
	case int(message.SERVICE_WORLDSERVER):
		base.GLOG.Printf("ADD WORLD SERVER: [%d]-[%s:%d]", info.SocketId, info.Ip, info.Port)
	}
}

func (this *ClusterServer) DelCluster(info *common.ServerInfo){
	this.m_ClusterLocker.RLock()
	_, bEx := this.m_ClusterMap[info.Type][info.SocketId]
	this.m_ClusterLocker.RUnlock()
	if bEx{
		this.m_ClusterLocker.Lock()
		delete(this.m_ClusterMap[info.Type], info.SocketId)
		for i, v := range this.m_ClusterList[info.Type].Array(){
			if v.(*common.ServerInfo).SocketId == info.SocketId{
				this.m_ClusterList[info.Type].Erase(i)
				break
			}
		}
		this.m_ClusterLocker.Unlock()
	}

	base.GLOG.Printf("服务器断开连接socketid[%d]",info.SocketId)
	switch info.Type {
	case int(message.SERVICE_GATESERVER):
		base.GLOG.Printf("与Gate服务器断开连接,id[%d]",info.SocketId)
	case int(message.SERVICE_WORLDSERVER):
		base.GLOG.Printf("与World服务器断开连接,id[%d]",info.SocketId)
	}
}

func (this *ClusterServer) GetCluster(socketId int) *common.ServerInfo{
	this.m_ClusterLocker.RLock()
	defer this.m_ClusterLocker.RUnlock()
	for i := 0; i < MAX_CLUSTER_TYPE; i++ {
		pClient, bEx := this.m_ClusterMap[i][socketId]
		if bEx{
			return pClient
		}
	}
	return nil
}

func (this *ClusterServer) BindServer(pService *network.ServerSocket){
	this.m_pService = pService
}

func (this *ClusterServer) SendMsg(socketId int, funcName string, params  ...interface{}){
	this.m_pService.SendMsgById(socketId, funcName, params...)
}

func (this *ClusterServer) BalanceMsg(nType int, funcName string, params  ...interface{}){
	pCluster:= this.GetCluster(this.RandomCluster(nType))
	if pCluster != nil{
		this.m_pService.SendMsgById(pCluster.SocketId, funcName, params...)
	}
}

func (this *ClusterServer) BoardCastMsg(nType int, funcName string, params  ...interface{}){
	clusterList := []*common.ServerInfo{}
	this.m_ClusterLocker.RLock()
	for _ ,v := range this.m_ClusterMap[nType]{
		clusterList = append(clusterList, v)
	}
	this.m_ClusterLocker.RUnlock()
	for _, v := range clusterList{
		this.m_pService.SendMsgById(v.SocketId, funcName, params...)
	}
}

func (this *ClusterServer) Send(socketId int, buff []byte){
	this.m_pService.SendById(socketId, buff)
}

func (this *ClusterServer) BalanceSend(nType int, buff []byte){
	pCluster := this.GetCluster(this.RandomCluster(nType))
	if pCluster != nil{
		this.m_pService.SendById(pCluster.SocketId, buff)
	}
}

func (this *ClusterServer) BalanceCluster(nType int) int{
	nIndex := 0
	this.m_ClusterLocker.RLock()
	for _, v := range this.m_ClusterList[nType].Array(){
		pClusterInfo := v.(*common.ServerInfo)
		if pClusterInfo.Weight <= 10000{
			nIndex = pClusterInfo.SocketId
			break
		}
	}
	this.m_ClusterLocker.RUnlock()
	if nIndex == 0{
		nIndex = this.RandomCluster(nType)
	}
	return nIndex
}

func (this *ClusterServer) RandomCluster(nType int) int{
	nIndex := 0
	this.m_ClusterLocker.RLock()
	if this.m_ClusterList[nType].Len() > 0{
		nLen := int(math.Max(float64(this.m_ClusterList[nType].Len()-1), 0))
		nRand := base.RAND.RandI(0, nLen)
		nIndex = this.m_ClusterList[nType].Get(nRand).(*common.ServerInfo).SocketId
	}
	this.m_ClusterLocker.RUnlock()
	return nIndex
}