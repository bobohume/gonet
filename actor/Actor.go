package actor

import (
	"context"
	"fmt"
	"gonet/base"
	"gonet/base/cron"
	"gonet/base/mpsc"
	"gonet/base/timer"
	"gonet/rpc"
	"log"
	"reflect"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
)

var (
	g_IdSeed int64
)

const (
	ASF_NULL = iota
	ASF_RUN  = iota
	ASF_STOP = iota //已经关闭
)

// ********************************************************
// actor 核心actor模式
// ********************************************************
type (
	ActorBase struct {
		actorName string
		rType     reflect.Type
		rVal      reflect.Value
		actorType ACTOR_TYPE
		Self      IActor //when parent interface class call interface, it call parent not child  use for virtual
	}

	Actor struct {
		ActorBase
		actorChan chan int //use for states
		id        int64
		state     int32
		trace     traceInfo //trace func
		mailBox   *mpsc.Queue[*CallIO]
		mailIn    [8]int64
		mailChan  chan bool
		timerId   *int64
		pool      IActorPool //ACTOR_TYPE_VIRTUAL,ACTOR_TYPE_POOL
		timerMap  map[int64]*timerUnit
		cronMap   map[int64]*cronUnit
	}

	IActor interface {
		Init()
		Stop()
		Start()
		SendMsg(head rpc.RpcHead, funcName string, params ...interface{})
		Send(head rpc.RpcHead, packet rpc.Packet)
		RegisterTimer(duration time.Duration, fun func(), opts ...timer.OpOption) int64 //注册定时器,时间为纳秒 1000 * 1000 * 1000
		RegisterCron(cronStr string, fun func()) int64
		StopTimer(int64)
		GetId() int64
		GetState() int32
		GetRpcHead(ctx context.Context) rpc.RpcHead //rpc is safe
		GetName() string
		GetActorType() ACTOR_TYPE
		HasRpc(string) bool
		Acotr() *Actor
		register(IActor, Op)
		setState(state int32)
		bindPool(IActorPool) //ACTOR_TYPE_VIRTUAL,ACTOR_TYPE_POOL
		getPool() IActorPool //ACTOR_TYPE_VIRTUAL,ACTOR_TYPE_POOL
	}

	IActorPool interface {
		SendAcotr(head rpc.RpcHead, packet rpc.Packet) bool //ACTOR_TYPE_VIRTUAL,ACTOR_TYPE_POOL特殊判断
	}

	CallIO struct {
		rpc.RpcHead
		*rpc.RpcPacket
		Buff []byte
	}

	traceInfo struct {
		funcName  string
		fileName  string
		filePath  string
		className string
	}

	timerUnit struct {
		timer.Op
		*timer.TimerNode
		fun func()
	}

	cronUnit struct {
		cron.Schedule
		nextTime time.Time
		timerId  int64
		fun      func()
	}
)

const (
	DESDORY_EVENT = iota
)

func (a *ActorBase) IsActorType(actorType ACTOR_TYPE) bool {
	return a.actorType == actorType
}

func AssignActorId() int64 {
	return atomic.AddInt64(&g_IdSeed, 1)
}

func (a *Actor) GetId() int64 {
	return a.id
}

func (a *Actor) SetId(id int64) {
	a.id = id
}

func (a *Actor) GetName() string {
	return a.actorName
}

func (a *Actor) GetRpcHead(ctx context.Context) rpc.RpcHead {
	rpcHead := ctx.Value("rpcHead").(rpc.RpcHead)
	return rpcHead
}

func (a *Actor) GetState() int32 {
	return atomic.LoadInt32(&a.state)
}

func (a *Actor) GetActorType() ACTOR_TYPE {
	return a.actorType
}

func (a *Actor) setState(state int32) {
	atomic.StoreInt32(&a.state, state)
}

func (a *Actor) bindPool(pPool IActorPool) {
	a.pool = pPool
}

func (a *Actor) getPool() IActorPool {
	return a.pool
}

