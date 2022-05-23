package db

import (
	"context"
	"gonet/actor"
	"gonet/base"
	"gonet/common/cluster"
	"gonet/rpc"
	"gonet/server/model"
	"reflect"
)

type (
	PlayerSaveMgr struct {
		actor.ActorPool
		cluster.Stub
	}

	PlayerMgr struct {
		actor.Actor
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
	p.Actor.Start()
}

//玩家加载数据
func (p *PlayerMgr) Load_Player_DB(ctx context.Context, playerId int64, mailbox rpc.MailBox) {
	pPlayer := &Player{}
	err := pPlayer.LoadPlayerDB(playerId)
	if err == nil {
		cluster.MGR.SendMsg(rpc.RpcHead{ClusterId: p.GetRpcHead(ctx).SrcClusterId, Id: playerId}, "game<-Player.Load_Player_DB_Finish", pPlayer.PlayerData)
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
}
