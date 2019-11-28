package toprank

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"gonet/actor"
	"gonet/base"
	"gonet/db"
	"gonet/rd"
	"gonet/server/common"
	"gonet/server/world"
)

type(
	TopMgrR struct {
		actor.Actor

		m_db *sql.DB
		m_Log *base.CLog
		m_topRankTimer *common.SimpleTimer
	}
)

func ZRdKey(nType int) string{
	return fmt.Sprintf("z_%s_%d", sqlTable, nType)
}

func HRdKey(nType int) string{
	return fmt.Sprintf("h_%s_%d", sqlTable, nType)
}

func (this *TopMgrR) loadDB(nType int) {
	this.m_Log.Println("读取排行榜")
	result := int64(0)
	rd.Do(world.RdID, func(c redis.Conn) {
		result, _ = redis.Int64(c.Do("ZCARD", ZRdKey(nType)))
	})
	if result == 0{
		fmt.Println(db.LoadSql(&TopRank{}, sqlTable, fmt.Sprintf("type = %d order by `score` limit 0, %d", nType,TOP_RANK_MAX)))
		rows, err := this.m_db.Query(db.LoadSql(&TopRank{}, sqlTable, fmt.Sprintf("type = %d order by `score` limit 0, %d", nType,TOP_RANK_MAX)));
		if err != nil{
			common.DBERROR("toprankr LoadDB", err)
		}
		rs := db.Query(rows, err)
		topList := make([]*TopRank, 0)
		rs.Obj(&topList)
		for _, v := range topList{
			rd.Do(world.RdID, func(c redis.Conn) {
				c.Send("ZADD", ZRdKey(nType), v.Score, v.Id)
				data, _ := json.Marshal(v)
				c.Send("HSET", HRdKey(nType), v.Id, data)
				c.Flush()
			})
		}
	}

	this.m_Log.Println("读取排行榜加载完成")
}

//分布式考虑直接数据库
func (this *TopMgrR) Init(num int){
	this.m_db = world.SERVER.GetDB()
	this.m_Log = world.SERVER.GetLog()
	this.m_topRankTimer = common.NewSimpleTimer(TOP_RANK_SYNC_TIME)
	this.Actor.Init(num)
	actor.MGR.AddActor(this)

	this.RegisterTimer(1000 * 1000 * 1000, this.update)//定时器
	this.RegisterCall("InTopRank", func(nType int, id int64, name string, score,val0,val1 int) {
		this.newInData(nType, id, name, score, val0, val1)
	})

	for i := ETopType_Start; i < ETopType_End; i++{
		this.loadDB(i)
	}

	this.Actor.Start()
}

func (this *TopMgrR) newInData(nType int, id int64, name string, score,val0,val1 int){
	pData := this.createRank(nType, id, name, score, val0, val1)
	rd.Do(world.RdID, func(c redis.Conn) {
		c.Send("ZADD", ZRdKey(nType), score, id)
		data, _ := json.Marshal(pData)
		c.Send("HSET", HRdKey(nType), id, data)
	})
	bExist := false
	row := this.m_db.QueryRow(fmt.Sprintf("select 1 from %s where id=%d and type=%d", sqlTable, id, nType))
	if row != nil{
		bExist = true
	}

	if bExist{
		this.m_db.Exec(db.UpdateSqlEx(pData, sqlTable, "score", "name", "value", "last_time", "id", "type"))
	}else{
		this.m_db.Exec(db.InsertSql(pData, sqlTable))
	}
}

func (this *TopMgrR) clearTop(nType int){
	rd.Do(world.RdID, func(c redis.Conn) {
		c.Send("DEL", ZRdKey(nType))
		c.Send("DEL", HRdKey(nType))
		c.Flush()
	})
	this.m_db.Exec(fmt.Sprintf("delete %s where type=%d", sqlTable, nType))
}

func (this *TopMgrR) getRank(nType int, id int64) *TopRank{
	pData := &TopRank{}
	data := []byte{}
	rd.Do(world.RdID, func(c redis.Conn) {
		data, _ = redis.Bytes(c.Do("HGET", HRdKey(nType), id))
	})

	if json.Unmarshal(data, pData) == nil{
		return pData
	}

	return nil
}

func (this *TopMgrR) createRank(nType int, id int64, name string, score,val0,val1 int) *TopRank{
	pData := &TopRank{}
	pData.Type = int8(nType)
	pData.Id = id
	pData.Name = name
	pData.Score = score
	pData.Value[0] = val0
	pData.Value[1] = val1
	return pData
}

func (this* TopMgrR) getPlayerRank(nType int, playerId int64) int{
	rank := int(-1)
	rd.Do(world.RdID, func(c redis.Conn) {
		rank, _ =redis.Int(c.Do("ZREVRANK", ZRdKey(nType), playerId))
	})
	return rank
}

func (this* TopMgrR) update(){
	//每隔一定时间同步sql的数据
}