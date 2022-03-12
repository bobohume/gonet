package db

import (
	"context"
	"gonet/actor"
	"gonet/rpc"
	"gonet/server/model"
)

type (
	PlayerMgr struct {
		actor.Actor

		m_PlayerMap map[int64]*Player
	}

	IPlayerMgr interface {
		actor.IActor
	}

	Player struct {
		model.PlayerData
		Raft rpc.PlayerClusterInfo
	}
)

var (
	PLAYERMGR PlayerMgr
)

func (this *PlayerMgr) Init() {
	this.Actor.Init()
	this.m_PlayerMap = make(map[int64]*Player)
	this.RegisterTimer(60 * 1000*1000*1000, this.SavePlayerToDB) //定时器
	actor.MGR.RegisterActor(this)
	this.Actor.Start()
}

func (this *PlayerMgr) SavePlayerToDB() {
	// 存储玩家数据
	for _, v := range this.m_PlayerMap {
		v.SavePlayerDB()
	}
}

//玩家加载数据
func (this *PlayerMgr) Load_Player_DB(ctx context.Context, playerId int64, clusterInfo rpc.PlayerClusterInfo) {
	pPlayer, bEx := this.m_PlayerMap[playerId]
	if bEx{
		pPlayer.Raft = clusterInfo
		SERVER.GetCluster().SendMsg(rpc.RpcHead{DestServerType:rpc.SERVICE_GAME, ClusterId:this.GetRpcHead(ctx).SrcClusterId, Id:playerId}, "Load_Player_DB_Finish", pPlayer.PlayerData)
	}else{
		pPlayer = &Player{}
		pPlayer.Raft = clusterInfo
		err := pPlayer.LoadPlayerDB(playerId)
		if err == nil{
			this.m_PlayerMap[playerId] = pPlayer
			SERVER.GetCluster().SendMsg(rpc.RpcHead{DestServerType:rpc.SERVICE_GAME, ClusterId:this.GetRpcHead(ctx).SrcClusterId, Id:playerId}, "Load_Player_DB_Finish", pPlayer.PlayerData)
		}else{
			SERVER.GetLog().Printf("Player Load_Player_DB [%d] err[%s]", playerId, err.Error())
		}
	}
}

//lease过期
func (this *PlayerMgr) Player_Lease_Expire(ctx context.Context) {
	playerId := this.GetRpcHead(ctx).Id
	pPlayer, bEx := this.m_PlayerMap[playerId]
	if bEx{
		pPlayer.SavePlayerDB()
	}
	delete(this.m_PlayerMap, playerId)
	SERVER.GetLog().Printf("[%d] 过期删除玩家", playerId)
}

