package common

import (
	"time"
	"gonet/base"
)

//状态机
type (
	default_handle func()
	eventHandle func(...interface{})
	State struct{
		OnEnter default_handle
		OnLeave default_handle
		OnExpire default_handle
		OnExec default_handle
		OnTrigle eventHandle
	}

	TimerStateMachine struct{
		m_preState int
		m_curState int
		m_expireInterval int64
		m_expireTime int64
		m_maxState int
		m_isLoop bool
		m_onStateChanged default_handle
		m_states []State
	}

	ITimerStateMachine interface {
		SetStateHandle(int,interface{}, interface{}, interface{}, interface{})
		SetStateChanged(interface{})
		SetState(int, int64, bool)
		GetState() int
		GetPreState() int
		Update(int64)
		Trigger(...interface{})
	}
)

func NewCTimerStateMachine(_maxState int, _eventHandle interface{}) *TimerStateMachine{
	timeState := new(TimerStateMachine)
	timeState.m_maxState = _maxState
	timeState.m_states = make([]State, _maxState)

	if _eventHandle != nil{
		for _,v := range timeState.m_states{
			v.OnTrigle = _eventHandle.(func(...interface{}))
		}
	}
	return timeState
}

func (this *TimerStateMachine) SetStateChanged(onStateChanged interface{}){
	this.m_onStateChanged = onStateChanged.(func())
}

func (this *TimerStateMachine) SetStateHandle(index int,OnEnter interface{}, OnLeave interface{}, OnExpire interface{}, OnExec interface{}){
	if OnEnter != nil{
		this.m_states[index].OnEnter  = OnEnter.(func())
	}
	if OnLeave != nil{
		this.m_states[index].OnLeave  = OnLeave.(func())
	}
	if OnExpire != nil{
		this.m_states[index].OnExpire = OnExpire.(func())
	}
	if OnExec != nil{
		this.m_states[index].OnExec   = OnExec.(func())
	}
}

func (this *TimerStateMachine) GetState() int{
	if this.m_curState < 0{
		return 0
	}
	return this.m_curState
}

func (this *TimerStateMachine) SetState(state int, expireTime int64, isLoop bool){
	base.Assert(state >= 0 && state < this.m_maxState,"invalid state")

	if state >= this.m_maxState{
		return
	}

	if -1 == this.m_curState && this.m_curState != state{
		s := this.m_states[this.m_curState]
		if s.OnLeave != nil{
			s.OnLeave()
		}
	}

	this.m_preState = this.m_curState
	this.m_curState = state

	s := this.m_states[state]
	if s.OnEnter != nil{
		s.OnEnter()
	}

	this.m_isLoop = isLoop
	this.m_expireTime = expireTime
	this.m_expireInterval = expireTime

	if (0 != expireTime){
		expireTime += time.Now().Unix()
	}
}

//获取前面的一个状态
func (this *TimerStateMachine) GetPreState()int{
	return this.m_preState
}

func (this *TimerStateMachine) Update(curTime int64){
	if (-1 == this.m_curState){
		return
	}

	if (0 != this.m_expireTime && curTime > this.m_expireTime){
		s := this.m_states[this.m_curState]

		//先重新设置超时的时间，然后切换状态
		this.m_expireTime =0
		if(this.m_isLoop){
			this.m_expireTime = curTime + this.m_expireInterval
		}

		if (nil != s.OnExpire){
			 s.OnExpire()
		}
	}else{
		s := this.m_states[this.m_curState]
		if (nil != s.OnExec){
			s.OnExec()
		}
	}
}

func (this *TimerStateMachine) Trigger(params ...interface{}){
	if(this.m_curState < 0){
		return
	}

	s := this.m_states[this.m_curState]
	if (nil != s.OnTrigle){
		s.OnTrigle(params...)
	}
}