func (a *Actor) HasRpc(funcName string) bool {
	_, bEx := a.rType.MethodByName(funcName)
	return bEx
}

func (a *Actor) Acotr() *Actor {
	return a
}

func (a *Actor) Init() {
	a.mailChan = make(chan bool, 1)
	a.mailBox = mpsc.New[*CallIO]()
	a.actorChan = make(chan int, 1)
	a.timerMap = make(map[int64]*timerUnit)
	a.cronMap = make(map[int64]*cronUnit)
	//trance
	a.trace.Init()
	if a.id == 0 {
		a.id = AssignActorId()
	}
	a.timerId = new(int64)
}

func (a *Actor) register(ac IActor, op Op) {
	rType := reflect.TypeOf(ac)
	a.ActorBase = ActorBase{rType: rType, rVal: reflect.ValueOf(ac), Self: ac, actorName: op.name, actorType: op.actorType}
}

func (a *Actor) RegisterTimer(duration time.Duration, fun func(), opts ...timer.OpOption) int64 {
	actorName := a.actorName
	node, op := timer.RegisterTimer(duration, func(Id int64) {
		a.SendMsg(rpc.RpcHead{ActorName: actorName}, "UpdateTimer", Id)
	}, opts...)
	a.timerMap[node.Id] = &timerUnit{Op: op, TimerNode: node, fun: fun}
	return node.Id
}

func (a *Actor) StopTimer(timerId int64) {
	timerUnit, bEx := a.timerMap[timerId]
	if bEx {
		timerUnit.Stop()
		delete(a.timerMap, timerId)
	}
}

func (a *Actor) RegisterCron(cronStr string, fun func()) int64 {
	sched, err := cron.ParseStandard(cronStr)
	if err != nil {
		base.LOG.Fatalf("RegisterCron [%s] cronStr err", a.actorName)
	} else {
		now := time.Now()
		nextTime := sched.Next(now)
		if !nextTime.IsZero() {
			id := AssignActorId()
			timerId := a.RegisterTimer(nextTime.Sub(now)+timer.TICK_INTERVAL, func() { a.updateCron(id) }, timer.WithCount(1))
			a.cronMap[id] = &cronUnit{Schedule: sched, nextTime: nextTime, timerId: timerId, fun: fun}
			return id
		}
	}
	return 0
}

func (a *Actor) updateCron(id int64) {
	cronUnit, bEx := a.cronMap[id]
	if bEx {
		now := time.Now()
		nextTime := cronUnit.Schedule.Next(cronUnit.nextTime)
		if !nextTime.IsZero() {
			cronUnit.nextTime = nextTime
			if nextTime.Sub(now) >= 0 {
				cronUnit.timerId = a.RegisterTimer(nextTime.Sub(now)+timer.TICK_INTERVAL, func() { a.updateCron(id) }, timer.WithCount(1))
			} else {
				cronUnit.timerId = a.RegisterTimer(timer.TICK_INTERVAL, func() { a.updateCron(id) }, timer.WithCount(1))
			}
			(cronUnit.fun)()
		} else {
			a.StopCron(id)
		}
	}
}

func (a *Actor) StopCron(id int64) {
	cronUnit, bEx := a.cronMap[id]
	if bEx {
		delete(a.cronMap, id)
		a.StopTimer(cronUnit.timerId)
	}
}

func (a *Actor) clear() {
	a.id = 0
	a.setState(ASF_NULL)
	//close(a.acotrChan)
	//close(a.mailChan)
}

func (a *Actor) Stop() {
	for _, v := range a.timerMap {
		v.Stop()
	}
	if atomic.CompareAndSwapInt32(&a.state, ASF_RUN, ASF_STOP) {
		a.actorChan <- DESDORY_EVENT
	}
}

func (a *Actor) Start() {
	if atomic.CompareAndSwapInt32(&a.state, ASF_NULL, ASF_RUN) {
		go a.run()
	}
}

func (a *Actor) SendMsg(head rpc.RpcHead, funcName string, params ...interface{}) {
	head.SocketId = 0
	a.Send(head, rpc.Marshal(&head, &funcName, params...))
}

