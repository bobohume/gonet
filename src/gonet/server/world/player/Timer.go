package player

import (
	"gonet/base"
	"gonet/db"
	"gonet/server/world/common"
	"sort"
	"time"
)

type TIMERID int
const(
	TIMER_ID1  TIMERID = iota
)

type(
	Timer struct {
		Id int `sql:"primary;name:id"`						//定时器Id
		PlayerId int64 `sql:"primary;name:player_id"`		//玩家Id
		Flag int64	`sql:"name:flag"`						//定时器数据
		ExpireTime int64 `sql:datetime;name:expire_time`	//定时器过期时间
	}

	TimerSet struct {
		base.Vector
	}//定时器序列

	TimerMgr struct {
		common.IPlayer
		m_TimerSet TimerSet//定时器set
		m_TimerMap map[int] *Timer//方便查找
	}

	ITimerMgr interface {
		Init(common.IPlayer)
		GetTimer(Id int) *Timer//获取定时器
		AddTimer(Id int, Flag int64, ExpireTime int64)//添加定时器
		DelTimer(Id int)//删除定时器
		Update()
		sort()
	}
)

func (this *TimerMgr) Init(pPlayer common.IPlayer){
	this.IPlayer = pPlayer
	this.m_TimerMap = map[int] *Timer{}
}

func (this *TimerMgr)  GetTimer(Id int) *Timer{
	pTimer, bEx := this.m_TimerMap[Id]
	if bEx{
		return pTimer
	}
	return nil
}

func (this *TimerMgr) AddTimer(Id int, Flag int64, ExpireTime int64){
	pTimer, bEx := this.m_TimerMap[Id]
	if bEx && pTimer != nil{
		pTimer.Flag = Flag
		pTimer.ExpireTime = ExpireTime
		this.sort()
		this.GetDB().Exec(db.UpdateSqlEx(pTimer, "tbl_timerset", "flag", "expire_time"))
	}else{
		pTimer = &Timer{}
		pTimer.Id = Id
		pTimer.PlayerId = this.GetPlayerId()
		pTimer.ExpireTime = ExpireTime
		pTimer.Flag = Flag
		this.m_TimerSet.Push_back(pTimer)
		this.m_TimerMap[Id] = pTimer
		this.sort()
		this.GetDB().Exec(db.InsertSql(pTimer, "tbl_timerset"))
	}
}

func (this *TimerMgr) DelTimer(Id int){
	_, bEx := this.m_TimerMap[Id]
	if bEx{
		for i, v := 0, this.m_TimerSet.Begin(); v != this.m_TimerSet.End(); v = this.m_TimerSet.Next(&i) {
			if v != nil && (*v).(*Timer).Id == Id{
				delete(this.m_TimerMap, Id)
				this.m_TimerSet.Erase(i)
				i--
				this.GetDB().Exec(db.DeleteSql((*v).(*Timer), "tbl_timerset"))
				break
			}
		}
	}
}

func (this *TimerMgr) Update() {
	nCurTime := time.Now().Unix()
	//定时器排过期时间排序
	for i, v := 0, this.m_TimerSet.Begin(); v != this.m_TimerSet.End(); v = this.m_TimerSet.Next(&i) {
		if v != nil && (*v).(*Timer).ExpireTime <= nCurTime{//活动过期
			delete(this.m_TimerMap, (*v).(*Timer).Id)
			this.m_TimerSet.Erase(i)
			i--
			this.GetDB().Exec(db.DeleteSql((*v).(*Timer), "tbl_timerset"))
			continue
		}else{
			break
		}
	}
}

func (this *TimerMgr) sort(){
	aa := this.m_TimerSet
	sort.Sort(&aa)
	sort.Sort(&this.m_TimerSet)
}

//sort interface
func (t *TimerSet) Less(i, j int) bool{
	return t.Get(i).(*Timer).ExpireTime < t.Get(j).(*Timer).ExpireTime
}