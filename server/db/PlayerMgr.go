package db

import (
	"context"
	"gonet/actor"
	"gonet/base"
	"gonet/base/cluster"
	"gonet/rpc"
	"gonet/server/model"
	"reflect"
	"time"
)

type (
	PlayerSaveMgr struct {
		actor.ActorPool
		cluster.Stub
	}

	PlayerMgr struct {
		actor.Actor
		PlayerMap map[int64]*Player
	}

	IPlayerMgr interface {
		actor.IActor
	}

	Player struct {
		model.PlayerData
	}
)

var (
	PLAYERSAVEMGR PlayerSaveMgr
)

const (
	MAX_PLAYER_MGR_COUNT = 10
)

func (p *PlayerSaveMgr) Init() {
	p.InitPool(p, reflect.TypeOf(PlayerMgr{}), MAX_PLAYER_MGR_COUNT)
	p.Stub.InitStub(rpc.STUB_PlayerMgr)
}

func (p *PlayerMgr) Init() {
	p.Actor.Init()
	//actor.MGR.RegisterActor(p)
	p.PlayerMap = make(map[int64]*Player)
	p.RegisterTimer(60*time.Second, p.SaveDB) //定时器
	p.Actor.Start()
}

func (p *PlayerMgr) SaveDB() {
	for _, player := range p.PlayerMap {
		player.SavePlayerDB()
	}
}

func (p *PlayerMgr) Load_Player_DB(ctx context.Context, playerId int64, mailbox rpc.MailBox) {
	player := p.GetPlayer(playerId)
	if player != nil {
		cluster.MGR.SendMsg(rpc.RpcHead{ClusterId: p.GetRpcHead(ctx).SrcClusterId, Id: playerId}, "game<-Player.Load_Player_DB_Finish", player.PlayerData)
	}
}

func (p *PlayerMgr) GetPlayer(Id int64) *Player {
	player, _ := p.PlayerMap[Id]
	if player == nil {
		player = &Player{}
		err := player.LoadPlayerDB(Id)
		if err == nil {
			p.PlayerMap[Id] = player
			return player
		} else {
			base.LOG.Printf("PlayerMgr GetData [%d] err[%s]", Id, err.Error())
		}
	}
	return player
}

func (p *PlayerMgr) Load(ctx context.Context, playerId int64, mailbox rpc.MailBox) {
	player := &Player{}
	err := player.LoadPlayerDB(playerId)
	if err == nil {
		cluster.MGR.SendMsg(rpc.RpcHead{ClusterId: p.GetRpcHead(ctx).SrcClusterId, Id: playerId}, "game<-Player.Load_Player_DB_Finish", player.PlayerData)
	} else {
		base.LOG.Printf("Player Load_Player_DB [%d] err[%s]", playerId, err.Error())
	}
}

func (p *PlayerMgr) OnStubRegister(ctx context.Context) {
	//这里可以是加载db数据
	base.LOG.Println("Stub db register sucess")
}

func (p *PlayerMgr) OnStubUnRegister(ctx context.Context) {
	//lease一致性这里要清理缓存数据了
	base.LOG.Println("Stub db unregister sucess")
	p.SaveDB()
	p.PlayerMap = make(map[int64]*Player)
}
