package worlddb

import (
	"database/sql"
	"gonet/actor"
	"gonet/db"
	"gonet/server/world"
	"time"
)

type (
	Player struct {
		UpdateTime int64
		UpdateTTLTime int64
		PlayerBlob []byte
	}

	PlayerMgr struct {
		actor.Actor

		m_PlayerMap map[int64] *Player
		m_db *sql.DB
	}

	IPlayerMgr interface {
		actor.IActor
	}
)

var (
	PLAYERMGR PlayerMgr
)

func (this* PlayerMgr) Init(num int){
	this.m_db = SERVER.GetDB()
	this.Actor.Init(1000)
	this.m_PlayerMap = make(map[int64] *Player)
	this.RegisterTimer(1000 * 1000 * 1000, this.Update)//定时器
	//load blob
	this.RegisterCall("Load_Player", func(playerId int64, accountId int64) {
		pPlayer, bEx := this.m_PlayerMap[playerId]
		if !bEx {
			//加载人物数据
			rows, err := world.SERVER.GetDB().Query("select `player_blob` from tbl_player where player_id = ?", playerId)
			//加载错误
			if err != nil{

			}
			rs := db.Query(rows, err)
			if rs.Next() {
				pPlayer = &Player{UpdateTime: time.Now().Unix(), PlayerBlob: rs.Row().Byte("player_blob"), UpdateTTLTime:time.Now().Unix()}
				this.m_PlayerMap[playerId] = pPlayer
			}
		}

		//发送人物数据
		SERVER.GetServer().SendMsgByID(this.GetSocketId(), "Load_Player_Finish", accountId, pPlayer.PlayerBlob)
	})

	//save blob
	this.RegisterCall("Save_Player", func(playerId int64, playerBlob []byte, accountId int64) {
		pPlayer, bEx := this.m_PlayerMap[playerId]
		if bEx {
			pPlayer.PlayerBlob = playerBlob
			pPlayer.UpdateTTLTime = time.Now().Unix()
			//设置redis ttl
		}
	})

	this.Actor.Start()
}


func (this *PlayerMgr) Update(){
	nTime := time.Now().Unix()
	DeletePlayers := []int64{}
	for i, v := range this.m_PlayerMap{
		//更新到数据库
		if nTime > v.UpdateTime + int64(3 * time.Minute){
			world.SERVER.GetDB().Exec("update tbl_player set `player_blob` = ? where player_id = ?", v.PlayerBlob, i)
			v.UpdateTime = nTime
		}

		//更新redis ttl时间
		if nTime > v.UpdateTTLTime + int64(3 * time.Minute){
			DeletePlayers = append(DeletePlayers, i)
		}
	}
	//删除过期缓存
	for _, v := range DeletePlayers{
		delete(this.m_PlayerMap, v)
	}
}

