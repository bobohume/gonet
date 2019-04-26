package cluster

import (
	"gonet/base"
	"gonet/server/common"
	"math"
	"sync"
)

type (
	//集群信息
	ClusterInfo struct {
		*common.ServerInfo
		ConnectNum int//连接数量
	}


	HashClusterMap map[int] *ClusterInfo

	ClusterManager struct{
		m_ClusterMap 	HashClusterMap
		m_ClusterList	base.IVector//保存地址，map的hash性质
		m_ClusterLocker	*sync.RWMutex
	}

	IClusterManager interface {
		Init()
		AddCluster(*common.ServerInfo)//添加集群
		DelCluster(*common.ServerInfo)//删除集群
		UpdateCluster(int, int)//更新集群信息
		GetAllCluster(base.IVector)//获取所有集群
		GetCluster(int) *ClusterInfo
		BalanceCluster()int//负载均衡
		RandomCluster()int//随机分配

		ClusterInitPacket(int) []byte//集群初始包
		ClusterAddPacket(*common.ServerInfo) []byte//集群添加member包
		ClusterDelPacket(*common.ServerInfo) []byte//集群删除membre包
	}
)


func (this *ClusterManager) Init(){
	this.m_ClusterLocker = &sync.RWMutex{}
	this.m_ClusterMap = make(HashClusterMap)
	this.m_ClusterList = &base.Vector{}
}

func (this *ClusterManager) AddCluster(pServerInfo *common.ServerInfo){
	this.m_ClusterLocker.Lock()
	pClusterInfo, bEx := this.m_ClusterMap[pServerInfo.SocketId]
	if bEx{
		pClusterInfo.ServerInfo = pServerInfo
	}else{
		pClusterInfo = &ClusterInfo{pServerInfo, 0}
		this.m_ClusterMap[pServerInfo.SocketId] =pClusterInfo
		this.m_ClusterList.Push_back(pClusterInfo)
	}
	this.m_ClusterLocker.Unlock()
}

func (this *ClusterManager) DelCluster(pServerInfo *common.ServerInfo){
	this.m_ClusterLocker.Lock()
	for i, v := range this.m_ClusterList.Array(){
		if v.(*ClusterInfo).SocketId == pServerInfo.SocketId{
			this.m_ClusterList.Erase(i)
			break
		}
	}
	delete(this.m_ClusterMap, pServerInfo.SocketId)
	this.m_ClusterLocker.Unlock()
}

func (this *ClusterManager) UpdateCluster(socketId int, Num int){
	this.m_ClusterLocker.RLock()
	pClusterInfo, bEx := this.m_ClusterMap[socketId]
	if bEx{
		pClusterInfo.ConnectNum = Num
	}
	this.m_ClusterLocker.RUnlock()
}

func (this *ClusterManager) GetCluster(socketId int) *ClusterInfo{
	this.m_ClusterLocker.RLock()
	pClusterInfo, _ := this.m_ClusterMap[socketId]
	this.m_ClusterLocker.RUnlock()
	return  pClusterInfo
}

func (this *ClusterManager) GetAllCluster(vec base.IVector){
	this.m_ClusterLocker.RLock()
	for _, v := range this.m_ClusterMap{
		vec.Push_back(v)
	}
	this.m_ClusterLocker.RUnlock()
}

func (this *ClusterManager) BalanceCluster() int{
	nIndex := -1
	this.m_ClusterLocker.RLock()
	for _, v := range this.m_ClusterList.Array(){
		pClusterInfo := v.(*ClusterInfo)
		if pClusterInfo.ConnectNum <= 10000{
			nIndex = pClusterInfo.SocketId
			break
		}
	}
	this.m_ClusterLocker.RUnlock()
	if nIndex == -1{
		nIndex = this.RandomCluster()
	}
	return nIndex
}

func (this *ClusterManager) RandomCluster() int{
	nIndex := -1
	this.m_ClusterLocker.RLock()
	if this.m_ClusterList.Len() > 0{
		nLen := int(math.Max(float64(this.m_ClusterList.Len()-1), 0))
		nRand := base.RAND().RandI(0, nLen)
		nIndex = this.m_ClusterList.Get(nRand).(*ClusterInfo).SocketId
	}
	this.m_ClusterLocker.RUnlock()
	return nIndex
}

func (this *ClusterManager) ClusterInitPacket(nType int) []byte{
	aIp := []string{}
	aPort := []int{}
	aSocketId := []int{}
	vec := base.NewVector()
	this.GetAllCluster(vec)
	for _, v := range vec.Array(){
		pInfo, bOk := v.(*ClusterInfo)
		if bOk && pInfo != nil{
			aIp = append(aIp, pInfo.Ip)
			aPort = append(aPort, pInfo.Port)
			aSocketId = append(aSocketId, pInfo.SocketId)
		}
	}
	return base.GetPacket(GetClusterMsg("Cluster_Socket_Init", nType), aIp, aPort, aSocketId)
}

func (this *ClusterManager) ClusterAddPacket(pServerInfo *common.ServerInfo) []byte{
	return base.GetPacket(GetClusterMsg("Cluster_Socket_Add", pServerInfo.Type), pServerInfo.Ip, pServerInfo.Port, pServerInfo.SocketId)
}

func (this *ClusterManager) ClusterDelPacket(pServerInfo *common.ServerInfo) []byte{
	return base.GetPacket(GetClusterMsg("Cluster_Socket_Del", pServerInfo.Type), pServerInfo.Ip, pServerInfo.Port, pServerInfo.SocketId)
}