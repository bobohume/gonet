package player

import (
	"context"
	"database/sql"
	"gonet/actor"
	"gonet/base"
	"gonet/common"
	"gonet/rpc"
	"gonet/server/game"
)

//********************************************************
// 玩家管理
//********************************************************
var (
	MGR PlayerMgr
	PLAYER Player
)

type (
	PlayerMgr struct {
		actor.Actor

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
	actor.MGR.RegisterActor(this)
	actor.MGR.RegisterActor(&PLAYER, actor.WithType(actor.ACTOR_TYPE_PLAYER))
	this.Actor.Start()
}

//玩家登录
func (this *PlayerMgr) LoginPlayerRequset(ctx context.Context, playerId int64, gateClusterId uint32, socketId uint32) {
	pPlayerCluster := game.SERVER.GetPlayerRaft().GetPlayer(playerId)
	if pPlayerCluster == nil {
		info := &rpc.PlayerClusterInfo{}
		info.Id = playerId
		info.GClusterId = game.SERVER.GetCluster().Id()
		info.ZClusterId = game.SERVER.GetCluster().RandomCluster(rpc.RpcHead{Id: playerId, DestServerType: rpc.SERVICE_ZONE}).ClusterId
		info.DClusterId = game.SERVER.GetCluster().RandomCluster(rpc.RpcHead{Id: playerId, DestServerType: rpc.SERVICE_DB}).ClusterId
		if game.SERVER.GetPlayerRaft().Publish(info) {
			pPlayerCluster = info
		}else{
			return
		}
	}

	game.SERVER.GetCluster().SendMsg(rpc.RpcHead{ClusterId: gateClusterId, DestServerType: rpc.SERVICE_GATE}, "A_G_Account_Login", socketId, *pPlayerCluster)
}

//玩家登录
func (this *PlayerMgr) Player_Login(ctx context.Context, gateClusterId uint32, clusterInfo rpc.PlayerClusterInfo) {
	playerId := clusterInfo.Id
	if actor.MGR.GetPlayer(playerId) != nil{
		actor.MGR.SendMsg(rpc.RpcHead{Id:playerId}, "ReLogin", gateClusterId, clusterInfo)
		return
	}

	pPlayer := &Player{}
	pPlayer.PlayerId = playerId
	pPlayer.SetId(playerId)
	pPlayer.Init()
	actor.MGR.AddPlayer(pPlayer)
	actor.MGR.SendMsg(rpc.RpcHead{Id:playerId}, "Login", gateClusterId, clusterInfo)
}

//玩家断开链接
func (this *PlayerMgr) G_ClientLost(ctx context.Context, playerId int64) {
	actor.MGR.SendMsg(rpc.RpcHead{Id:playerId}, "Logout", playerId)
}
