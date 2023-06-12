package base

// 状态机
type (
	FsmHandle func()
	State     struct {
		OnEnter FsmHandle
		OnExec  FsmHandle
		OnExit  FsmHandle
	}

	Fsm struct {
		preState int
		curState int
		maxState int
		states   []State
	}

	IFsm interface {
		Init(int)
		SetStateHandle(state int, st *State)
		SetState(int)
		GetState() int
		GetPreState() int
		Update()
	}
)

func (f *Fsm) Init(_maxState int) {
	f.curState = 0
	f.maxState = _maxState
	f.states = make([]State, _maxState)
}

func (f *Fsm) SetStateHandle(state int, st *State) {
	if st == nil {
		return
	}

	if st.OnEnter != nil {
		f.states[state].OnEnter = st.OnEnter
	}
	if st.OnExec != nil {
		f.states[state].OnExec = st.OnExec
	}
	if st.OnExit != nil {
		f.states[state].OnExit = st.OnExit
	}
}

func (f *Fsm) GetState() int {
	return f.curState
}

func (f *Fsm) SetState(state int) {
	Assert(state >= 0 && state < f.maxState, "invalid state")

	if state >= f.maxState {
		return
	}

	if f.curState != state {
		s := f.states[f.curState]
		if s.OnExit != nil {
			s.OnExit()
		}
	}

	f.preState = f.curState
	f.curState = state

	s := f.states[state]
	if s.OnEnter != nil {
		s.OnEnter()
	}
}

// 获取前面的一个状态
func (f *Fsm) GetPreState() int {
	return f.preState
}

func (f *Fsm) Update() {
	s := f.states[f.curState]
	if nil != s.OnExec {
		s.OnExec()
	}
}
