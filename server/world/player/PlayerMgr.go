package player

import (
	"context"
	"database/sql"
	"fmt"
	"gonet/actor"
	"gonet/base"
	"gonet/common"
	"gonet/db"
	"gonet/rpc"
	"gonet/server/world"
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
		AddPlayer(accountId int64) actor.IActor
		RemovePlayer(accountId int64)
		Update()
	}
)

func (this *PlayerMgr) Init() {
	this.Actor.Init()
	this.m_db = world.SERVER.GetDB()
	this.m_Log = world.SERVER.GetLog()
	this.m_PingTimer = common.NewSimpleTimer(120)
	this.m_PingTimer.Start()
	actor.MGR.RegisterActor(this)
	actor.MGR.RegisterActor(&PLAYER, actor.WithType(actor.ACTOR_TYPE_PLAYER))
	this.Actor.Start()
}


func (this *PlayerMgr) AddPlayer(accountId int64) actor.IActor {
	LoadPlayerDB := func(accountId int64) ([]int64, int) {
		PlayerList := make([]int64, 0)
		PlayerNum := 0
		rows, err := this.m_db.Query(fmt.Sprintf("select player_id from tbl_player where account_id=%d", accountId))
		rs := db.Query(rows, err)
		for rs.Next() {
			PlayerId := rs.Row().Int64("player_id")
			PlayerList = append(PlayerList, PlayerId)
			PlayerNum++
		}
		return PlayerList, PlayerNum
	}

	fmt.Printf("玩家[%d]登录", accountId)
	PlayerList, PlayerNum := LoadPlayerDB(accountId)
	pPlayer := &Player{}
	pPlayer.AccountId = accountId
	pPlayer.SetId(accountId)
	pPlayer.PlayerIdList = PlayerList
	pPlayer.PlayerNum = PlayerNum
	pPlayer.Init()
	actor.MGR.AddPlayer(pPlayer)
	return pPlayer
}

//玩家登录
func (this *PlayerMgr) G_W_CLoginRequest(ctx context.Context, accountId int64, gateClusterId uint32, clusterInfo rpc.PlayerClusterInfo) {
	pPlayer := actor.MGR.GetPlayer(accountId)
	if pPlayer != nil {
		actor.MGR.SendMsg(rpc.RpcHead{Id:accountId}, "Logout", accountId)
	}

	this.AddPlayer(accountId)
	actor.MGR.SendMsg(rpc.RpcHead{Id:accountId}, "Login", gateClusterId, clusterInfo)
}

//玩家断开链接
func (this *PlayerMgr) G_ClientLost(ctx context.Context, accountId int64) {
	actor.MGR.SendMsg(rpc.RpcHead{Id:accountId}, "Logout", accountId)
	actor.MGR.DelPlayer(accountId)
}

//account创建玩家反馈， 考虑到在创建角色的时候退出的情况
func (this *PlayerMgr) A_W_CreatePlayer(ctx context.Context, accountId int64, playerId int64, playername string, sex int32, gClusterId uint32) {
	//查询playerid是否唯一
	error := 1
	rows, err := this.m_db.Query(fmt.Sprintf("select 1 from tbl_player where player_id = %d", playerId))
	if err == nil {
		rs := db.Query(rows, err)
		if !rs.Next() {
			//查找账号玩家数量
			rows, err = this.m_db.Query(fmt.Sprintf("select count(player_id) as player_count from tbl_player where account_id = %d", accountId))
			if err == nil {
				rs = db.Query(rows, err)
				if rs.Next() {
					player_count := rs.Row().Int("player_count")
					if player_count >= 1 { //创建玩家上限
						this.m_Log.Printf("账号[%d]创建玩家数量上限", accountId)
					} else { //创建玩家
						_, err = this.m_db.Exec(fmt.Sprintf("insert into tbl_player (account_id, player_id, player_name, sex, level, gold, draw_gold)"+
							"values(%d, %d, '%s', %d, 1, 0,	0)", accountId, playerId, playername, sex))
						if err == nil {
							this.m_Log.Printf("账号[%d]创建玩家[%d]", accountId, playerId)
							error = 0
						}

						//通知玩家`
						actor.MGR.SendMsg(rpc.RpcHead{Id:accountId}, "CreatePlayer", playerId, gClusterId, error)
					}
				}
			}
		}
	}

	if error == 1 { //创建失败通知accout删除player
		this.m_Log.Printf("账号[%d]创建玩家[%d]失败", accountId, playerId)
		world.SendToAccount("W_A_DeletePlayer", accountId, playerId)
	}
}
