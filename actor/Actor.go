package actor

import (
	"context"
	"fmt"
	"gonet/base"
	"gonet/base/mpsc"
	"gonet/rpc"
	"log"
	"reflect"
	"runtime"
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
		m_AcotrChan chan int//use for states
		m_Id       	 int64
		m_CallMap	 map[string] *CallFunc//rpc
		m_pTimer 	 *time.Ticker//定时器
		m_TimerCall	func()//定时器触发函数
		m_bStart	bool
		m_Trace  	traceInfo//trace func
		m_MailBox	*mpsc.Queue
		m_bMailIn	int32
		m_MailChan  chan bool
	}

	IActor interface {
		Init()
		Stop()
		Start()
		FindCall(funcName string) *CallFunc
		RegisterCall(funcName string, call interface{})
		SendMsg(head rpc.RpcHead, funcName string, params ...interface{})
		Send(head rpc.RpcHead, buff []byte)
		PacketFunc(packet rpc.Packet) bool//回调函数
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

	traceInfo struct {
		funcName string
		fileName string
		filePath string
		className string
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

func (this *Actor) Init() {
	this.m_MailChan = make(chan bool)
	this.m_MailBox = mpsc.New()
	this.m_AcotrChan = make(chan int, 1)
	this.m_Id = AssignActorId()
	this.m_CallMap = make(map[string] *CallFunc)
	this.m_pTimer = time.NewTicker(1<<63-1)//默认没有定时器
	this.m_TimerCall = nil
	//trance
	this.m_Trace.Init()
}

func (this *Actor)  RegisterTimer(duration time.Duration, fun interface{}){
	this.m_pTimer.Stop()
	this.m_pTimer = time.NewTicker(duration)
	this.m_TimerCall = fun.(func())
}

func (this *Actor) clear() {
	this.m_Id = 0
	this.m_bStart = false
	//close(this.m_AcotrChan)
	//close(this.m_MailChan)
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

	callfunc := &CallFunc{Func:call, FuncVal:reflect.ValueOf(call), FuncType:reflect.TypeOf(call), FuncParams:reflect.TypeOf(call).String()}
	this.m_CallMap[funcName] = callfunc
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
	this.m_MailBox.Push(io)
	if atomic.CompareAndSwapInt32(&this.m_bMailIn, 0, 1){
		this.m_MailChan <- true
	}
}

func (this *Actor) PacketFunc(packet rpc.Packet) bool{
	rpcPacket, head := rpc.UnmarshalHead(packet.Buff)
	if this.FindCall(rpcPacket.FuncName) != nil{
		head.SocketId = packet.Id
		head.Reply = packet.Reply
		this.Send(head, packet.Buff)
		return true
	}

	return false
}

func (this *Actor) Trace(funcName string){
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

		if len(params) >= 1{
			in := make([]reflect.Value, len(params))
			for i, param := range params {
				in[i] = reflect.ValueOf(param)
			}

			this.Trace(funcName)
			ret := f.Call(in)
			this.Trace("")
			if ret != nil && head.Reply != ""{
				ret = append([]reflect.Value{reflect.ValueOf(&head)}, ret...)
				rpc.GCall.Call(ret)
			}
		}else{
			log.Printf("func [%s] params at least one context", funcName)
			//f.Call([]reflect.Value{reflect.ValueOf(ctx)})
		}
	}
}

func (this *Actor) consume() {
	for data :=  this.m_MailBox.Pop(); data != nil; data =  this.m_MailBox.Pop(){
		this.call(data.(CallIO))
	}
	atomic.StoreInt32(&this.m_bMailIn, 0)
}

func (this *Actor) loop() bool{
	defer func() {
		if err := recover(); err != nil{
			base.TraceCode(this.m_Trace.ToString(), err)
		}
	}()

	select {
	case  <-this.m_MailChan:
		this.consume()
	case msg := <-this.m_AcotrChan :
		if msg == DESDORY_EVENT{
			return false
		}
	case <- this.m_pTimer.C:
		if this.m_TimerCall != nil{
			this.Trace("timer")
			this.m_TimerCall()
			this.Trace("")
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

func (this *traceInfo) Init() {
	_, file, _,  bOk := runtime.Caller(2)
	if bOk{
		index := strings.LastIndex(file, "/")
		if index!= -1{
			this.fileName = file[index+1:]
			this.filePath = file[:index]
			index1 := strings.LastIndex(this.fileName, ".")
			if index1!= -1 {
				this.className = strings.ToLower(this.fileName[:index1])
			}
		}
	}
}

func (this *traceInfo) ToString() string{
	return fmt.Sprintf("trace go file[%s] call[%s]\n", this.fileName, this.funcName)
}