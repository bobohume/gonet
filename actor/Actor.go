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
		m_ActorName  string
		m_RType      reflect.Type
		m_RVal       reflect.Value
		m_ActorType  ACTOR_TYPE
		Self     	 IActor //when parent interface class call interface, it call parent not child  use for virtual
		m_RpcMethodMap  map[string] string
	}

	Actor struct {
		ActorBase
		m_AcotrChan chan int //use for states
		m_Id        int64
		m_State     int32
		m_Trace     traceInfo //trace func
		m_MailBox   *mpsc.Queue
		m_bMailIn   [8]int64
		m_MailChan  chan bool
		m_TimerId   *int64
		share_rpc	int `rpc:"GetRpcHead;UpdateTimer"`
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
		setState(state int32)
		GetRpcHead(ctx context.Context) rpc.RpcHead //rpc is safe
		GetName() string
		GetActorType() ACTOR_TYPE
		Register(IActor,  Op)
		HasRpc(string) bool
		GetAcotr() *Actor
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

func (this *ActorBase) IsActorType(actorType ACTOR_TYPE) bool {
	return this.m_ActorType == actorType
}

func AssignActorId() int64 {
	return atomic.AddInt64(&g_IdSeed, 1)
}

func (this *Actor) GetId() int64 {
	return this.m_Id
}

func (this *Actor) SetId(id int64)  {
	this.m_Id = id
}

func (this *Actor) GetName() string {
	return this.m_ActorName
}

func (this *Actor) GetRpcHead(ctx context.Context) rpc.RpcHead {
	rpcHead := ctx.Value("rpcHead").(rpc.RpcHead)
	return rpcHead
}

func (this *Actor) GetState() int32 {
	return atomic.LoadInt32(&this.m_State)
}

func (this *Actor) GetActorType() ACTOR_TYPE{
	return this.m_ActorType
}

func (this *Actor) setState(state int32) {
	atomic.StoreInt32(&this.m_State, state)
}

func (this *Actor) HasRpc(funcName string) bool{
	_, bEx := this.m_RpcMethodMap[funcName]
	return bEx
}

func (this *Actor) GetAcotr() *Actor{
	return this
}

func (this *Actor) Init() {
	this.m_MailChan = make(chan bool)
	this.m_MailBox = mpsc.New()
	this.m_AcotrChan = make(chan int, 1)
	//trance
	this.m_Trace.Init()
	if this.m_Id == 0 {
		this.m_Id = AssignActorId()
	}
}

func (this *Actor) Register(pActor IActor, op Op){
	rType := reflect.TypeOf(pActor)
	this.ActorBase = ActorBase{m_RType: rType, m_RVal: reflect.ValueOf(pActor), Self:pActor, m_ActorName:op.m_name, m_ActorType:op.m_type, m_RpcMethodMap: op.m_RpcMethodMap}
}

func (this *Actor) RegisterTimer(duration time.Duration, fun func(), opts ...timer.OpOption) {
	if this.m_TimerId == nil {
		this.m_TimerId = new(int64)
		*this.m_TimerId = this.m_Id
	}

	timer.RegisterTimer(this.m_TimerId, duration, func() {
		this.SendMsg(rpc.RpcHead{ActorName:this.m_ActorName}, "UpdateTimer", (*int64)(unsafe.Pointer(&fun)))
	}, opts...)
}

func (this *Actor) clear() {
	this.m_Id = 0
	this.setState(ASF_NULL)
	//close(this.m_AcotrChan)
	//close(this.m_MailChan)
	timer.StopTimer(this.m_TimerId)
}

func (this *Actor) Stop() {
	if atomic.CompareAndSwapInt32(&this.m_State, ASF_RUN, ASF_STOP) {
		this.m_AcotrChan <- DESDORY_EVENT
	}
}

func (this *Actor) Start() {
	if atomic.CompareAndSwapInt32(&this.m_State, ASF_NULL, ASF_RUN) {
		go this.run()
	}
}

func (this *Actor) SendMsg(head rpc.RpcHead, funcName string, params ...interface{}) {
	head.SocketId = 0
	this.Send(head, rpc.Marshal(head, funcName, params...))
}

func (this *Actor) Send(head rpc.RpcHead, packet rpc.Packet) {
	defer func() {
		if err := recover(); err != nil {
			base.TraceCode(err)
		}
	}()

	var io CallIO
	io.RpcHead = head
	io.RpcPacket = packet.RpcPacket
	io.Buff = packet.Buff
	this.m_MailBox.Push(io)
	if atomic.LoadInt64(&this.m_bMailIn[0]) == 0 && atomic.CompareAndSwapInt64(&this.m_bMailIn[0], 0, 1) {
		this.m_MailChan <- true
	}
}

func (this *Actor) Trace(funcName string) {
	this.m_Trace.funcName = funcName
}

func (this *Actor) call(io CallIO) {
	rpcPacket := io.RpcPacket
	head := io.RpcHead
	funcName := rpcPacket.FuncName
	if !this.HasRpc(funcName) {
		log.Printf("func [%s] has no method", funcName)
		return
	}

	methodName := this.m_RpcMethodMap[funcName]
	m,_ := this.m_RType.MethodByName(methodName)
	rpcPacket.RpcHead.SocketId = io.SocketId
	params := rpc.UnmarshalBody(rpcPacket, m.Type)
	if len(params) >= 1 {
		in := make([]reflect.Value, len(params))
		in[0] = this.m_RVal
		for i, param := range params{
			if i == 0{
				continue
			}
			in[i] = reflect.ValueOf(param)
		}

		this.Trace(funcName)
		ret := m.Func.Call(in)
		this.Trace("")
		if ret != nil && head.Reply != "" {
			ret = append([]reflect.Value{reflect.ValueOf(&head)}, ret...)
			rpc.GCall.Call(ret)
		}
	} else {
		log.Printf("func [%s] params at least one context", funcName)
		//f.Call([]reflect.Value{reflect.ValueOf(ctx)})
	}
}

func (this *Actor) UpdateTimer(ctx context.Context, p *int64) {
	func1 := (*func())(unsafe.Pointer(p))
	this.Trace("timer")
	(*func1)()
	this.Trace("")
}

func (this *Actor) consume() {
	atomic.StoreInt64(&this.m_bMailIn[0], 0)
	for data := this.m_MailBox.Pop(); data != nil; data = this.m_MailBox.Pop() {
		this.call(data.(CallIO))
	}
}

func (this *Actor) loop() bool {
	defer func(){
		if err := recover(); err != nil {
			base.TraceCode(this.m_Trace.ToString(), err)
		}
	}()

	select {
	case <-this.m_MailChan:
		this.consume()
	case msg := <-this.m_AcotrChan:
		if msg == DESDORY_EVENT {
			return true
		}
	}
	return false
}

func (this *Actor) run() {
	for {
		if this.loop() {
			break
		}
	}

	this.clear()
}

func (this *traceInfo) Init() {
	_, file, _, bOk := runtime.Caller(2)
	if bOk {
		index := strings.LastIndex(file, "/")
		if index != -1 {
			this.fileName = file[index+1:]
			this.filePath = file[:index]
			index1 := strings.LastIndex(this.fileName, ".")
			if index1 != -1 {
				this.className = strings.ToLower(this.fileName[:index1])
			}
		}
	}
}

func (this *traceInfo) ToString() string {
	return fmt.Sprintf("trace go file[%s] call[%s]\n", this.fileName, this.funcName)
}

func GetRpcMethodMap(rType reflect.Type, tagName string) map[string] string{
	rpcMethod := map[string] string{}
	sf, bEx := rType.Elem().FieldByName(tagName)
	if !bEx{
		return rpcMethod
	}
	tag := sf.Tag.Get("rpc")
	methodNames := strings.Split(tag, ";")
	for _, methodName := range methodNames{
		funcName := strings.ToLower(methodName)
		rpcMethod[funcName] = methodName
	}

	return rpcMethod
}