package player

import (
	"context"
	"database/sql"
	"gonet/actor"
	"gonet/base"
	"gonet/common"
	"gonet/common/cluster"
	"gonet/rpc"
	"reflect"
)

//********************************************************
// 玩家管理
//********************************************************
var (
	MGR PlayerMgr
)

type (
	PlayerMgr struct {
		actor.Actor
		actor.VirtualActor
		m_db        *sql.DB
		m_Log       *base.CLog
		m_PingTimer common.ISimpleTimer
	}

	IPlayerMgr interface {
		Update()
	}
)

func (this *PlayerMgr) Init() {
	this.Actor.Init()
	this.m_PingTimer = common.NewSimpleTimer(120)
	this.m_PingTimer.Start()
	this.InitActor(this, reflect.TypeOf(Player{}))
	actor.MGR.RegisterActor(this)
	this.Actor.Start()
}

//玩家登录
func (this *PlayerMgr) LoginPlayerRequset(ctx context.Context, playerId int64, gateClusterId uint32, socketId uint32) {
	pMailBox := cluster.MGR.MailBox.Get(playerId)
	if pMailBox == nil {
		info := &rpc.MailBox{}
		info.Id = playerId
		info.ClusterId = cluster.MGR.Id()
		info.MailType = rpc.MAIL_Player
		if cluster.MGR.MailBox.Create(info) {
			pMailBox = info
		} else {
			return
		}
	}

	cluster.MGR.SendMsg(rpc.RpcHead{ClusterId: gateClusterId}, "gate<-EventProcess.G_Player_Login", socketId, *pMailBox)
}

//玩家登录
func (this *PlayerMgr) Player_Login(ctx context.Context, gateClusterId uint32, mailbox rpc.MailBox) {
	playerId := mailbox.Id
	if this.GetActor(playerId) != nil {
		actor.MGR.SendMsg(rpc.RpcHead{Id: playerId}, "Player.ReLogin", gateClusterId, mailbox)
		return
	}

	pPlayer := &Player{}
	pPlayer.PlayerId = playerId
	pPlayer.SetId(playerId)
	pPlayer.Init()
	this.AddActor(pPlayer)
	actor.MGR.SendMsg(rpc.RpcHead{Id: playerId}, "Player.Login", gateClusterId, mailbox)
}

//玩家断开链接
func (this *PlayerMgr) G_ClientLost(ctx context.Context, playerId int64) {
	actor.MGR.SendMsg(rpc.RpcHead{Id: playerId}, "Player.Logout", playerId)
}
