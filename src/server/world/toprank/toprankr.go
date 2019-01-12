package toprank

import (
	"actor"
	"base"
	"database/sql"
	"db"
	"fmt"
	"rd"
	"server/common"
	"server/world"
)

type(
	TopMgrR struct {
		actor.Actor

		m_db *sql.DB
		m_Log *base.CLog
		m_topRankTimer *common.SimpleTimer
	}
)

func getZRdKey(nType int) string{
	return fmt.Sprintf("z_%s_%d", sqlTable, nType)
}

func getHRdKey(nType int) string{
	return fmt.Sprintf("h_%s_%d", sqlTable, nType)
}

func (this *TopMgrR) loadDB(nType int) {
	this.m_Log.Println("读取排行榜")
	result := new(int64)
	rd.Query(world.RdID, "ZCARD", getZRdKey(nType), result)
	if *result == 0{
		fmt.Println(db.LoadSql(&TopRank{}, sqlTable, fmt.Sprintf("type = %d order by `score` limit 0, %d", nType,TOP_RANK_MAX)))
		rows, err := this.m_db.Query(db.LoadSql(&TopRank{}, sqlTable, fmt.Sprintf("type = %d order by `score` limit 0, %d", nType,TOP_RANK_MAX)));
		if err != nil{
			common.DBERROR("toprankr LoadDB", err)
		}
		rs := db.Query(rows)
		topList := make([]*TopRank, 0)
		rs.Obj(&topList)
		sdata := []interface{}{}
		data  := []interface{}{}
		for _, v := range topList{
			sdata = append(sdata, v.Score, v.Id)
			data  = append(data, v.Id, v)
		}

		rd.ExecKV(-1, world.RdID, "ZADD", getZRdKey(nType), sdata...)
		rd.ExecKV(-1, world.RdID, "HMSET", getHRdKey(nType), data...)
	}

	this.m_Log.Println("读取排行榜加载完成")
}

//分布式考虑直接数据库
func (this *TopMgrR) Init(num int){
	this.m_db = world.SERVER.GetDB()
	this.m_Log = world.SERVER.GetLog()
	this.m_topRankTimer = common.NewSimpleTimer(TOP_RANK_SYNC_TIME)
	this.Actor.Init(num)
	actor.MGR().AddActor(this)

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
	rd.ExecKV(-1, world.RdID, "ZADD", getZRdKey(nType), score, id)
	rd.ExecKV(-1, world.RdID, "HMSET", getHRdKey(nType), id, pData)
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
	rd.Exec(-1, world.RdID, "DEL", getZRdKey(nType))
	rd.Exec(-1, world.RdID, "DEL", getHRdKey(nType))
	this.m_db.Exec(fmt.Sprintf("delete %s where type=%d"), sqlTable, nType)
}

func (this *TopMgrR) getRank(nType int, id int64) *TopRank{
	pData := &TopRank{}
	if rd.Query(world.RdID, "HGET", getHRdKey(nType), pData, id) == nil{
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
	rank := new(int64)
	if rd.Query(world.RdID, "ZRANGE", getZRdKey(nType), rank) == nil{
		return int(*rank)
	}
	return -1
}

func (this* TopMgrR) update(){
	//每隔一定时间同步sql的数据
}