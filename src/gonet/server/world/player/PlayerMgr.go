package player

import (
	"database/sql"
	"fmt"
	"gonet/actor"
	"gonet/base"
	"gonet/db"
	"gonet/message"
	"gonet/rpc"
	"gonet/server/common"
	"gonet/server/world"
	"strings"
)
//********************************************************
// 玩家管理
//********************************************************
var(
	PLAYERMGR PlayerMgr
)

type(
	PlayerMgr struct{
		actor.ActorPool//玩家actor线城池

		m_db         *sql.DB
		m_Log        *base.CLog
		m_PingTimer  common.ISimpleTimer
	}

	IPlayerMgr interface {
		actor.IActor

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
		//查询playerid是否唯一
		error := 1
		rows, err := this.m_db.Query(fmt.Sprintf("select 1 from tbl_player where player_id = %d", playerId))
		if err == nil{
			rs := db.Query(rows, err)
			if !rs.Next(){
				//查找账号玩家数量
				rows, err = this.m_db.Query(fmt.Sprintf("select count(player_id) as player_count from tbl_player where account_id = %d", accountId))
				if err == nil {
					rs = db.Query(rows, err)
					if rs.Next(){
						player_count := rs.Row().Int("player_count")
						if player_count >= 1{//创建玩家上限
							this.m_Log.Printf("账号[%d]创建玩家数量上限", accountId)
						}else{//创建玩家
							_, err = this.m_db.Exec(fmt.Sprintf("insert into tbl_player (account_id, player_id, player_name, sex, level, gold, draw_gold)" +
								"values(%d, %d, '%s', %d, 1, 0,	0)", accountId, playerId, playername, sex))
							if err == nil{
								this.m_Log.Printf("账号[%d]创建玩家[%d]", accountId, playerId)
								error = 0
							}

							//通知玩家`
							pPlayer := this.GetPlayer(accountId)
							if pPlayer != nil {
								pPlayer.SendMsg("CreatePlayer", playerId, socketId, error)
							}
						}
					}
				}
			}
		}

		if error == 1 {//创建失败通知accout删除player
			this.m_Log.Printf("账号[%d]创建玩家[%d]失败", accountId, playerId)
			world.SERVER.GetAccountCluster().BalanceMsg("W_A_DeletePlayer", accountId, playerId)
		}
	})

	//this.RegisterTimer(1000 * 1000 * 1000, this.Update)//定时器
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

func (this *PlayerMgr) PacketFunc(id int, buff []byte) bool{
	defer func() {
		if err := recover(); err != nil {
			base.TraceCode(err)
		}
	}()

	var io actor.CallIO
	io.Buff = buff
	io.SocketId = id

	bitstream := base.NewBitStream(io.Buff, len(io.Buff))
	funcName := bitstream.ReadString()
	funcName = strings.ToLower(funcName)
	if this.FindCall(funcName) != nil{
		this.Send(io)
		return true
	}else{
		bitstream.ReadInt(base.Bit8)
		nType := bitstream.ReadInt(base.Bit8)
		if nType == rpc.RPC_INT64 || nType == rpc.RPC_INT64_PTR{
			nId := rpc.ReadInt64(bitstream)
			return this.SendById(nId, funcName, io)
		}else if nType == rpc.RPC_MESSAGE{
			packet, err := rpc.UnmarshalPB(bitstream)
			if err != nil{
				return false
			}
			packetHead := packet.(message.Packet).GetPacketHead()
			nId := packetHead.Id
			return this.SendById(nId, funcName, io)
		}
	}

	return false
}