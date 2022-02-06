package toprank

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"gonet/actor"
	"gonet/base"
	"gonet/common"
	"gonet/db"
	"gonet/rd"
	"gonet/server/model"
	"gonet/server/world"

	"github.com/gomodule/redigo/redis"
)

type (
	TopMgrR struct {
		actor.Actor

		m_db           *sql.DB
		m_Log          *base.CLog
		m_topRankTimer *common.SimpleTimer
	}
)

func ZRdKey(nType int) string {
	return fmt.Sprintf("z_%s_%d", sqlTable, nType)
}

func HRdKey(nType int) string {
	return fmt.Sprintf("h_%s_%d", sqlTable, nType)
}

func (this *TopMgrR) loadDB(nType int) {
	this.m_Log.Println("读取排行榜")
	result, _ := redis.Int64(rd.Do(world.RdID, func(c redis.Conn) (reply interface{}, err error) {
		return c.Do("ZCARD", ZRdKey(nType))
	}))
	if result == 0 {
		fmt.Println(db.LoadSql(&model.TopRank{}, db.WithWhere(&model.TopRank{Type:int8(nType)}), db.WithLimit(TOP_RANK_MAX)))
		rows, err := this.m_db.Query(db.LoadSql(&model.TopRank{}, db.WithWhere(&model.TopRank{Type:int8(nType)}), db.WithLimit(TOP_RANK_MAX)))
		if err != nil {
			common.DBERROR("toprankr LoadDB", err)
		}
		rs := db.Query(rows, err)
		topList := make([]*model.TopRank, 0)
		rs.Obj(&topList)
		for _, v := range topList {
			data, _ := json.Marshal(v)
			rd.Do(world.RdID, func(c redis.Conn) (reply interface{}, err error) {
				c.Send("ZADD", ZRdKey(nType), v.Score, v.Id)
				c.Send("HSET", HRdKey(nType), v.Id, data)
				c.Flush()
				return nil, nil
			})
		}
	}

	this.m_Log.Println("读取排行榜加载完成")
}

//分布式考虑直接数据库
func (this *TopMgrR) Init() {
	this.Actor.Init()
	this.m_db = world.SERVER.GetDB()
	this.m_Log = world.SERVER.GetLog()
	this.m_topRankTimer = common.NewSimpleTimer(TOP_RANK_SYNC_TIME)
	this.RegisterTimer(1000*1000*1000, this.update) //定时器
	for i := ETopType_Start; i < ETopType_End; i++ {
		this.loadDB(i)
	}
	actor.MGR.RegisterActor(this)
	this.Actor.Start()
}

func (this *TopMgrR) newInData(nType int, id int64, name string, score, val0, val1 int) {
	pData := this.createRank(nType, id, name, score, val0, val1)
	data, _ := json.Marshal(pData)
	rd.Do(world.RdID, func(c redis.Conn) (reply interface{}, err error) {
		c.Send("ZADD", ZRdKey(nType), score, id)
		c.Send("HSET", HRdKey(nType), id, data)
		return nil, nil
	})
	bExist := false
	row := this.m_db.QueryRow(fmt.Sprintf("select 1 from %s where id=%d and type=%d", sqlTable, id, nType))
	if row != nil {
		bExist = true
	}

	if bExist {
		this.m_db.Exec(db.UpdateSql(pData))
	} else {
		this.m_db.Exec(db.InsertSql(pData))
	}
}

func (this *TopMgrR) clearTop(nType int) {
	rd.Do(world.RdID, func(c redis.Conn) (reply interface{}, err error) {
		c.Send("DEL", ZRdKey(nType))
		c.Send("DEL", HRdKey(nType))
		c.Flush()
		return nil, nil
	})
	this.m_db.Exec(fmt.Sprintf("delete %s where type=%d", sqlTable, nType))
}

func (this *TopMgrR) getRank(nType int, id int64) *model.TopRank {
	pData := &model.TopRank{}
	data, _ := redis.Bytes(rd.Do(world.RdID, func(c redis.Conn) (reply interface{}, err error) {
		return c.Do("HGET", HRdKey(nType), id)
	}))

	if json.Unmarshal(data, pData) == nil {
		return pData
	}

	return nil
}

func (this *TopMgrR) createRank(nType int, id int64, name string, score, val0, val1 int) *model.TopRank {
	pData := &model.TopRank{}
	pData.Type = int8(nType)
	pData.Id = id
	pData.Name = name
	pData.Score = score
	pData.Value[0] = val0
	pData.Value[1] = val1
	return pData
}

func (this *TopMgrR) getPlayerRank(nType int, playerId int64) int {
	rank, _ := redis.Int(rd.Do(world.RdID, func(c redis.Conn) (reply interface{}, err error) {
		return c.Do("ZREVRANK", ZRdKey(nType), playerId)
	}))
	return rank
}

func (this *TopMgrR) update() {
	//每隔一定时间同步sql的数据
}

func (this *TopMgrR) InTopRank(ctx context.Context, nType int, id int64, name string, score, val0, val1 int) {
	this.newInData(nType, id, name, score, val0, val1)
}
