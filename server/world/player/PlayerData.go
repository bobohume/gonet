package player

import (
	"database/sql"
	"gonet/base"
	"gonet/db"
	"gonet/rpc"
	"gonet/server/model"
	"gonet/server/world"
)

const (
	EVENT_WAIT_KICKPLAYER = 1 << 0
	EVENT_CLIENTLOST      = 1 << 1
	EVENT_LOGOUTED        = 1 << 2
	MAX_PALYER_COUNT      = 1
)

const (
	STATUS_IDEL      = iota
	STATUS_LOGIN     = iota
	STATUS_IN_SELECT = iota
	STATUS_IN_GAME   = iota
	STATUS_LOGOUT    = iota
	STATUS_OFFLINE   = iota
	STATUS_COUNT     = iota
)

type (
	PlayerData struct {
		model.SimplePlayerData

		AccountId            int64
		PlayerId             int64
		GateClusterId        uint32
		m_PlayerRaft         rpc.PlayerClusterInfo
		AccountName          string
		PlayerNum            int
		PlayerIdList         []int64
		PlayerSimpleDataList []*model.SimplePlayerData
		m_PlayerKVMap        map[int]*model.PlayerKvData
		m_db                 *sql.DB
		m_Log                *base.CLog
	}

	IPlayerData interface {
		Init()

		SetGateClusterId(uint32)
		GetGateClusterId() uint32
		GetZoneClusterId() uint32
		GetAccountId() int64
		GetPlayerCount() int
		GetLeaseId() int64
		SetPlayerId(int64) bool
		GetPlayerId() int64
		GetPlayerName() string

		LoadPlayerData() //加载其他数据
		//----KV---//
		LoadKV()                    //加载kv
		SetKV(key int, value int64) //设置kv
		DelKV(key int)              //删除key
		GetKV(key int) int64        //获取key
		//----KV---//
	}
)

func (this *PlayerData) Init() {
	this.m_db = world.SERVER.GetDB()
	this.m_Log = world.SERVER.GetLog()
	this.m_PlayerKVMap = map[int]*model.PlayerKvData{}
	//this.PlayerIdList = make([]int, 0)
	//this.PlayerSimpleDataList = make([]*SimplePlayerData, 0)
}

func (this *PlayerData) GetZoneClusterId() uint32 {
	return this.m_PlayerRaft.ZClusterId
}

func (this *PlayerData) SetGateClusterId(clusterId uint32) {
	this.GateClusterId = clusterId
}

func (this *PlayerData) GetGateClusterId() uint32 {
	return this.GateClusterId
}

func (this *PlayerData) GetLeaseId() int64 {
	return this.m_PlayerRaft.LeaseId
}

func (this *PlayerData) GetAccountId() int64 {
	return this.AccountId
}

func (this *PlayerData) SetPlayerId(PlayerId int64) bool {
	for i := 0; i < len(this.PlayerIdList); i++ {
		if this.PlayerIdList[i] == PlayerId {
			this.PlayerId = PlayerId
			this.SimplePlayerData = *this.PlayerSimpleDataList[i]
			return true
		}
	}
	return false
}

func (this *PlayerData) GetPlayerId() int64 {
	return this.PlayerId
}

func (this *PlayerData) GetPlayerName() string {
	return this.PlayerName
}

func (this *PlayerData) GetPlayerCount() int {
	count := 0
	for i := 0; i < len(this.PlayerIdList); i++ {
		if this.PlayerIdList[i] != 0 {
			count++
		}
	}
	return count
}

func (this *PlayerData) LoadPlayerData() {
	//加载kv数据
	this.LoadKV()
}

//-------------kv--------------//
func (this *PlayerData) LoadKV() {
	pData := &model.PlayerKvData{}
	rows, err := this.m_db.Query(db.LoadSql(pData, db.WithWhere(model.PlayerKvData{PlayerId:this.GetPlayerId()})))
	rs := db.Query(rows, err)
	for rs.Next() {
		pData := &model.PlayerKvData{}
		db.LoadObjSql(pData, rs.Row())
		this.m_PlayerKVMap[pData.Key] = pData
	}
}

func (this *PlayerData) SetKV(key int, value int64) {
	pDdata, bEx := this.m_PlayerKVMap[key]
	if bEx && pDdata != nil {
		pDdata.Value = value
		this.m_db.Exec(db.UpdateSql(pDdata))
	} else {
		pDdata = &model.PlayerKvData{PlayerId: this.GetPlayerId(), Key: key, Value: value}
		this.m_PlayerKVMap[key] = pDdata
		this.m_db.Exec(db.InsertSql(pDdata))
	}
}

func (this *PlayerData) DelKV(key int) {
	pDdata, bEx := this.m_PlayerKVMap[key]
	if bEx && pDdata != nil {
		this.m_db.Exec(db.DeleteSql(pDdata))
		delete(this.m_PlayerKVMap, key)
	}
}

func (this *PlayerData) GetKV(key int) int64 {
	pDdata, bEx := this.m_PlayerKVMap[key]
	if bEx && pDdata != nil {
		return pDdata.Value
	}
	return 0
}
