package player

import (
	"actor"
	"time"
	"database/sql"
	"base"
	"fmt"
	"db"
	"sync"
	"server/common"
	"server/world"
)

type(
	PlayerSimpleMgr struct{
		actor.Actor
		m_SimplePlayerMap map[int] *SimplePlayerData
		m_SimplePlayerNameMap map[string] *SimplePlayerData
		m_Locker *sync.RWMutex
		m_db *sql.DB
		m_Log *base.CLog
	}

	IPlayerSimpleMgr interface {
		actor.IActor

		LoadSimplePlayerDatas()
		GetPlayerDataByName(string) *SimplePlayerData
		GetPlayerDataById(int) *SimplePlayerData
		GetPlayerName(int) string
	}
)

var(
	PLAYERSIMPLEMGR PlayerSimpleMgr
)

func loadSimple(row db.IRow, s *SimplePlayerData){
	s.AccountId = row.Int("account_id")
	s.PlayerId = row.Int("player_id")
	s.PlayerName = row.String("player_name")
	s.Level = row.Int("level")
	s.Sex = row.Int("sex")
	s.Gold = row.Int("gold")
	s.DrawGold = row.Int("draw_gold")
	s.Vip = row.Int("vip")
	s.LastLoginTime = row.Time("last_login_time")
	s.LastLogoutTime = row.Time("last_logout_time")
}

func (this *PlayerSimpleMgr) Init(num int) {
	this.Actor.Init(num)
	//注册结构体
	base.RegisterMessage(&SimplePlayerData{})

	this.m_Locker = &sync.RWMutex{}
	this.m_db = world.SERVER.GetDB()
	this.m_Log = world.SERVER.GetLog()
	this.m_SimplePlayerMap = make(map[int] *SimplePlayerData)
	this.m_SimplePlayerNameMap = make(map[string] *SimplePlayerData)

	this.Actor.Start()
}

func (this *PlayerSimpleMgr) LoadSimplePlayerDatas() {
	startTime := time.Now().Unix()
	var simpledata SimplePlayerData
	rows, err := this.m_db.Query(db.LoadSql(simpledata, "tbl_player", ""));
	if err != nil{
		common.DBERROR("LoadSimplePlayerDatas", err)
	}
	rs := db.Query(rows)
	for rs.Next(){
		pData := &SimplePlayerData{}
		loadSimple(rs.Row(), pData)
		this.m_Locker.Lock()
		this.m_SimplePlayerMap[pData.PlayerId] = pData
		this.m_SimplePlayerNameMap[pData.PlayerName] = pData
		this.m_Locker.Unlock()
	}

	endTime := time.Now().Unix()
	this.m_Log.Println("结束读取玩家的简单信息[%d],timecost[%d]", startTime, endTime-startTime)
}

func (this *PlayerSimpleMgr) GetPlayerDataByName(name string) *SimplePlayerData{
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

func (this *PlayerSimpleMgr) GetPlayerDataById(playerId int) *SimplePlayerData{
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

func (this *PlayerSimpleMgr) GetPlayerName(playerId int) string{
	pData := this.GetPlayerDataById(playerId)
	if pData != nil {
		return pData.PlayerName
	}

	return  ""
}

func LoadSimplePlayerData(playerId int) *SimplePlayerData{
	pData := new(SimplePlayerData)
	rows, err := world.SERVER.GetDB().Query(db.LoadSql(pData, "tbl_player", fmt.Sprintf("player_id =%d", playerId)))
	rs := db.Query(rows)
	if err == nil && rs.Next(){
		loadSimple(rs.Row(), pData)
		return pData
	}else if err != nil{
		common.DBERROR("LoadSimplePlayerData",err)
	}
	return nil
}

func LoadSimplePlayerDataByName(name string) *SimplePlayerData{
	pData := new(SimplePlayerData)
	var LoginTime string
	var LogoutTime string
	row := world.SERVER.GetDB().QueryRow(db.LoadSql(pData, "tbl_player", fmt.Sprintf("player_name='%s'", name)))
	if row != nil{
		err := row.Scan(&pData.AccountId, &pData.PlayerId, &pData.PlayerName, &pData.Level, &pData.Sex, &pData.Gold, &pData.DrawGold, &pData.Vip, &LoginTime, &LogoutTime)
		if err == nil {
			pData.LastLoginTime = base.GetDBTime(LoginTime).Unix()
			pData.LastLogoutTime = base.GetDBTime(LogoutTime).Unix()
			return pData
		}else{
			common.DBERROR("LoadSimplePlayerDataByName",err)
		}
	}
	return nil
}

func LoadSimplePlayerDatas(accountId int) []*SimplePlayerData{
	pList := make([]*SimplePlayerData, 0)
	nPlayerNum := 0
	pData := new(SimplePlayerData)
	rows, err := world.SERVER.GetDB().Query(db.LoadSql(pData, "tbl_player", fmt.Sprintf("account_id=%d", accountId)))
	rs := db.Query(rows)
	if err == nil{
		for rs.Next(){
			loadSimple(rs.Row(), pData)
			pList = append(pList, pData)
			nPlayerNum++
		}
	}
	return pList
}