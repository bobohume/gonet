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
		m_Locker sync.Locker
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

func (this *PlayerSimpleMgr) Init(num int) {
	this.Actor.Init(num)
	//注册结构体
	base.RegisterMessage(&SimplePlayerData{})

	this.m_Locker = &sync.Mutex{}
	this.m_db = world.SERVER.GetDB()
	this.m_Log = world.SERVER.GetLog()
	this.m_SimplePlayerMap = make(map[int] *SimplePlayerData)
	this.m_SimplePlayerNameMap = make(map[string] *SimplePlayerData)

	this.Actor.Start()
}

func (this *PlayerSimpleMgr) LoadSimplePlayerDatas() {
	startTime := time.Now().Unix()
	var simpledata SimplePlayerData
	var LoginTime, LogoutTime string
	row, err := this.m_db.Query(db.LoadSql(simpledata, "tbl_player", ""));
	if err != nil{
		common.DBERROR("LoadSimplePlayerDatas", err)
	}
	for row.Next(){
		pData := &SimplePlayerData{}
		err := row.Scan(&pData.AccountId, &pData.PlayerId, &pData.PlayerName, &pData.Level, &pData.Sex, &pData.Gold, &pData.DrawGold, &pData.Vip, &LoginTime, &LogoutTime)
		if err != nil{
			common.DBERROR("LoadSimplePlayerDatas", err)
		}else{
			pData.LastLoginTime = db.GetDBTime(LoginTime).Unix()
			pData.LastLogoutTime = db.GetDBTime(LogoutTime).Unix()
			this.m_Locker.Lock()
			this.m_SimplePlayerMap[pData.PlayerId] = pData
			this.m_SimplePlayerNameMap[pData.PlayerName] = pData
			this.m_Locker.Unlock()
		}
	}

	endTime := time.Now().Unix()
	this.m_Log.Println("结束读取玩家的简单信息[%d],timecost[%d]", startTime, endTime-startTime)
}

func (this *PlayerSimpleMgr) GetPlayerDataByName(name string) *SimplePlayerData{
	this.m_Locker.Lock()
	pData, exist := this.m_SimplePlayerNameMap[name]
	this.m_Locker.Unlock()
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
	this.m_Locker.Lock()
	pData, exist := this.m_SimplePlayerMap[playerId]
	this.m_Locker.Unlock()
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
	var LoginTime string
	var LogoutTime string
	row := world.SERVER.GetDB().QueryRow(db.LoadSql(pData, "tbl_player", fmt.Sprintf("playerId =%d", playerId)))
	if row != nil{
		err := row.Scan(&pData.AccountId, &pData.PlayerId, &pData.PlayerName, &pData.Level, &pData.Sex, &pData.Gold, &pData.DrawGold, &pData.Vip, &LoginTime, &LogoutTime)
		if err == nil {
			pData.LastLoginTime = db.GetDBTime(LoginTime).Unix()
			pData.LastLogoutTime = db.GetDBTime(LogoutTime).Unix()
			return pData
		}else{
			common.DBERROR("LoadSimplePlayerData",err)
		}
	}
	return nil
}

func LoadSimplePlayerDataByName(name string) *SimplePlayerData{
	pData := new(SimplePlayerData)
	var LoginTime string
	var LogoutTime string
	row := world.SERVER.GetDB().QueryRow(db.LoadSql(pData, "tbl_player", fmt.Sprintf("playerName='%s'", name)))
	if row != nil{
		err := row.Scan(&pData.AccountId, &pData.PlayerId, &pData.PlayerName, &pData.Level, &pData.Sex, &pData.Gold, &pData.DrawGold, &pData.Vip, &LoginTime, &LogoutTime)
		if err == nil {
			pData.LastLoginTime = db.GetDBTime(LoginTime).Unix()
			pData.LastLogoutTime = db.GetDBTime(LogoutTime).Unix()
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
	var LoginTime string
	var LogoutTime string
	pData := new(SimplePlayerData)
	rows, err := world.SERVER.GetDB().Query(db.LoadSql(pData, "tbl_player", fmt.Sprintf("accountId=%d", accountId)))
	if err == nil{
		for rows.Next(){
			err := rows.Scan(&pData.AccountId, &pData.PlayerId, &pData.PlayerName, &pData.Level, &pData.Sex, &pData.Gold, &pData.DrawGold, &pData.Vip, &LoginTime, &LogoutTime)
			if err == nil {
				pData.LastLoginTime = db.GetDBTime(LoginTime).Unix()
				pData.LastLogoutTime = db.GetDBTime(LogoutTime).Unix()
				pList = append(pList, pData)
				nPlayerNum++
			}else{
				common.DBERROR("LoadSimplePlayerDatas", err)
			}
		}
	}
	return pList
}