package player

import (
	"database/sql"
	"base"
	"server/world"
)

const(
	EVENT_WAIT_KICKPLAYER	=	1<<0
	EVENT_CLIENTLOST			=	1<<1
	EVENT_LOGOUTED          	= 	1<<2
	MAX_PALYER_COUNT			= 1
)

const(
	STATUS_IDEL 				= 	iota
	STATUS_LOGIN				=	iota
	STATUS_IN_SELECT			=	iota
	STATUS_IN_GAME				=	iota
	STATUS_LOGOUT				=	iota
	STATUS_OFFLINE				=	iota
	STATUS_COUNT				=	iota
)

const(
	MAX_PLAYER_CHAN = 32
)

type (
	PlayerData struct{
		AccountId int64
		PlayerId int64
		SocketId int
		AccountName string
		PlayerSimpleData *SimplePlayerData

		PlayerNum int
		PlayerIdList []int64
		PlayerSimpleDataList []*SimplePlayerData
		m_db *sql.DB
		m_Log *base.CLog
	}

	IPlayerData interface {
		Init()
		GetPlayerId() int64
		GetPlayerCount() int
		SetPlayerId(int64) bool
		GetPlayerName() string
	}
)

func (this *PlayerData) Init(){
	this.m_db = world.SERVER.GetDB()
	this.m_Log = world.SERVER.GetLog()
	//this.PlayerIdList = make([]int, 0)
	//this.PlayerSimpleDataList = make([]*SimplePlayerData, 0)
}

func (this *PlayerData) GetPlayerId()int64{
	return this.PlayerId
}

func (this *PlayerData) GetPlayerName() string{
	if this.PlayerSimpleData != nil{
		return this.PlayerSimpleData.PlayerName
	}
	return ""
}

func (this *PlayerData) GetPlayerCount()int{
	count := 0
	for i := 0; i < len(this.PlayerIdList); i++ {
		if this.PlayerIdList[i] != 0 {
			count++
		}
	}
	return count
}

func (this *PlayerData) SetPlayerId(PlayerId int64) bool{
	for i := 0; i < len(this.PlayerIdList); i++ {
		if this.PlayerIdList[i] == PlayerId {
			this.PlayerId = PlayerId
			this.PlayerSimpleData = this.PlayerSimpleDataList[i]
			return  true
		}
	}
	return  false
}

