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
		socketMap map[uint32]int64
		playerMap map[int64]*Player
		locker    *sync.RWMutex
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

func (p *PlayerMgr) ReleaseSocketMap(socketId uint32, bClose bool) {
	p.locker.RLock()
	playerId, _ := p.socketMap[socketId]
	p.locker.RUnlock()
	p.locker.Lock()
	delete(p.playerMap, playerId)
	delete(p.socketMap, socketId)
	p.locker.Unlock()
	//if bClose{
	SERVER.GetServer().StopClient(socketId)
	//}
}

func (p *PlayerMgr) AddPlayerMap(socketId uint32, mailbox rpc.MailBox) int {
	playerId := mailbox.Id
	Id := p.GetSocket(playerId)
	p.ReleaseSocketMap(Id, Id != socketId)

	player := NewPlayer(socketId, playerId)
	player.GClusterId = mailbox.ClusterId
	//player.ZClusterId = clusterInfo.ZClusterId
	p.locker.Lock()
	p.playerMap[playerId] = player
	p.socketMap[socketId] = playerId
	p.locker.Unlock()
	SERVER.GetCluster().SendMsg(rpc.RpcHead{ClusterId: player.GClusterId, Id: playerId}, "game<-PlayerMgr.Player_Login", SERVER.GetCluster().Id(), mailbox)
	return base.NONE_ERROR
}

func (p *PlayerMgr) GetSocket(playerId int64) uint32 {
	socketId := uint32(0)
	p.locker.RLock()
	player, exist := p.playerMap[playerId]
	p.locker.RUnlock()
	if exist {
		socketId = player.SocketId
	}
	return socketId
}

func (p *PlayerMgr) GetPlayerId(socketId uint32) int64 {
	playerId := int64(0)
	p.locker.RLock()
	id, exist := p.socketMap[socketId]
	p.locker.RUnlock()
	if exist {
		playerId = id
	}
	return playerId
}

func (p *PlayerMgr) GetPlayer(socketId uint32) *Player {
	playerId := p.GetPlayerId(socketId)
	p.locker.RLock()
	player, exist := p.playerMap[playerId]
	p.locker.RUnlock()
	if exist {
		return player
	}
	return nil
}

func (p *PlayerMgr) Init() {
	p.Actor.Init()
	p.socketMap = make(map[uint32]int64)
	p.playerMap = make(map[int64]*Player)
	p.locker = &sync.RWMutex{}
	actor.MGR.RegisterActor(p)
	p.Actor.Start()
}

func (p *PlayerMgr) ADD_ACCOUNT(ctx context.Context, socketId uint32, mailbox rpc.MailBox) {
	base.LOG.Printf("login incoming  Socket:%d PlayerId:%d GClusterId:%d ", socketId, mailbox.Id, mailbox.ClusterId)
	p.AddPlayerMap(socketId, mailbox)
}

func (p *PlayerMgr) DEL_ACCOUNT(ctx context.Context, socketid uint32) {
	base.LOG.Printf("DELACCOUNT Socket:%d ", socketid)
	playerId := p.GetPlayerId(socketid)
	p.ReleaseSocketMap(socketid, true)
	SERVER.GetCluster().SendMsg(rpc.RpcHead{SendType: rpc.SEND_BOARD_CAST}, "game<-PlayerMgr.G_ClientLost", playerId)
}

// 重连世界服务器，账号重新登录
func (p *PlayerMgr) World_Relogin(ctx context.Context) {
	playerMap := make(map[int64]uint32)
	p.locker.RLock()
	for i, v := range p.playerMap {
		playerMap[i] = v.GClusterId
	}
	p.locker.RUnlock()

	if len(playerMap) != 0 {
		for i, v := range playerMap {
			SERVER.GetCluster().SendMsg(rpc.RpcHead{Id: i, ClusterId: v}, "G_W_Relogin")
		}
	}
}
