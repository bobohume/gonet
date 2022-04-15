package gate

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
		m_SocketMap map[uint32]int64
		m_PlayerMap map[int64]*Player
		m_Locker    *sync.RWMutex
	}

	IPlayerMgr interface {
		actor.IActor
		ReleaseSocketMap(uint32, bool)
		AddPlayerMap(uint32, rpc.MailBox) int
		GetSocket(int64) uint32
		GetPlayerId(uint32) int64
		GetPlayer(uint32) *Player
	}

	Player struct {
		PlayerID   int64
		LastTime   int64
		SocketId   uint32
		GClusterId uint32
		ZClusterId uint32
	}
)

var (
	g_pPlayer = &Player{}
)

func NewPlayer(socket uint32, playerId int64) *Player {
	player := Player{LastTime: time.Now().Unix(), SocketId: socket, GClusterId: 0, PlayerID: playerId, ZClusterId: 0}
	return &player
}

func (this *PlayerMgr) ReleaseSocketMap(socketId uint32, bClose bool) {
	this.m_Locker.RLock()
	playerId, _ := this.m_SocketMap[socketId]
	this.m_Locker.RUnlock()
	this.m_Locker.Lock()
	delete(this.m_PlayerMap, playerId)
	delete(this.m_SocketMap, socketId)
	this.m_Locker.Unlock()
	//if bClose{
	SERVER.GetServer().StopClient(socketId)
	//}
}

func (this *PlayerMgr) AddPlayerMap(socketId uint32, mailbox rpc.MailBox) int {
	playerId := mailbox.Id
	Id := this.GetSocket(playerId)
	this.ReleaseSocketMap(Id, Id != socketId)

	pPlayer := NewPlayer(socketId, playerId)
	pPlayer.GClusterId = mailbox.ClusterId
	//pPlayer.ZClusterId = clusterInfo.ZClusterId
	this.m_Locker.Lock()
	this.m_PlayerMap[playerId] = pPlayer
	this.m_SocketMap[socketId] = playerId
	this.m_Locker.Unlock()
	SERVER.GetCluster().SendMsg(rpc.RpcHead{ClusterId: pPlayer.GClusterId, Id: playerId}, "game<-PlayerMgr.Player_Login", SERVER.GetCluster().Id(), mailbox)
	return base.NONE_ERROR
}

func (this *PlayerMgr) GetSocket(playerId int64) uint32 {
	socketId := uint32(0)
	this.m_Locker.RLock()
	pPlayer, exist := this.m_PlayerMap[playerId]
	this.m_Locker.RUnlock()
	if exist {
		socketId = pPlayer.SocketId
	}
	return socketId
}

func (this *PlayerMgr) GetPlayerId(socketId uint32) int64 {
	playerId := int64(0)
	this.m_Locker.RLock()
	id, exist := this.m_SocketMap[socketId]
	this.m_Locker.RUnlock()
	if exist {
		playerId = id
	}
	return playerId
}

func (this *PlayerMgr) GetPlayer(socketId uint32) *Player {
	playerId := this.GetPlayerId(socketId)
	this.m_Locker.RLock()
	pPlayer, exist := this.m_PlayerMap[playerId]
	this.m_Locker.RUnlock()
	if exist {
		return pPlayer
	}
	return nil
}

func (this *PlayerMgr) Init() {
	this.Actor.Init()
	this.m_SocketMap = make(map[uint32]int64)
	this.m_PlayerMap = make(map[int64]*Player)
	this.m_Locker = &sync.RWMutex{}
	actor.MGR.RegisterActor(this)
	this.Actor.Start()
}

func (this *PlayerMgr) ADD_ACCOUNT(ctx context.Context, socketId uint32, mailbox rpc.MailBox) {
	base.LOG.Printf("login incoming  Socket:%d PlayerId:%d GClusterId:%d ", socketId, mailbox.Id, mailbox.ClusterId)
	this.AddPlayerMap(socketId, mailbox)
}

func (this *PlayerMgr) DEL_ACCOUNT(ctx context.Context, socketid uint32) {
	playerId := this.GetPlayerId(socketid)
	this.ReleaseSocketMap(socketid, true)
	SERVER.GetCluster().SendMsg(rpc.RpcHead{SendType: rpc.SEND_BOARD_CAST}, "game<-PlayerMgr.G_ClientLost", playerId)
}

//重连世界服务器，账号重新登录
func (this *PlayerMgr) World_Relogin(ctx context.Context) {
	playerMap := make(map[int64]uint32)
	this.m_Locker.RLock()
	for i, v := range this.m_PlayerMap {
		playerMap[i] = v.GClusterId
	}
	this.m_Locker.RUnlock()

	if len(playerMap) != 0 {
		for i, v := range playerMap {
			SERVER.GetCluster().SendMsg(rpc.RpcHead{Id: i, ClusterId: v}, "G_W_Relogin")
		}
	}
}
