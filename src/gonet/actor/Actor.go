package actor

import (
	"context"
	"gonet/base"
	"gonet/rpc"
	"log"
	"reflect"
	"strings"
	"sync/atomic"
	"time"
)

var(
	g_IdSeed int64
)

//********************************************************
// actor 核心actor模式
//********************************************************
type (
	Actor struct {
		m_CallChan  chan CallIO//rpc chan
		m_AcotrChan chan int//use for states
		m_Id       	 int64
		m_CallMap	 map[string] *CallFunc//rpc
		m_pTimer 	 *time.Ticker//定时器
		m_TimerCall	func()//定时器触发函数
		m_bStart	bool
	}

	IActor interface {
		Init(chanNum int)
		Stop()
		Start()
		FindCall(funcName string) *CallFunc
		RegisterCall(funcName string, call interface{})
		SendMsg(head rpc.RpcHead, funcName string, params ...interface{})
		Send(head rpc.RpcHead, buff []byte)
		PacketFunc(id uint32, buff []byte) bool//回调函数
		RegisterTimer(duration time.Duration, fun interface{})//注册定时器,时间为纳秒 1000 * 1000 * 1000
		GetId() int64
		GetRpcHead(ctx context.Context) rpc.RpcHead//rpc is safe
	}

	CallIO struct {
		rpc.RpcHead
		Buff []byte
	}

	CallFunc struct {
		Func interface{}
		FuncType reflect.Type
		FuncVal reflect.Value
		FuncParams string
	}
)

const (
	DESDORY_EVENT = iota
)

func AssignActorId() int64 {
	atomic.AddInt64(&g_IdSeed, 1)
	return int64(g_IdSeed)
}

func (this *Actor) GetId() int64 {
	return this.m_Id
}

func (this *Actor) GetRpcHead(ctx context.Context) rpc.RpcHead{
	rpcHead := ctx.Value("rpcHead").(rpc.RpcHead)
	return rpcHead
}

func (this *Actor) Init(chanNum int) {
	this.m_CallChan = make(chan CallIO, chanNum)
	this.m_AcotrChan = make(chan int, 1)
	this.m_Id = AssignActorId()
	this.m_CallMap = make(map[string] *CallFunc)
	this.m_pTimer = time.NewTicker(1<<63-1)//默认没有定时器
	this.m_TimerCall = nil
}

func (this *Actor)  RegisterTimer(duration time.Duration, fun interface{}){
	this.m_pTimer.Stop()
	this.m_pTimer = time.NewTicker(duration)
	this.m_TimerCall = fun.(func())
}

func (this *Actor) clear() {
	this.m_Id = 0
	this.m_bStart = false
	close(this.m_AcotrChan)
	close(this.m_CallChan)
	if this.m_pTimer != nil{
		this.m_pTimer.Stop()
	}

	this.m_CallMap = make(map[string] *CallFunc)
}

func (this *Actor) Stop() {
	this.m_AcotrChan <- DESDORY_EVENT
}

func (this *Actor) Start(){
	if this.m_bStart == false{
		go this.run()
		this.m_bStart = true
	}
}

func (this *Actor) FindCall(funcName string) *CallFunc{
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

	this.m_CallMap[funcName] = &CallFunc{Func:call, FuncVal:reflect.ValueOf(call), FuncType:reflect.TypeOf(call), FuncParams:reflect.TypeOf(call).String()}
}

func (this *Actor) SendMsg(head rpc.RpcHead,funcName string, params ...interface{}) {
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
	this.m_CallChan <- io
}

func (this *Actor) PacketFunc(id uint32, buff []byte) bool{
	rpcPacket, head := rpc.UnmarshalHead(buff)
	if this.FindCall(rpcPacket.FuncName) != nil{
		head.SocketId = id
		this.Send(head, buff)
		return true
	}

	return false
}

func (this *Actor) call(io CallIO) {
	rpcPacket, _ := rpc.UnmarshalHead(io.Buff)
	funcName := rpcPacket.FuncName
	pFunc := this.FindCall(funcName)

	if pFunc != nil {
		f := pFunc.FuncVal
		k := pFunc.FuncType
		strParams := pFunc.FuncParams
		rpcPacket.RpcHead.SocketId = io.SocketId
		params := rpc.UnmarshalBody(rpcPacket, k)

		if k.NumIn()  != len(params){
			log.Printf("func [%s] can not call, func params [%s], params [%v]", funcName, strParams, params)
			return
		}

		if len(params) >= 1{
			in := make([]reflect.Value, len(params))
			for i, param := range params {
				in[i] = reflect.ValueOf(param)
			}

			f.Call(in)
		}else{
			log.Printf("func [%s] params at least one context", funcName)
			//f.Call([]reflect.Value{reflect.ValueOf(ctx)})
		}
	}
}

func (this *Actor) loop() bool{
	defer func() {
		if err := recover(); err != nil{
			base.TraceCode(err)
		}
	}()

	select {
	case io := <-this.m_CallChan:
		this.call(io)
	case msg := <-this.m_AcotrChan :
		if msg == DESDORY_EVENT{
			return false
		}
	case <- this.m_pTimer.C:
		if this.m_TimerCall != nil{
			this.m_TimerCall()
		}
	}
	return true
}

func (this *Actor) run(){
	for {
		if !this.loop(){
			break
		}
	}

	this.clear()
}