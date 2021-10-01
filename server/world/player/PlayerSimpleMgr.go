package player

import (
	"database/sql"
	"gonet/actor"
	"gonet/base"
	"gonet/common"
	"gonet/db"
	"gonet/server/model"
	"gonet/server/world"
	"sync"
	"time"
)

type(
	PlayerSimpleMgr struct{
		actor.Actor
		m_SimplePlayerMap map[int64] *model.SimplePlayerData
		m_SimplePlayerNameMap map[string] *model.SimplePlayerData
		m_Locker *sync.RWMutex
		m_db *sql.DB
		m_Log *base.CLog
	}

	IPlayerSimpleMgr interface {
		actor.IActor

		LoadSimplePlayerDatas()
		GetPlayerDataByName(string) *model.SimplePlayerData
		GetPlayerDataById(int64) *model.SimplePlayerData
		GetPlayerName(int64) string
	}
)

var(
	SIMPLEMGR PlayerSimpleMgr
)

func loadSimple(row db.IRow, s *model.SimplePlayerData){
	s.AccountId = row.Int64("account_id")
	s.PlayerId = row.Int64("player_id")
	s.PlayerName = row.String("player_name")
	s.Level = row.Int("level")
	s.Sex = row.Int("sex")
	s.Gold = row.Int("gold")
	s.DrawGold = row.Int("draw_gold")
	s.Vip = row.Int("vip")
	s.LastLoginTime = row.Time("last_login_time")
	s.LastLogoutTime = row.Time("last_logout_time")
}

func (this *PlayerSimpleMgr) Init() {
	this.Actor.Init()
	this.m_Locker = &sync.RWMutex{}
	this.m_db = world.SERVER.GetDB()
	this.m_Log = world.SERVER.GetLog()
	this.m_SimplePlayerMap = make(map[int64] *model.SimplePlayerData)
	this.m_SimplePlayerNameMap = make(map[string] *model.SimplePlayerData)
	this.Actor.Start()
}

func (this *PlayerSimpleMgr) LoadSimplePlayerDatas() {
	startTime := time.Now().Unix()
	var simpledata model.SimplePlayerData
	rows, err := this.m_db.Query(db.LoadSql(simpledata, db.WithOutWhere()))
	if err != nil{
		common.DBERROR("LoadSimplePlayerDatas", err)
	}
	rs := db.Query(rows, err)
	for rs.Next(){
		pData := &model.SimplePlayerData{}
		loadSimple(rs.Row(), pData)
		//rs.Row().Obj(pData)
		this.m_Locker.Lock()
		this.m_SimplePlayerMap[pData.PlayerId] = pData
		this.m_SimplePlayerNameMap[pData.PlayerName] = pData
		this.m_Locker.Unlock()
	}

	endTime := time.Now().Unix()
	this.m_Log.Printf("结束读取玩家的简单信息[%d],timecost[%d]", startTime, endTime-startTime)
}

func (this *PlayerSimpleMgr) GetPlayerDataByName(name string) *model.SimplePlayerData{
	this.m_Locker.RLock()
	pData, exist := this.m_SimplePlayerNameMap[name]
	this.m_Locker.RUnlock()
	if exist{
		return pData
	}

	pData = LoadSimplePlayerDataByName(name)
	if pData != nil{
		this.m_Locker.Lock()
		this.m_SimplePlayerMap[pData.PlayerId] = pData
		this.m_SimplePlayerNameMap[name] = pData
		this.m_Locker.Unlock()
	}

	return pData
}

func (this *PlayerSimpleMgr) GetPlayerDataById(playerId int64) *model.SimplePlayerData{
	this.m_Locker.RLock()
	pData, exist := this.m_SimplePlayerMap[playerId]
	this.m_Locker.RUnlock()
	if exist{
		return pData
	}

	pData = LoadSimplePlayerData(playerId)
	if pData != nil{
		this.m_Locker.Lock()
		this.m_SimplePlayerMap[pData.PlayerId] = pData
		this.m_SimplePlayerNameMap[pData.PlayerName] = pData
		this.m_Locker.Unlock()
	}

	return pData
}

func (this *PlayerSimpleMgr) GetPlayerName(playerId int64) string{
	pData := this.GetPlayerDataById(playerId)
	if pData != nil {
		return pData.PlayerName
	}

	return  ""
}

func LoadSimplePlayerData(playerId int64) *model.SimplePlayerData{
	pData := &model.SimplePlayerData{PlayerId:playerId}
	rows, err := world.SERVER.GetDB().Query(db.LoadSql(pData))
	rs := db.Query(rows, err)
	if err == nil && rs.Next(){
		loadSimple(rs.Row(), pData)
		return pData
	}else if err != nil{
		common.DBERROR("LoadSimplePlayerData",err)
	}
	return nil
}

func LoadSimplePlayerDataByName(name string) *model.SimplePlayerData{
	pData := new(model.SimplePlayerData)
	rows, err := world.SERVER.GetDB().Query(db.LoadSql(pData, db.WithWhere(model.SimplePlayerData{PlayerName:name})))
	rs := db.Query(rows, err)
	if rs.Next(){
		loadSimple(rs.Row(), pData)
		return pData
	}
	return nil
}

func LoadSimplePlayerDatas(accountId int64) []*model.SimplePlayerData{
	pList := make([]*model.SimplePlayerData, 0)
	nPlayerNum := 0
	pData := new(model.SimplePlayerData)
	rows, err := world.SERVER.GetDB().Query(db.LoadSql(pData,  db.WithWhere(model.SimplePlayerData{AccountId:accountId})))
	rs := db.Query(rows, err)
	for rs.Next(){
		loadSimple(rs.Row(), pData)
		pList = append(pList, pData)
		nPlayerNum++
	}
	return pList
}