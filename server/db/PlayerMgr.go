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

func (this *PlayerSaveMgr) Init() {
	this.InitPool(this, reflect.TypeOf(PlayerMgr{}), MAX_PLAYER_MGR_COUNT)
	this.Stub.InitStub(rpc.STUB_PlayerMgr)
}

func (this *PlayerMgr) Init() {
	this.Actor.Init()
	//actor.MGR.RegisterActor(this)
	this.Actor.Start()
}

//玩家加载数据
func (this *PlayerMgr) Load_Player_DB(ctx context.Context, playerId int64, mailbox rpc.MailBox) {
	pPlayer := &Player{}
	err := pPlayer.LoadPlayerDB(playerId)
	if err == nil {
		cluster.MGR.SendMsg(rpc.RpcHead{ClusterId: this.GetRpcHead(ctx).SrcClusterId, Id: playerId}, "game<-Player.Load_Player_DB_Finish", pPlayer.PlayerData)
	} else {
		base.LOG.Printf("Player Load_Player_DB [%d] err[%s]", playerId, err.Error())
	}
}

func (this *PlayerMgr) Stub_On_Register(ctx context.Context, Id int64) {
	//这里可以是加载db数据
	base.LOG.Println("Stub db register sucess")
}

func (this *PlayerMgr) Stub_On_UnRegister(ctx context.Context, Id int64) {
	//lease一致性这里要清理缓存数据了
	base.LOG.Println("Stub db unregister sucess")
}
