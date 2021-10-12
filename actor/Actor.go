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
	Actor struct {
		m_AcotrChan chan int //use for states
		m_Id        int64
		m_CallMap   map[string]*CallFunc //rpc
		m_State     int32
		m_Trace     traceInfo //trace func
		m_MailBox   *mpsc.Queue
		m_bMailIn   int32
		m_MailChan  chan bool
		m_TimerId   *int64
	}

	IActor interface {
		Init()
		Stop()
		Start()
		FindCall(funcName string) *CallFunc
		RegisterCall(funcName string, call interface{})
		SendMsg(head rpc.RpcHead, funcName string, params ...interface{})
		Send(head rpc.RpcHead, buff []byte)
		PacketFunc(packet rpc.Packet) bool                                        //回调函数
		RegisterTimer(duration time.Duration, fun func(), opts ...timer.OpOption) //注册定时器,时间为纳秒 1000 * 1000 * 1000
		GetId() int64
		GetState() int32
		setState(state int32)
		GetRpcHead(ctx context.Context) rpc.RpcHead //rpc is safe
	}

	CallIO struct {
		rpc.RpcHead
		Buff []byte
	}

	CallFunc struct {
		Func       interface{}
		FuncType   reflect.Type
		FuncVal    reflect.Value
		FuncParams string
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

func AssignActorId() int64 {
	return atomic.AddInt64(&g_IdSeed, 1)
}

func (this *Actor) GetId() int64 {
	return this.m_Id
}

func (this *Actor) GetRpcHead(ctx context.Context) rpc.RpcHead {
	rpcHead := ctx.Value("rpcHead").(rpc.RpcHead)
	return rpcHead
}

func (this *Actor) GetState() int32 {
	return atomic.LoadInt32(&this.m_State)
}

func (this *Actor) setState(state int32) {
	atomic.StoreInt32(&this.m_State, state)
}

func (this *Actor) Init() {
	this.m_MailChan = make(chan bool)
	this.m_MailBox = mpsc.New()
	this.m_AcotrChan = make(chan int, 1)
	this.m_Id = AssignActorId()
	this.m_CallMap = make(map[string]*CallFunc)
	//trance
	this.RegisterCall("UpdateTimer", func(ctx context.Context, p *int64) {
		func1 := (*func())(unsafe.Pointer(p))
		this.Trace("timer")
		(*func1)()
		this.Trace("")
	})
	this.m_Trace.Init()
}

func (this *Actor) RegisterTimer(duration time.Duration, fun func(), opts ...timer.OpOption) {
	if this.m_TimerId == nil {
		this.m_TimerId = new(int64)
		*this.m_TimerId = this.m_Id
	}
	timer.RegisterTimer(this.m_TimerId, duration, func() {
		this.SendMsg(rpc.RpcHead{}, "UpdateTimer", (*int64)(unsafe.Pointer(&fun)))
	}, opts...)
}

func (this *Actor) clear() {
	this.m_Id = 0
	this.setState(ASF_NULL)
	//close(this.m_AcotrChan)
	//close(this.m_MailChan)
	timer.StopTimer(this.m_TimerId)
	this.m_CallMap = make(map[string]*CallFunc)
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

func (this *Actor) FindCall(funcName string) *CallFunc {
	funcName = strings.ToLower(funcName)
	fun, exist := this.m_CallMap[funcName]
	if exist == true {
		return fun
	}
	return nil
}

func (this *Actor) RegisterCall(funcName string, call interface{}) {
	funcName = strings.ToLower(funcName)
	if this.FindCall(funcName) != nil {
		log.Fatalln("actor error [%s] 消息重复定义", funcName)
	}

	callfunc := &CallFunc{Func: call, FuncVal: reflect.ValueOf(call), FuncType: reflect.TypeOf(call), FuncParams: reflect.TypeOf(call).String()}
	this.m_CallMap[funcName] = callfunc
}

func (this *Actor) SendMsg(head rpc.RpcHead, funcName string, params ...interface{}) {
	head.SocketId = 0
	this.Send(head, rpc.Marshal(head, funcName, params...))
}

func (this *Actor) Send(head rpc.RpcHead, buff []byte) {
	defer func() {
		if err := recover(); err != nil {
			base.TraceCode(err)
		}
	}()

	var io CallIO
	io.RpcHead = head
	io.Buff = buff
	this.m_MailBox.Push(io)
	if atomic.CompareAndSwapInt32(&this.m_bMailIn, 0, 1) {
		this.m_MailChan <- true
	}
}

func (this *Actor) PacketFunc(packet rpc.Packet) bool {
	rpcPacket, head := rpc.UnmarshalHead(packet.Buff)
	if this.FindCall(rpcPacket.FuncName) != nil {
		head.SocketId = packet.Id
		head.Reply = packet.Reply
		this.Send(head, packet.Buff)
		return true
	}

	return false
}

func (this *Actor) Trace(funcName string) {
	this.m_Trace.funcName = funcName
}

func (this *Actor) call(io CallIO) {
	rpcPacket, _ := rpc.Unmarshal(io.Buff)
	head := io.RpcHead
	funcName := rpcPacket.FuncName
	pFunc := this.FindCall(funcName)
	if pFunc != nil {
		f := pFunc.FuncVal
		k := pFunc.FuncType
		rpcPacket.RpcHead.SocketId = io.SocketId
		params := rpc.UnmarshalBody(rpcPacket, k)

		if len(params) >= 1 {
			in := make([]reflect.Value, len(params))
			for i, param := range params {
				in[i] = reflect.ValueOf(param)
			}

			this.Trace(funcName)
			ret := f.Call(in)
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
}

func (this *Actor) consume() {
	atomic.StoreInt32(&this.m_bMailIn, 0)
	for data := this.m_MailBox.Pop(); data != nil; data = this.m_MailBox.Pop() {
		this.call(data.(CallIO))
	}
}

func (this *Actor) loop() bool {
	defer func() {
		if err := recover(); err != nil {
			base.TraceCode(this.m_Trace.ToString(), err)
		}
	}()

	select {
	case <-this.m_MailChan:
		this.consume()
	case msg := <-this.m_AcotrChan:
		if msg == DESDORY_EVENT {
			return false
		}
	}
	return true
}

func (this *Actor) run() {
	for {
		if !this.loop() {
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
