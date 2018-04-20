package toprank

import (
	"actor"
	"database/sql"
	"base"
	"db"
	"fmt"
	"time"
	"sort"
	"server/common"
	"server/world"
)

const(
	ETopType_Start = iota
	//eTopType_PVE = eTopType_Start
	ETopType_End = iota

	sqlTable = "tbl_toprank"
	oVERDUETIME = (30*24*60*60)  //30天不上线清排行榜
)

type(
	TopRank struct{
		Id uint64	`primary`
		Type int8	`primary`
		Name string
		Score int
		Value [2]int
		LastTime int64 `datetime`
	}

	TOPRANKSET []*TopRank//排行榜队列
	TOPRANKMAP map[uint64] *TopRank//排行榜队列
	CTopMgr struct {
		actor.Actor
		m_db *sql.DB
		m_Log *base.CLog
		m_topRankMap[ETopType_End] TOPRANKMAP
		m_topRankSet[ETopType_End] TOPRANKSET
	}

	ITopMgr interface {
		actor.IActor

		newInData(int, uint64, string, int, int, int)
		getRank(int, uint64) *TopRank
		createRank(int, uint64, string, int, int, int) *TopRank
		clear(int)
		deleteOverDue(int)
		getPlayerRank(int, int) int //获取排名

		Update()
	}
)

var(
	TOPMGR CTopMgr
)

func (this *CTopMgr) Init(num int){
	this.m_db = world.SERVER.GetDB()
	this.m_Log = world.SERVER.GetLog()
	this.Actor.Init(num)
	actor.GetGActorList().RegisterGActorList("toprank", this)
	for i := ETopType_Start; i < ETopType_End; i++{
		this.m_topRankMap[i] = make(TOPRANKMAP)
	}

	this.RegisterTimer(1000 * 1000 * 1000, this.Update)//定时器
	this.RegisterCall("InTopRank", func(caller *actor.Caller, nType int, id uint64, name string, score,val0,val1 int) {
		this.newInData(nType, id, name, score, val0, val1)
	})

	LoadDB := func() {
		pData := &TopRank{}
		this.m_Log.Println("加载排行榜")
		row, err := this.m_db.Query(db.LoadSql(pData, sqlTable, ""));
		if err != nil{
			common.DBERROR("toprank LoadDB", err)
		}
		var LastTime string
		for row.Next(){
			pData := &TopRank{}
			err := row.Scan(&pData.Id, &pData.Type, &pData.Name, &pData.Score, &pData.Value[0], &pData.Value[1], &LastTime)
			if err != nil{
				common.DBERROR("toprank LoadDB", err)
			}else{
				pData.LastTime = db.GetDBTime(LastTime).Unix()
				this.m_topRankMap[pData.Type][pData.Id] = pData
				this.m_topRankSet[pData.Type] = append(this.m_topRankSet[pData.Type], pData)
			}
		}
		this.m_Log.Println("排行榜加载完成")
	}

	LoadDB()
	this.Actor.Start()
}

func (this *CTopMgr) showTest(nType int){
	/*for i, v := range this.m_topRankMap[nType]{
		fmt.Println(i, v)
	}*/
	for i, v := range this.m_topRankSet[nType]{
		fmt.Println(i, v)
	}
}

func (this *CTopMgr) newInData(nType int, id uint64, name string, score,val0,val1 int){
	pData := this.getRank(nType, id)
	if pData == nil{
		pData = this.createRank(nType, id, name, score, val0, val1)
		this.m_db.Exec(db.InsertSql(pData, sqlTable))
		this.m_topRankMap[nType][id] = pData
		this.m_topRankSet[nType] = append(this.m_topRankSet[nType], pData)
		sort.Sort(this.m_topRankSet[nType])
		this.showTest(nType)
	}else{
		pData.Score = score
		pData.Name = name
		pData.Value[0] = val0
		pData.Value[1] = val1
		this.m_db.Exec(db.UpdateSqlEx(pData, sqlTable, "Score", "Name", "Value", "LastTime", "Id", "Type"))
		sort.Sort(this.m_topRankSet[nType])
		this.showTest(nType)
	}
}

func (this *CTopMgr) clear(nType int){
	items := this.m_topRankMap[nType]
	for i,_ := range items{
		delete(items, i)
	}
	this.m_topRankSet[nType] = make(TOPRANKSET, 0)
	this.m_db.Exec(fmt.Sprintf("delete %s where type=%d"), sqlTable, nType)
}

func (this *CTopMgr) getRank(nType int, id uint64) *TopRank{
	items := this.m_topRankMap[nType]
	pData, exist := items[id]
	if exist{
		return pData
	}
	return nil
}

func (this *CTopMgr) createRank(nType int, id uint64, name string, score,val0,val1 int) *TopRank{
	pData := &TopRank{}
	pData.Type = int8(nType)
	pData.Id = id
	pData.Name = name
	pData.Score = score
	pData.Value[0] = val0
	pData.Value[1] = val1
	return pData
}

func (this *CTopMgr) deleteOverDue(nType int){
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

func (this* CTopMgr) getPlayerRank(nType, playerId int) int{
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

func (this* CTopMgr) Update(){

}

//sort interface
func (t TOPRANKSET) Len() int{
	return len(t)
}

func (t TOPRANKSET) Less(i, j int) bool{
	return t[i].Score > t[j].Score
}

func (t TOPRANKSET)Swap(i, j int){
	t[i], t[j] = t[j], t[i]
}