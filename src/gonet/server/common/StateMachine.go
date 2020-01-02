package common

import (
	"gonet/base"
)

//状态机
type (
	TimerHandle func()
	State struct{
		OnEnter TimerHandle
		OnLeave TimerHandle
		OnExec TimerHandle
	}

	StateMachine struct{
		m_preState int
		m_curState int
		m_maxState int
		m_states []State
	}

	IStateMachine interface {
		Init(int)
		SetStateHandle(state int, pState *State)
		SetState(int)
		GetState() int
		GetPreState() int
		Update()
	}
)

func (this *StateMachine) Init(_maxState int){
	this.m_curState = -1
	this.m_maxState = _maxState
	this.m_states = make([]State, _maxState)
}

func (this *StateMachine) SetStateHandle(state int, pState *State){
	if pState == nil{
		return
	}

	if pState.OnEnter != nil{
		this.m_states[state].OnEnter  = pState.OnEnter
	}
	if pState.OnExec != nil{
		this.m_states[state].OnExec   = pState.OnExec
	}
	if pState.OnLeave != nil{
		this.m_states[state].OnLeave  = pState.OnLeave
	}
}

func (this *StateMachine) GetState() int{
	if this.m_curState < 0{
		return 0
	}
	return this.m_curState
}

func (this *StateMachine) SetState(state int){
	base.Assert(state >= 0 && state < this.m_maxState,"invalid state")

	if state >= this.m_maxState{
		return
	}

	if -1 != this.m_curState && this.m_curState != state{
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
}

//获取前面的一个状态
func (this *StateMachine) GetPreState()int{
	return this.m_preState
}

func (this *StateMachine) Update(){
	if (-1 == this.m_curState){
		return
	}

	s := this.m_states[this.m_curState]
	if (nil != s.OnExec){
		s.OnExec()
	}
}