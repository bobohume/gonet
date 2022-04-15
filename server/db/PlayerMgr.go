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
	this.Stub.InitStub("", rpc.STUB_PLAYERDB)
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

//lease过期
func (this *PlayerMgr) Player_On_UnRegister(ctx context.Context) {
	playerId := this.GetRpcHead(ctx).Id
	base.LOG.Printf("[%d] 过期删除玩家", playerId)
}
