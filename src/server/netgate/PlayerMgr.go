package netgate

import (
	"time"
	"base"
	"actor"
	"sync"
)

type(
	CPlayerManager struct {
		actor.Actor
		m_SocketMap map[int] int
		m_AccountMap map[int] AccountInfo
		m_Locker sync.Locker
	}

	IPlayerMangaer interface {
		actor.IActor
		ReleaseSocketMap(int)
		AddAccountMap(int, int, int) int
		GetAccountSocket(int) int
		GetSocketAccount(int) (int, int)
		getSocket(int) (int,bool)
		getAccount(int) (AccountInfo,bool)
	}

	AccountInfo struct{
		LastTime int64
		SocketId int
		AccountName string
	}
)

func NewAccountInfo(socket int) AccountInfo{
	accountInfo := AccountInfo{time.Now().Unix(), socket, ""}
	return  accountInfo
}

func (this *CPlayerManager)getSocket(socketId int) (int,bool){
	this.m_Locker.Lock()
	accountid, exist :=  this.m_SocketMap[socketId]
	this.m_Locker.Unlock()
	return accountid, exist
}

func (this *CPlayerManager)getAccount(accountId int) (AccountInfo,bool){
	this.m_Locker.Lock()
	accountInfo, exist := this.m_AccountMap[accountId]
	this.m_Locker.Unlock()
	return accountInfo, exist
}

func (this *CPlayerManager)ReleaseSocketMap(socketId int, bClose bool){
	accountid,exist := this.getSocket(socketId)
	if exist == false{
		return
	}
	this.m_Locker.Lock()
	delete(this.m_AccountMap, accountid)
	delete(this.m_SocketMap, socketId)
	this.m_Locker.Unlock()
	if bClose == true{
		SERVER.GetServer().StopClient(socketId)
	}
}

func (this *CPlayerManager)AddAccountMap(socketId int, accountId int) int {
	accountInfo, exist := this.getAccount(accountId)
	if exist == true{
		this.ReleaseSocketMap(accountInfo.SocketId, accountInfo.SocketId != socketId)
	}
	this.m_Locker.Lock()
	this.m_AccountMap[accountId] = NewAccountInfo(socketId)
	this.m_SocketMap[socketId] = accountId
	this.m_Locker.Unlock()
	return  base.NONE_ERROR
}

func (this *CPlayerManager)GetAccountSocket(accountId int) int{
	accountInfo, exist := this.getAccount(accountId)
	if exist == true{
		return accountInfo.SocketId
	}
	return  0
}

func (this *CPlayerManager)GetSocketAccount(socketId int) (int){
	accountid, exist :=  this.getSocket(socketId)
	if exist == true{
		_, exist := this.getAccount(accountid)
		if exist == true{
			return accountid
		}
	}
	return   0
}

func (this *CPlayerManager) Init(num int){
	this.Actor.Init(num)
	this.m_SocketMap = make(map[int] int)
	this.m_AccountMap = make(map[int] AccountInfo)
	this.m_Locker = &sync.Mutex{}

	this.RegisterCall("ADD_ACCOUNT", func(caller *actor.Caller, socketid int, accountId int) {
		SERVER.GetLog().Printf("login incoming  Socket:%d Account:%d ",socketid, accountId)
		error := this.AddAccountMap(socketid, accountId)
		if error == base.NONE_ERROR{
			SendToWorld("G_W_CLoginRequest", accountId)
		}
		SERVER.GetLog().Printf("login incoming  Socket:%d Account:%d ",socketid, accountId)
	})

	this.RegisterCall("DEL_ACCOUNT", func(caller *actor.Caller, socketid int) {
		accountId := this.GetSocketAccount(socketid)
		this.ReleaseSocketMap(socketid, true)
		SERVER.GetWorldScoket().SendMsg("G_ClientLost", accountId)
		SERVER.GetAccountScoket().SendMsg("G_ClientLost", accountId)
	})

	//重连世界服务器，账号重新登录
	this.RegisterCall("Account_Relink", func(caller *actor.Caller) {
		accountMap := make(map [int] int)
		this.m_Locker.Lock()
		for i, v := range this.m_SocketMap {
			accountMap[i] = v
		}
		this.m_Locker.Unlock()

		if len(accountMap) != 0{
			for _, v := range accountMap {
				SendToWorld("G_W_CLoginRequest", v)
			}
		}
	})
	this.Actor.Start()
}