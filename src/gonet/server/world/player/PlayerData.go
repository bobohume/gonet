package player

import (
	"database/sql"
	"fmt"
	"gonet/base"
	"gonet/db"
	"gonet/server/world"
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
		m_PlayerKVMap map[int] *PlayerKvData
		m_db *sql.DB
		m_Log *base.CLog
	}

	IPlayerData interface {
		Init()
		GetGateSocketId() int
		GetAccountId() int64
		GetPlayerId() int64
		GetPlayerCount() int
		SetPlayerId(int64) bool
		GetPlayerName() string

		LoadPlayerData()//加载其他数据
		//----KV---//
		LoadKV()//加载kv
		SetKV(key int, value int64)//设置kv
		DelKV(key int)//删除key
		GetKV(key int) int64//获取key
		//----KV---//
	}
)

func (this *PlayerData) Init(){
	this.m_db = world.SERVER.GetDB()
	this.m_Log = world.SERVER.GetLog()
	this.m_PlayerKVMap = map[int]*PlayerKvData{}
	//this.PlayerIdList = make([]int, 0)
	//this.PlayerSimpleDataList = make([]*SimplePlayerData, 0)
}

func (this *PlayerData) GetGateSocketId() int{
	return this.SocketId
}

func (this *PlayerData) GetAccountId() int64{
	return this.AccountId
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

func (this *PlayerData) LoadPlayerData() {
	//加载kv数据
	this.LoadKV()
}
//-------------kv--------------//
func (this *PlayerData) LoadKV() {
	pData := &PlayerKvData{}
	rows, err := this.m_db.Query(db.LoadSql(pData, "tbl_player_kv", fmt.Sprintf("player_id = %d", this.GetPlayerId())))
	rs := db.Query(rows, err)
	for rs.Next(){
		pData := &PlayerKvData{}
		db.LoadObjSql(pData, rs.Row())
		this.m_PlayerKVMap[pData.Key] = pData
	}
}

func (this *PlayerData) SetKV(key int, value int64){
	pDdata, bEx := this.m_PlayerKVMap[key]
	if bEx && pDdata != nil{
		pDdata.Value = value
		this.m_db.Exec(db.UpdateSqlEx(pDdata, "tbl_player_kv", "value"))
	}else{
		pDdata = &PlayerKvData{PlayerId:this.GetPlayerId(), Key:key, Value:value}
		this.m_PlayerKVMap[key] = pDdata
		this.m_db.Exec(db.InsertSql(pDdata, "tbl_player_kv"))
	}
}

func (this *PlayerData) DelKV(key int){
	pDdata, bEx := this.m_PlayerKVMap[key]
	if bEx && pDdata != nil{
		this.m_db.Exec(db.DeleteSql(pDdata, "tbl_player_kv"))
	}
}

func (this *PlayerData) GetKV(key int) int64{
	pDdata, bEx := this.m_PlayerKVMap[key]
	if bEx && pDdata != nil{
		return pDdata.Value
	}
	return 0
}

