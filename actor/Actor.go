package actor

import (
	"context"
	"fmt"
	"gonet/base"
	"gonet/base/mpsc"
	"gonet/common/timer"
	"gonet/rpc"
	"log"
	"reflect"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
	"unsafe"
)

var (
	g_IdSeed int64
)

const (
	ASF_NULL = iota
	ASF_RUN  = iota
	ASF_STOP = iota //已经关闭
)

//********************************************************
// actor 核心actor模式
//********************************************************
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
		acotrChan chan int //use for states
		id        int64
		state     int32
		trace     traceInfo //trace func
		mailBox   *mpsc.Queue[*CallIO]
		mailIn    [8]int64
		mailChan  chan bool
		timerId   *int64
		pool      IActorPool //ACTOR_TYPE_VIRTUAL,ACTOR_TYPE_POOL
	}

	IActor interface {
		Init()
		Stop()
		Start()
		SendMsg(head rpc.RpcHead, funcName string, params ...interface{})
		Send(head rpc.RpcHead, packet rpc.Packet)
		RegisterTimer(duration time.Duration, fun func(), opts ...timer.OpOption) //注册定时器,时间为纳秒 1000 * 1000 * 1000
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
	a.mailChan = make(chan bool)
	a.mailBox = mpsc.New[*CallIO]()
	a.acotrChan = make(chan int, 1)
	//trance
	a.trace.Init()
	if a.id == 0 {
		a.id = AssignActorId()
	}
}

func (a *Actor) register(ac IActor, op Op) {
	rType := reflect.TypeOf(ac)
	a.ActorBase = ActorBase{rType: rType, rVal: reflect.ValueOf(ac), Self: ac, actorName: op.name, actorType: op.actorType}
}

func (a *Actor) RegisterTimer(duration time.Duration, fun func(), opts ...timer.OpOption) {
	if a.timerId == nil {
		a.timerId = new(int64)
		*a.timerId = a.id
	}

	timer.RegisterTimer(a.timerId, duration, func() {
		a.SendMsg(rpc.RpcHead{ActorName: a.actorName}, "UpdateTimer", (*int64)(unsafe.Pointer(&fun)))
	}, opts...)
}

func (a *Actor) clear() {
	a.id = 0
	a.setState(ASF_NULL)
	//close(a.acotrChan)
	//close(a.mailChan)
	timer.StopTimer(a.timerId)
}

func (a *Actor) Stop() {
	if atomic.CompareAndSwapInt32(&a.state, ASF_RUN, ASF_STOP) {
		a.acotrChan <- DESDORY_EVENT
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
	if !bEx{
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
			ret = append([]reflect.Value{reflect.ValueOf(&head)}, ret...)
			rpc.MGR.Call(ret)
		}
	} else {
		log.Printf("func [%s] params at least one context", funcName)
		//f.Call([]reflect.Value{reflect.ValueOf(ctx)})
	}
}

func (a *Actor) UpdateTimer(ctx context.Context, p *int64) {
	func1 := (*func())(unsafe.Pointer(p))
	a.Trace("timer")
	(*func1)()
	a.Trace("")
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
	case msg := <-a.acotrChan:
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
