package player

import (
	"database/sql"
	"fmt"
	"github.com/golang/protobuf/proto"
	"gonet/actor"
	"gonet/base"
	"gonet/db"
	"gonet/server/common"
	"gonet/server/world"
)
//********************************************************
// 玩家管理
//********************************************************
var(
	PLAYERMGR PlayerMgr
	PLAYER Player
)
type(
	PlayerMgr struct{
		actor.ActorPool

		m_db *sql.DB
		m_Log *base.CLog
		m_PingTimer common.ISimpleTimer
	}

	IPlayerMgr interface {
		actor.IActorPool

		GetPlayer(accountId int64) actor.IActor
		AddPlayer(accountId int64) actor.IActor
		RemovePlayer(accountId int64)
		Update()
	}
)

func (this* PlayerMgr) Init(num int){
	this.ActorPool.Init(num)
	this.m_db = world.SERVER.GetDB()
	this.m_Log = world.SERVER.GetLog()
	this.m_PingTimer = common.NewSimpleTimer(120)
	this.m_PingTimer.Start()
	actor.MGR.AddActor(this)
	this.RegisterTimer(1000 * 1000 * 1000, this.Update)//定时器
	//玩家登录
	this.RegisterCall("G_W_CLoginRequest", func(accountId int64) {
		pPlayer := this.GetPlayer(accountId)
		if pPlayer != nil{
			pPlayer.SendMsg("Logout", accountId)
			this.RemovePlayer(accountId)
		}

		pPlayer = this.AddPlayer(accountId)
		pPlayer.SendMsg("Login", this.GetSocketId())
	})

	//玩家断开链接
	this.RegisterCall("G_ClientLost", func(accountId int64) {
		pPlayer := this.GetPlayer(accountId)
		if pPlayer != nil{
			pPlayer.SendMsg("Logout", accountId)
		}

		this.RemovePlayer(accountId)
	})

	//account创建玩家反馈， 考虑到在创建角色的时候退出的情况
	this.RegisterCall("A_W_CreatePlayer", func(accountId int64, playerId int64, playername string, sex int32, socketId int) {
		rows, err := this.m_db.Query(fmt.Sprintf("call `sp_createplayer`(%d,'%s',%d, %d)", accountId, playername, sex, playerId))
		if err == nil && rows != nil{
			if rows.NextResultSet() && rows.NextResultSet(){
				rs := db.Query(rows, err)
				if rs.Next(){
					err := rs.Row().Int("@err")
					playerId := rs.Row().Int64("@playerId")
					//register
					if (err == 0) {
						this.m_Log.Printf("账号[%d]创建玩家[%d]", accountId, playerId)
					} else {
						this.m_Log.Printf("账号[%d]创建玩家失败", accountId)
						world.SERVER.GetAccountCluster().BalacaceMsg("W_A_DeletePlayer", accountId, playerId)
					}

					//通知玩家`
					pPlayer := this.GetPlayer(accountId)
					if pPlayer != nil {
						pPlayer.SendMsg("CreatePlayer", playerId, socketId, err)
					}
				}
			}
		}else{
			this.m_Log.Printf("账号[%d]创建玩家失败", accountId)
			world.SERVER.GetAccountCluster().BalacaceMsg("W_A_DeletePlayer", accountId, playerId)
		}
	})

	//this.RegisterTimer(1000 * 1000 * 1000, this.Update)//定时器
	PLAYER.Init(1)
	this.Actor.Start()
}

func (this *PlayerMgr) GetPlayer(accountId int64) actor.IActor{
	return this.GetActor(accountId)
}

func (this *PlayerMgr) AddPlayer(accountId int64) actor.IActor{
	LoadPlayerDB := func(accountId int64) ([]int64, int){
		PlayerList := make([]int64, 0)
		PlayerNum := 0
		rows, err := this.m_db.Query(fmt.Sprintf("select player_id from tbl_player where account_id=%d", accountId))
		rs := db.Query(rows, err)
		for rs.Next(){
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
	pPlayer.PlayerIdList = PlayerList
	pPlayer.PlayerNum = PlayerNum
	this.AddActor(accountId, pPlayer)
	pPlayer.Init(MAX_PLAYER_CHAN)
	return pPlayer
}

func (this *PlayerMgr) RemovePlayer(accountId int64){
	this.m_Log.Printf("移除帐号数据[%d]", accountId)
	this.DelActor(accountId)
}

func (this* PlayerMgr) Update(){
}

//--------------发送给客户端----------------------//
func SendToClient(AccountId int64, packet proto.Message){
	pPlayer := PLAYERMGR.GetPlayer(AccountId)
	if pPlayer != nil{
		 world.SendToClient(pPlayer.GetSocketId(), packet)
	}
}