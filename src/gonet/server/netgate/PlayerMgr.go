package netgate

import (
	"context"
	"gonet/actor"
	"gonet/base"
	"gonet/message"
	"gonet/rpc"
	"sync"
	"time"
)

type(
	PlayerManager struct {
		actor.Actor
		m_SocketMap map[uint32] int64
		m_AccountMap map[int64] *AccountInfo
		m_Locker *sync.RWMutex
	}

	IPlayerMangaer interface {
		actor.IActor
		ReleaseSocketMap(uint32, bool)
		AddAccountMap(int64, uint32) int
		GetSocket(int64) uint32
		GetAccount(uint32) int64
		GetAccountInfo(uint32) *AccountInfo
	}

	AccountInfo struct{
		AccountId int64
		LastTime int64
		SocketId uint32
		WClusterId uint32
		ZClusterId uint32
	}
)

var(
	g_pAccount = &AccountInfo{}
)

func NewAccountInfo(socket uint32, accountId int64) *AccountInfo{
	accountInfo := AccountInfo{LastTime:time.Now().Unix(), SocketId:socket, WClusterId:0, AccountId:accountId, ZClusterId:0}
	return  &accountInfo
}

func (this *PlayerManager) ReleaseSocketMap(socketId uint32, bClose bool){
	this.m_Locker.RLock()
	accountId, _ :=  this.m_SocketMap[socketId]
	this.m_Locker.RUnlock()
	this.m_Locker.Lock()
	delete(this.m_AccountMap, accountId)
	delete(this.m_SocketMap, socketId)
	this.m_Locker.Unlock()
	//if bClose{
	SERVER.GetServer().StopClient(socketId)
	//}
}

func (this *PlayerManager) AddAccountMap(accountId int64, socketId uint32) int {
	Id := this.GetSocket(accountId)
	this.ReleaseSocketMap(Id, Id != socketId)

	accountInfo := NewAccountInfo(socketId, accountId)
	accountInfo.WClusterId = SERVER.GetWorldCluster().RandomCluster(rpc.RpcHead{Id:accountId}).ClusterId
	accountInfo.ZClusterId = SERVER.GetZoneCluster().RandomCluster(rpc.RpcHead{Id:accountId}).ClusterId
	this.m_Locker.Lock()
	this.m_AccountMap[accountId] = accountInfo
	this.m_SocketMap[socketId] = accountId
	this.m_Locker.Unlock()
	SERVER.GetWorldCluster().SendMsg(rpc.RpcHead{ClusterId:accountInfo.WClusterId}, "G_W_CLoginRequest", accountId, SERVER.GetCluster().Id(), accountInfo.ZClusterId)
	return  base.NONE_ERROR
}

func (this *PlayerManager) GetSocket(accountId int64) uint32{
	socketId := uint32(0)
	this.m_Locker.RLock()
	accountInfo, exist := this.m_AccountMap[accountId]
	this.m_Locker.RUnlock()
	if exist{
		socketId = accountInfo.SocketId
	}
	return socketId
}

func (this *PlayerManager) GetAccount(socketId uint32) int64{
	accoundId := int64(0)
	this.m_Locker.RLock()
	id, exist :=  this.m_SocketMap[socketId]
	this.m_Locker.RUnlock()
	if exist{
		accoundId = id
	}
	return accoundId
}

func (this *PlayerManager) GetAccountInfo(socketId uint32) *AccountInfo{
	accountId := this.GetAccount(socketId)
	this.m_Locker.RLock()
	accountInfo, exist := this.m_AccountMap[accountId]
	this.m_Locker.RUnlock()
	if exist{
		return accountInfo
	}
	return nil
}

func (this *PlayerManager) Init(num int){
	this.Actor.Init(num)
	this.m_SocketMap = make(map[uint32] int64)
	this.m_AccountMap = make(map[int64] *AccountInfo)
	this.m_Locker = &sync.RWMutex{}
	this.RegisterCall("ADD_ACCOUNT", func(ctx context.Context, accountId int64, socketId uint32) {
		SERVER.GetLog().Printf("login incoming  Socket:%d Account:%d ",socketId, accountId)
		this.AddAccountMap(accountId, socketId)
	})

	this.RegisterCall("DEL_ACCOUNT", func(ctx context.Context, socketid uint32) {
		accountId := this.GetAccount(socketid)
		this.ReleaseSocketMap(socketid, true)
		SERVER.GetAccountCluster().SendMsg(rpc.RpcHead{SendType:message.SEND_BOARD_CAST}, "G_ClientLost", accountId)
	})

	//重连世界服务器，账号重新登录
	/*this.RegisterCall("World_Relogin", func(ctx context.Context) {
		accountMap := make(map [int64] uint32)
		this.m_Locker.RLock()
		for i, v := range this.m_AccountMap {
			accountMap[i] = v.WClusterId
		}
		this.m_Locker.RUnlock()

		if len(accountMap) != 0{
			for i, v := range accountMap {
				SERVER.GetWorldCluster().SendMsg(v, "G_W_Relogin", &message.RpcHead{Id:i})
			}
		}
	})*/
	this.Actor.Start()
}