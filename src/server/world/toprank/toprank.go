package toprank

import (
	"actor"
	"db"
	"fmt"
	"time"
	"sort"
	"server/common"
	"server/world"
	"database/sql"
	"base"
)

const(
	ETopType_Start = iota
	//eTopType_PVE = eTopType_Start
	ETopType_End = iota

	sqlTable = "tbl_toprank"
	oVERDUETIME = (30*24*60*60)  //30天不上线清排行榜
	TOP_RANK_MAX  = 100
	TOP_RANK_SYNC_TIME = 3 * 60
)

type(
	TopRank struct{
		Id uint64	`sql:"primary;name:id"`
		Type int8	`sql:"primary;name:type"`
		Name string `sql:"name:name"`
		Score int `sql:"name:score"`
		Value [2]int `sql:"name:value"`
		LastTime int64 `sql:"datetime;name:last_time"`
	}

	TOPRANKSET []*TopRank//排行榜队列
	TOPRANKMAP map[uint64] *TopRank//排行榜队列
	TopMgr struct {
		actor.Actor

		m_db *sql.DB
		m_Log *base.CLog
		m_topRankMap[ETopType_End] TOPRANKMAP
		m_topRankSet[ETopType_End] TOPRANKSET
		m_topRankTimer *common.SimpleTimer
	}

	ITopMgr interface {
		actor.IActor

		loadDB()
		newInData(int, uint64, string, int, int, int)
		getRank(int, uint64) *TopRank
		createRank(int, uint64, string, int, int, int) *TopRank
		clearTop(int)
		deleteOverDue(int)
		getPlayerRank(int, int) int //获取排名
		clearRank()
		Update()
	}
)

var(
	TOPMGR TopMgr
)

func loadTopRank(row db.IRow, t *TopRank){
	t.Id = uint64(row.Int64("id"))
	t.Type = int8(row.Int("type"))
	t.Name = row.String("name")
	t.Score = row.Int("score")
	t.Value[0] = row.Int("value0")
	t.Value[1] = row.Int("value1")
	t.LastTime = row.Time("last_time")
}

func (this *TopMgr) loadDB() {
	//pData := &TopRank{}
	this.m_Log.Println("读取排行榜")
	this.clearRank()
	fmt.Sprintf(fmt.Sprintf("select * from %s order by `score` limit 0, %d", sqlTable, TOP_RANK_MAX))
	rows, err := this.m_db.Query(fmt.Sprintf("select * from %s order by `score` limit 0, %d", sqlTable, TOP_RANK_MAX));
	//row, err := this.m_db.Query(db.LoadSql(pData, sqlTable, ""));
	if err != nil{
		common.DBERROR("toprank LoadDB", err)
	}
	rs := db.Query(rows)
	for rs.Next(){
		pData := &TopRank{}
		loadTopRank(rs.Row(), pData)
		this.m_topRankMap[pData.Type][pData.Id] = pData
		this.m_topRankSet[pData.Type] = append(this.m_topRankSet[pData.Type], pData)
	}
	this.m_Log.Println("读取排行榜加载完成")
}

//分布式考虑直接数据库
func (this *TopMgr) Init(num int){
	this.m_db = world.SERVER.GetDB()
	this.m_Log = world.SERVER.GetLog()
	this.m_topRankTimer = common.NewSimpleTimer(TOP_RANK_SYNC_TIME)
	this.Actor.Init(num)
	actor.GetGActorList().RegisterGActorList("toprank", this)
	this.clearRank()

	this.RegisterTimer(1000 * 1000 * 1000, this.Update)//定时器
	this.RegisterCall("InTopRank", func(nType int, id uint64, name string, score,val0,val1 int) {
		this.newInData(nType, id, name, score, val0, val1)
	})

	this.loadDB()
	this.Actor.Start()
}

