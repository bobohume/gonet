package gm

import (
	"gonet/actor"
	"gonet/base"
	"gonet/common"
	"gonet/orm"
	"gonet/server/model"
	"sync"
	"time"
)

type (
	PlayerSimpleMgr struct {
		actor.Actor
		simplePlayerMap     map[int64]*model.SimplePlayerData
		simplePlayerNameMap map[string]*model.SimplePlayerData
		locker              *sync.RWMutex
	}

	IPlayerSimpleMgr interface {
		actor.IActor

		LoadSimplePlayerDatas()
		GetPlayerDataByName(string) *model.SimplePlayerData
		GetPlayerDataById(int64) *model.SimplePlayerData
		GetPlayerName(int64) string
	}
)

var (
	SIMPLEMGR PlayerSimpleMgr
)

func loadSimple(row orm.IRow, s *model.SimplePlayerData) {
	s.PlayerId = row.Int64("player_id")
	s.PlayerName = row.String("player_name")
	s.AccountId = row.Int64("account_id")
	s.Level = row.Int("level")
	s.Sex = row.Int("sex")
	s.Gold = row.Int("gold")
	s.DrawGold = row.Int("draw_gold")
	s.Vip = row.Int("vip")
	s.LastLoginTime = row.Time("last_login_time")
	s.LastLogoutTime = row.Time("last_logout_time")
}

func (p *PlayerSimpleMgr) Init() {
	p.Actor.Init()
	p.locker = &sync.RWMutex{}
	p.simplePlayerMap = make(map[int64]*model.SimplePlayerData)
	p.simplePlayerNameMap = make(map[string]*model.SimplePlayerData)
	actor.MGR.RegisterActor(p)
	p.Actor.Start()
}

func (p *PlayerSimpleMgr) LoadSimplePlayerDatas() {
	startTime := time.Now().Unix()
	var simpledata model.SimplePlayerData
	rows, err := orm.DB.Query(orm.LoadSql(simpledata, orm.WithOutWhere()))
	if err != nil {
		common.DBERROR("LoadSimplePlayerDatas", err)
	}
	rs, err := orm.Query(rows, err)
	for err == nil && rs.Next() {
		data := &model.SimplePlayerData{}
		loadSimple(rs.Row(), data)
		//rs.Row().Obj(data)
		p.locker.Lock()
		p.simplePlayerMap[data.PlayerId] = data
		p.simplePlayerNameMap[data.PlayerName] = data
		p.locker.Unlock()
	}

	endTime := time.Now().Unix()
	base.LOG.Printf("结束读取玩家的简单信息[%d],timecost[%d]", startTime, endTime-startTime)
}

func (p *PlayerSimpleMgr) GetPlayerDataByName(name string) *model.SimplePlayerData {
	p.locker.RLock()
	data, exist := p.simplePlayerNameMap[name]
	p.locker.RUnlock()
	if exist {
		return data
	}

	data = LoadSimplePlayerDataByName(name)
	if data != nil {
		p.locker.Lock()
		p.simplePlayerMap[data.PlayerId] = data
		p.simplePlayerNameMap[name] = data
		p.locker.Unlock()
	}

	return data
}

func (p *PlayerSimpleMgr) GetPlayerDataById(playerId int64) *model.SimplePlayerData {
	p.locker.RLock()
	data, exist := p.simplePlayerMap[playerId]
	p.locker.RUnlock()
	if exist {
		return data
	}

	data = LoadSimplePlayerData(playerId)
	if data != nil {
		p.locker.Lock()
		p.simplePlayerMap[data.PlayerId] = data
		p.simplePlayerNameMap[data.PlayerName] = data
		p.locker.Unlock()
	}

	return data
}

func (p *PlayerSimpleMgr) GetPlayerName(playerId int64) string {
	data := p.GetPlayerDataById(playerId)
	if data != nil {
		return data.PlayerName
	}

	return ""
}

func LoadSimplePlayerData(playerId int64) *model.SimplePlayerData {
	data := &model.SimplePlayerData{PlayerId: playerId}
	rows, err := orm.DB.Query(orm.LoadSql(data))
	rs, err := orm.Query(rows, err)
	if err == nil && rs.Next() {
		loadSimple(rs.Row(), data)
		return data
	} else if err != nil {
		common.DBERROR("LoadSimplePlayerData", err)
	}
	return nil
}

func LoadSimplePlayerDataByName(name string) *model.SimplePlayerData {
	data := new(model.SimplePlayerData)
	rows, err := orm.DB.Query(orm.LoadSql(data, orm.WithWhere(model.SimplePlayerData{PlayerName: name})))
	rs, err := orm.Query(rows, err)
	if err == nil && rs.Next() {
		loadSimple(rs.Row(), data)
		return data
	}
	return nil
}
