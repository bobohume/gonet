package common

import (
	"gonet/base"
)

//状态机
type (
	TimerHandle func()
	State       struct {
		OnEnter TimerHandle
		OnExit  TimerHandle
		OnExec  TimerHandle
	}

	StateMachine struct {
		preState int
		curState int
		maxState int
		states   []State
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

func (this *StateMachine) Init(_maxState int) {
	this.curState = 0
	this.maxState = _maxState
	this.states = make([]State, _maxState)
}

func (this *StateMachine) SetStateHandle(state int, st *State) {
	if st == nil {
		return
	}

	if st.OnEnter != nil {
		this.states[state].OnEnter = st.OnEnter
	}
	if st.OnExec != nil {
		this.states[state].OnExec = st.OnExec
	}
	if st.OnExit != nil {
		this.states[state].OnExit = st.OnExit
	}
}

func (this *StateMachine) GetState() int {
	return this.curState
}

func (this *StateMachine) SetState(state int) {
	base.Assert(state >= 0 && state < this.maxState, "invalid state")

	if state >= this.maxState {
		return
	}

	if this.curState != state {
		s := this.states[this.curState]
		if s.OnExit != nil {
			s.OnExit()
		}
	}

	this.preState = this.curState
	this.curState = state

	s := this.states[state]
	if s.OnEnter != nil {
		s.OnEnter()
	}
}

//获取前面的一个状态
func (this *StateMachine) GetPreState() int {
	return this.preState
}

func (this *StateMachine) Update() {
	s := this.states[this.curState]
	if nil != s.OnExec {
		s.OnExec()
	}
}
