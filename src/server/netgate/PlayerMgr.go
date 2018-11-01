package netgate

import (
	"time"
	"base"
	"actor"
	"sync"
)

type(
	PlayerManager struct {
		actor.Actor
		m_SocketMap map[int] int
		m_AccountMap map[int] *AccountInfo
		m_Locker *sync.RWMutex
	}

	IPlayerMangaer interface {
		actor.IActor
		ReleaseSocketMap(int, bool)
		AddAccountMap(int, int) int
		GetAccountSocket(int) int
		GetSocketAccount(int) int
	}

	AccountInfo struct{
		LastTime int64
		SocketId int
		AccountName string
	}
)

func NewAccountInfo(socket int) *AccountInfo{
	accountInfo := AccountInfo{time.Now().Unix(), socket, ""}
	return  &accountInfo
}

func (this *PlayerManager)ReleaseSocketMap(socketId int, bClose bool){
	this.m_Locker.RLock()
	accountid, _ :=  this.m_SocketMap[socketId]
	this.m_Locker.RUnlock()
	this.m_Locker.Lock()
	delete(this.m_AccountMap, accountid)
	delete(this.m_SocketMap, socketId)
	this.m_Locker.Unlock()
	//if bClose{
	SERVER.GetServer().StopClient(socketId)
	//}
}

func (this *PlayerManager)AddAccountMap(socketId int, accountId int) int {
	Id := this.getAccountId(accountId)
	this.ReleaseSocketMap(Id, Id != socketId)

	accountInfo := NewAccountInfo(socketId)
	this.m_Locker.Lock()
	this.m_AccountMap[accountId] = accountInfo
	this.m_SocketMap[socketId] = accountId
	this.m_Locker.Unlock()
	return  base.NONE_ERROR
}

func (this *PlayerManager)GetAccountSocket(accountId int) int{
	return this.getAccountId(accountId)
}

func (this *PlayerManager) GetSocketAccount(socketId int) int{
	return this.getSocketId(socketId)
}

func (this *PlayerManager)getSocketId(socketId int) (int){
	accoundId := 0
	this.m_Locker.RLock()
	id, exist :=  this.m_SocketMap[socketId]
	this.m_Locker.RUnlock()
	if exist{
		accoundId = id
	}
	return accoundId
}

func (this *PlayerManager)getAccountId(accountId int) (int){
	socketId := 0
	this.m_Locker.RLock()
	accountInfo, exist := this.m_AccountMap[accountId]
	this.m_Locker.RUnlock()
	if exist{
		socketId = accountInfo.SocketId
	}
	return socketId
}

func (this *PlayerManager) Init(num int){
	this.Actor.Init(num)
	this.m_SocketMap = make(map[int] int)
	this.m_AccountMap = make(map[int] *AccountInfo)
	this.m_Locker = &sync.RWMutex{}

	this.RegisterCall("ADD_ACCOUNT", func(socketid int, accountId int) {
		SERVER.GetLog().Printf("login incoming  Socket:%d Account:%d ",socketid, accountId)
		error := this.AddAccountMap(socketid, accountId)
		if error == base.NONE_ERROR{
			SendToWorld("G_W_CLoginRequest", accountId)
		}
	})

	this.RegisterCall("DEL_ACCOUNT", func(socketid int) {
		accountId := this.GetSocketAccount(socketid)
		this.ReleaseSocketMap(socketid, true)
		SERVER.GetWorldSocket().SendMsg("G_ClientLost", accountId)
		SERVER.GetAccountSocket().SendMsg("G_ClientLost", accountId)
	})

	//重连世界服务器，账号重新登录
	this.RegisterCall("Account_Relink", func() {
		accountMap := make(map [int] int)
		this.m_Locker.RLock()
		for i, v := range this.m_SocketMap {
			accountMap[i] = v
		}
		this.m_Locker.RUnlock()

		if len(accountMap) != 0{
			for _, v := range accountMap {
				SendToWorld("G_W_CLoginRequest", v)
			}
		}
	})
	this.Actor.Start()
}