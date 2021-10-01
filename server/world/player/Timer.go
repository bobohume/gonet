package player

import (
	"gonet/base/vector"
	"gonet/db"
	"gonet/server/model"
	"sort"
	"time"
)

type TIMERID int
const(
	TIMER_ID1  TIMERID = iota
)

type(
	TimerSet struct {
		vector.Vector
	}//定时器序列

	TimerMgr struct {
		IPlayer

		m_TimerSet TimerSet//定时器set
		m_TimerMap map[int] *model.Timer//方便查找
	}

	ITimerMgr interface {
		Init(IPlayer)
		GetTimer(Id int) *model.Timer//获取定时器
		AddTimer(Id int, Flag int64, ExpireTime int64)//添加定时器
		DelTimer(Id int)//删除定时器
		Update()
		sort()
	}
)

func (this *TimerMgr) Init(pPlayer IPlayer){
	this.IPlayer = pPlayer
	this.m_TimerMap = map[int] *model.Timer{}
}

func (this *TimerMgr)  GetTimer(Id int) *model.Timer{
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
		this.GetDB().Exec(db.UpdateSql(pTimer))
	}else{
		pTimer = &model.Timer{}
		pTimer.Id = Id
		pTimer.PlayerId = this.GetPlayerId()
		pTimer.ExpireTime = ExpireTime
		pTimer.Flag = Flag
		this.m_TimerSet.PushBack(pTimer)
		this.m_TimerMap[Id] = pTimer
		this.sort()
		this.GetDB().Exec(db.InsertSql(pTimer))
	}
}

func (this *TimerMgr) DelTimer(Id int){
	_, bEx := this.m_TimerMap[Id]
	if bEx{
		itr := this.m_TimerSet.Iterator()
		itr.Begin()
		for itr.Next(){
			v := itr.Value()
			if v != nil && (v).(*model.Timer).Id == Id{
				delete(this.m_TimerMap, Id)
				this.m_TimerSet.Erase(itr.Index())
				itr.Prev()
				this.GetDB().Exec(db.DeleteSql((v).(*model.Timer)))
				break
			}
		}
	}
}

func (this *TimerMgr) Update() {
	nCurTime := time.Now().Unix()
	//定时器排过期时间排序
	itr := this.m_TimerSet.Iterator()
	for itr.Next() {
		v := itr.Value()
		if v != nil && (v).(*model.Timer).ExpireTime <= nCurTime{//活动过期
			delete(this.m_TimerMap, (v).(*model.Timer).Id)
			this.m_TimerSet.Erase(itr.Index())
			itr.Prev()
			this.GetDB().Exec(db.DeleteSql((v).(*model.Timer)))
			continue
		}else{
			break
		}
	}
}

func (this *TimerMgr) sort(){
	sort.Sort(&this.m_TimerSet)
}

//sort interface
func (t *TimerSet) Less(i, j int) bool{
	return t.Get(i).(*model.Timer).ExpireTime < t.Get(j).(*model.Timer).ExpireTime
}