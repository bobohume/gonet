package netgate

import (
	"context"
	"gonet/actor"
	"gonet/base"
	"gonet/rpc"
	"sync"
	"time"
)

type (
	PlayerMgr struct {
		actor.Actor
		m_SocketMap  map[uint32]int64
		m_AccountMap map[int64]*AccountInfo
		m_Locker     *sync.RWMutex
	}

	IPlayerMgr interface {
		actor.IActor
		ReleaseSocketMap(uint32, bool)
		AddAccountMap(uint32, rpc.PlayerClusterInfo) int
		GetSocket(int64) uint32
		GetAccount(uint32) int64
		GetAccountInfo(uint32) *AccountInfo
	}

	AccountInfo struct {
		AccountId  int64
		LastTime   int64
		SocketId   uint32
		WClusterId uint32
		ZClusterId uint32
	}
)

var (
	g_pAccount = &AccountInfo{}
)

func NewAccountInfo(socket uint32, accountId int64) *AccountInfo {
	accountInfo := AccountInfo{LastTime: time.Now().Unix(), SocketId: socket, WClusterId: 0, AccountId: accountId, ZClusterId: 0}
	return &accountInfo
}

func (this *PlayerMgr) ReleaseSocketMap(socketId uint32, bClose bool) {
	this.m_Locker.RLock()
	accountId, _ := this.m_SocketMap[socketId]
	this.m_Locker.RUnlock()
	this.m_Locker.Lock()
	delete(this.m_AccountMap, accountId)
	delete(this.m_SocketMap, socketId)
	this.m_Locker.Unlock()
	//if bClose{
	SERVER.GetServer().StopClient(socketId)
	//}
}

func (this *PlayerMgr) AddAccountMap(socketId uint32, clusterInfo rpc.PlayerClusterInfo) int {
	accountId := clusterInfo.Id
	Id := this.GetSocket(accountId)
	this.ReleaseSocketMap(Id, Id != socketId)

	accountInfo := NewAccountInfo(socketId, accountId)
	accountInfo.WClusterId = clusterInfo.WClusterId
	accountInfo.ZClusterId = clusterInfo.ZClusterId
	this.m_Locker.Lock()
	this.m_AccountMap[accountId] = accountInfo
	this.m_SocketMap[socketId] = accountId
	this.m_Locker.Unlock()
	SERVER.GetCluster().SendMsg(rpc.RpcHead{ClusterId: accountInfo.WClusterId, DestServerType: rpc.SERVICE_WORLDSERVER}, "G_W_CLoginRequest", accountId, SERVER.GetCluster().Id(), clusterInfo)
	return base.NONE_ERROR
}

func (this *PlayerMgr) GetSocket(accountId int64) uint32 {
	socketId := uint32(0)
	this.m_Locker.RLock()
	accountInfo, exist := this.m_AccountMap[accountId]
	this.m_Locker.RUnlock()
	if exist {
		socketId = accountInfo.SocketId
	}
	return socketId
}

func (this *PlayerMgr) GetAccount(socketId uint32) int64 {
	accoundId := int64(0)
	this.m_Locker.RLock()
	id, exist := this.m_SocketMap[socketId]
	this.m_Locker.RUnlock()
	if exist {
		accoundId = id
	}
	return accoundId
}

func (this *PlayerMgr) GetAccountInfo(socketId uint32) *AccountInfo {
	accountId := this.GetAccount(socketId)
	this.m_Locker.RLock()
	accountInfo, exist := this.m_AccountMap[accountId]
	this.m_Locker.RUnlock()
	if exist {
		return accountInfo
	}
	return nil
}

func (this *PlayerMgr) Init() {
	this.Actor.Init()
	this.m_SocketMap = make(map[uint32]int64)
	this.m_AccountMap = make(map[int64]*AccountInfo)
	this.m_Locker = &sync.RWMutex{}
	actor.MGR.RegisterActor(this)
	this.Actor.Start()
}

func (this *PlayerMgr) ADD_ACCOUNT(ctx context.Context, socketId uint32, clusterInfo rpc.PlayerClusterInfo) {
	SERVER.GetLog().Printf("login incoming  Socket:%d Account:%d WClusterId:%d ", socketId, clusterInfo.Id, clusterInfo.WClusterId)
	this.AddAccountMap(socketId, clusterInfo)
}

func (this *PlayerMgr) DEL_ACCOUNT(ctx context.Context, socketid uint32) {
	accountId := this.GetAccount(socketid)
	this.ReleaseSocketMap(socketid, true)
	SERVER.GetCluster().SendMsg(rpc.RpcHead{SendType: rpc.SEND_BOARD_CAST, DestServerType: rpc.SERVICE_WORLDSERVER}, "G_ClientLost", accountId)
}

//重连世界服务器，账号重新登录
func (this *PlayerMgr) World_Relogin(ctx context.Context) {
	accountMap := make(map [int64] uint32)
	this.m_Locker.RLock()
	for i, v := range this.m_AccountMap {
		accountMap[i] = v.WClusterId
	}
	this.m_Locker.RUnlock()

	if len(accountMap) != 0{
		for i, v := range accountMap {
			SERVER.GetCluster().SendMsg(rpc.RpcHead{Id:i, ClusterId:v}, "G_W_Relogin")
		}
	}
}