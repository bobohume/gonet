package netgate

import (
	"gonet/actor"
	"gonet/base"
	"sync"
	"time"
)

type(
	PlayerManager struct {
		actor.Actor
		m_SocketMap map[int] int64
		m_AccountMap map[int64] *AccountInfo
		m_Locker *sync.RWMutex
	}

	IPlayerMangaer interface {
		actor.IActor
		ReleaseSocketMap(int, bool)
		AddAccountMap(int, int64) int
		GetSocket(int64) int
		GetAccount(int) int64
		GetAccountInfo(int) *AccountInfo
	}

	AccountInfo struct{
		AccountId int64
		LastTime int64
		SocketId int
		WClusterId uint32
	}
)

var(
	g_pAccount = &AccountInfo{}
)

func NewAccountInfo(socket int, accountId int64) *AccountInfo{
	accountInfo := AccountInfo{LastTime:time.Now().Unix(), SocketId:socket, WClusterId:0, AccountId:accountId}
	return  &accountInfo
}

func (this *PlayerManager) ReleaseSocketMap(socketId int, bClose bool){
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

func (this *PlayerManager) AddAccountMap(accountId int64, socketId int) int {
	Id := this.GetSocket(accountId)
	this.ReleaseSocketMap(Id, Id != socketId)

	accountInfo := NewAccountInfo(socketId, accountId)
	accountInfo.WClusterId = SERVER.GetWorldCluster().RandomCluster()
	this.m_Locker.Lock()
	this.m_AccountMap[accountId] = accountInfo
	this.m_SocketMap[socketId] = accountId
	this.m_Locker.Unlock()
	SERVER.GetWorldCluster().SendMsg(accountInfo.WClusterId, "G_W_CLoginRequest", accountId, SERVER.GetCluster().Id())
	return  base.NONE_ERROR
}

func (this *PlayerManager) GetSocket(accountId int64) int{
	socketId := 0
	this.m_Locker.RLock()
	accountInfo, exist := this.m_AccountMap[accountId]
	this.m_Locker.RUnlock()
	if exist{
		socketId = accountInfo.SocketId
	}
	return socketId
}

func (this *PlayerManager) GetAccount(socketId int) int64{
	accoundId := int64(0)
	this.m_Locker.RLock()
	id, exist :=  this.m_SocketMap[socketId]
	this.m_Locker.RUnlock()
	if exist{
		accoundId = id
	}
	return accoundId
}

func (this *PlayerManager) GetAccountInfo(socketId int) *AccountInfo{
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
	this.m_SocketMap = make(map[int] int64)
	this.m_AccountMap = make(map[int64] *AccountInfo)
	this.m_Locker = &sync.RWMutex{}
	this.RegisterCall("ADD_ACCOUNT", func(accountId int64, socketId int) {
		SERVER.GetLog().Printf("login incoming  Socket:%d Account:%d ",socketId, accountId)
		this.AddAccountMap(accountId, socketId)
	})

	this.RegisterCall("DEL_ACCOUNT", func(socketid int) {
		accountId := this.GetAccount(socketid)
		this.ReleaseSocketMap(socketid, true)
		SERVER.GetAccountCluster().BoardCastMsg("G_ClientLost", accountId)
	})

	//重连世界服务器，账号重新登录
	/*this.RegisterCall("World_Relogin", func() {
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