func (a *Actor) Send(head rpc.RpcHead, packet rpc.Packet) {
	defer func() {
		if err := recover(); err != nil {
			base.TraceCode(err)
		}
	}()

	var io CallIO
	io.RpcHead = head
	io.RpcPacket = packet.RpcPacket
	io.Buff = packet.Buff
	a.mailBox.Push(&io)
	if atomic.LoadInt64(&a.mailIn[0]) == 0 && atomic.CompareAndSwapInt64(&a.mailIn[0], 0, 1) {
		a.mailChan <- true
	}
}

func (a *Actor) Trace(funcName string) {
	a.trace.funcName = funcName
}

func (a *Actor) call(io *CallIO) {
	rpcPacket := io.RpcPacket
	head := io.RpcHead
	funcName := rpcPacket.FuncName
	m, bEx := a.rType.MethodByName(funcName)
	if !bEx {
		log.Printf("func [%s] has no method", funcName)
		return
	}
	rpcPacket.RpcHead.SocketId = io.SocketId
	params := rpc.UnmarshalBody(rpcPacket, m.Type)
	if len(params) >= 1 {
		in := make([]reflect.Value, len(params))
		in[0] = a.rVal
		for i, param := range params {
			if i == 0 {
				continue
			}
			in[i] = reflect.ValueOf(param)
		}

		a.Trace(funcName)
		ret := m.Func.Call(in)
		a.Trace("")
		if ret != nil && head.Reply != "" {
			rets := make([]interface{}, len(ret)+1)
			rets[0] = &head
			for i, param := range ret {
				rets[i+1] = param.Interface()
			}
			rpc.MGR.Call(rets...)
		}
	} else {
		log.Printf("func [%s] params at least one context", funcName)
		//f.Call([]reflect.Value{reflect.ValueOf(ctx)})
	}
}

func (a *Actor) UpdateTimer(ctx context.Context, timerId int64) {
	timerUnit, bEx := a.timerMap[timerId]
	if bEx {
		if timerUnit.IsCount {
			timerUnit.Count--
		}
		a.Trace("timer")
		(timerUnit.fun)()
		a.Trace("")
		if timerUnit.IsCount && timerUnit.Count <= 0 {
			a.StopTimer(timerId)
		}
	}
}

func (a *Actor) consume() {
	atomic.StoreInt64(&a.mailIn[0], 0)
	for data := a.mailBox.Pop(); data != nil; data = a.mailBox.Pop() {
		a.call(data)
	}
}

func (a *Actor) loop() bool {
	defer func() {
		if err := recover(); err != nil {
			base.TraceCode(a.trace.ToString(), err)
		}
	}()

	select {
	case <-a.mailChan:
		a.consume()
	case msg := <-a.actorChan:
		if msg == DESDORY_EVENT {
			return true
		}
	}
	return false
}

func (a *Actor) run() {
	for {
		if a.loop() {
			break
		}
	}

	a.clear()
}

func (a *traceInfo) Init() {
	_, file, _, bOk := runtime.Caller(2)
	if bOk {
		index := strings.LastIndex(file, "/")
		if index != -1 {
			a.fileName = file[index+1:]
			a.filePath = file[:index]
			index1 := strings.LastIndex(a.fileName, ".")
			if index1 != -1 {
				a.className = strings.ToLower(a.fileName[:index1])
			}
		}
	}
}

func (a *traceInfo) ToString() string {
	return fmt.Sprintf("trace go file[%s] call[%s]\n", a.fileName, a.funcName)
}

func GetRpcMethodMap(rType reflect.Type, tagName string) map[string]string {
	rpcMethod := map[string]string{}
	sf, bEx := rType.Elem().FieldByName(tagName)
	if !bEx {
		return rpcMethod
	}
	tag := sf.Tag.Get("rpc")
	methodNames := strings.Split(tag, ";")
	for _, methodName := range methodNames {
		funcName := strings.ToLower(methodName)
		rpcMethod[funcName] = methodName
	}

	return rpcMethod
}