func (this *TopMgr) clearRank(){
	for i := ETopType_Start; i < ETopType_End; i++{
		this.m_topRankMap[i] = make(TOPRANKMAP)
		this.m_topRankSet[i] = make(TOPRANKSET, 0)
	}
}

func (this *TopMgr) showTest(nType int){
	/*for i, v := range this.m_topRankMap[nType]{
		fmt.Println(i, v)
	}*/
	for i, v := range this.m_topRankSet[nType]{
		base.GLOG.Println(i, v)
	}
}

func (this *TopMgr) newInData(nType int, id uint64, name string, score,val0,val1 int){
	pData := this.getRank(nType, id)
	if pData == nil{
		pData = this.createRank(nType, id, name, score, val0, val1)
		this.m_topRankMap[nType][id] = pData
		this.m_topRankSet[nType] = append(this.m_topRankSet[nType], pData)
		sort.Sort(&this.m_topRankSet[nType])
		this.showTest(nType)
	}else{
		pData.Score = score
		pData.Name = name
		pData.Value[0] = val0
		pData.Value[1] = val1
		sort.Sort(&this.m_topRankSet[nType])
		this.showTest(nType)
	}

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

func (this *TopMgr) clearTop(nType int){
	items := this.m_topRankMap[nType]
	for i,_ := range items{
		delete(items, i)
	}
	this.m_topRankSet[nType] = make(TOPRANKSET, 0)
	this.m_db.Exec(fmt.Sprintf("delete %s where type=%d"), sqlTable, nType)
}

func (this *TopMgr) getRank(nType int, id uint64) *TopRank{
	items := this.m_topRankMap[nType]
	pData, exist := items[id]
	if exist{
		return pData
	}
	return nil
}

func (this *TopMgr) createRank(nType int, id uint64, name string, score,val0,val1 int) *TopRank{
	pData := &TopRank{}
	pData.Type = int8(nType)
	pData.Id = id
	pData.Name = name
	pData.Score = score
	pData.Value[0] = val0
	pData.Value[1] = val1
	return pData
}

func (this *TopMgr) deleteOverDue(nType int){
	isNeedOverDue := func() bool{
		return false
	}

	if !isNeedOverDue(){
		return
	}

	curtime := time.Now().Unix()
	items := this.m_topRankMap[nType]
	for i,v := range items{
		pData := v
		if pData != nil && curtime - pData.LastTime >= oVERDUETIME {
			delete(items, i)
			this.m_Log.Printf("删除过期的排行榜项[type]=%d,[uid]=%s", nType, pData.Id)
			this.m_db.Exec(db.DeleteSql(pData, sqlTable))
		}
	}

	for i,v := range this.m_topRankSet[nType]{
		pData := v
		if pData != nil && curtime - pData.LastTime >= oVERDUETIME {
			this.m_topRankSet[nType] = append(this.m_topRankSet[nType][:i], this.m_topRankSet[nType][i+1:]...)
		}
	}
}

func (this* TopMgr) getPlayerRank(nType, playerId int) int{
	Id := uint64(playerId)
	pData := this.getRank(nType, Id)
	if pData != nil{
		items := this.m_topRankSet[nType]
		index := sort.Search(items.Len(), func(i int) bool {
			return items[i].Score <= pData.Score
		})

		if index <= items.Len() && items[index].Score == pData.Score{
			for i := index; i < items.Len(); i++{
				if(items[i].Id == Id){
					return  i
				}else if(items[i].Score != pData.Score){
					break
				}
			}
		}
	}
	return -1
}

func (this* TopMgr) Update(){
	//每隔一定时间同步sql的数据
	if this.m_topRankTimer.CheckTimer(){
		this.loadDB()
	}
}

//sort interface
func (t *TOPRANKSET) Len() int{
	return len(*t)
}

func (t *TOPRANKSET) Less(i, j int) bool{
	return (*t)[i].Score > (*t)[j].Score
}

func (t *TOPRANKSET)Swap(i, j int){
	(*t)[i], (*t)[j] = (*t)[j], (*t)[i]
}