package actor

import (
	"base"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"
	"sync/atomic"
	"message"
)

var(
	g_IdSeed int32
)
type (
	Actor struct {
		m_CallChan  chan CallIO//rpc chan
		m_AcotrChan chan int//use for states
		m_Id       	 int
		m_CallMap	 map[string]interface{}//rpc
		//m_pActorMgr  *ActorMgr
		m_pTimer 	 *time.Ticker//定时器
		m_TimerCall	func()//定时器触发函数
		m_bStart	bool
	}

	IActor interface {
		Init(int)
		Clear()
		Stop()
		Start()
		FindCall(string) interface{}
		RegisterCall(string, interface{})//这里在回调的第一个参数为默认附加参数，为CALLER 信息， 同线程为ACOTRID,remote为SOCKETID
		SendMsg(int,string, ...interface{})
		Send(CallIO)
		PacketFunc(int,[]byte) bool//回调函数
		RegisterTimer(time.Duration, interface{})//注册定时器,时间为纳秒 1000 * 1000 * 1000
		GetId() int
	}

	CallIO struct {
		SocketId int
		ActorId int
		Buff []byte
	}

	Caller struct{
		SocketId int
		ActorId int
	}
)

const (
	DESDORY_EVENT = iota
)

func SendActor(pActor IActor, io CallIO){
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("SendActor", err)
		}
	}()

	if pActor != nil{
		pActor.Send(io)
	}
}

func SendMsg(pActor *Actor, sokcetId int, funcName string, params ...interface{}){
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("SendMsg", err)
		}
	}()

	var io CallIO
	io.ActorId = pActor.m_Id
	io.SocketId = sokcetId
	io.Buff = base.GetPacket(funcName, params...)

	if pActor != nil{
		pActor.Send(io)
	}
}

func AssignActorId() int {
	atomic.AddInt32(&g_IdSeed, 1)
	return int(g_IdSeed)
}

func (this *Actor) GetId() int {
	return this.m_Id
}

func (this *Actor) Init(chanNum int) {
	this.m_CallChan = make(chan CallIO, chanNum)
	this.m_AcotrChan = make(chan int, 1)
	this.m_Id = AssignActorId()
	this.m_CallMap = make(map[string]interface{})
	//this.m_pActorMgr = nil
	this.m_pTimer = time.NewTicker(1<<63-1)//默认没有定时器
	this.m_TimerCall = nil
}

func (this *Actor)  RegisterTimer(duration time.Duration, fun interface{}){
	this.m_pTimer = time.NewTicker(duration)
	this.m_TimerCall = fun.(func())
}

func (this *Actor) Clear() {
	this.m_Id = 0
	this.m_bStart = false
	//this.m_pActorMgr = nil
	close(this.m_AcotrChan)
	close(this.m_CallChan)
	if this.m_pTimer != nil{
		this.m_pTimer.Stop()
	}

	for i := range this.m_CallMap{
		delete(this.m_CallMap, i)
	}
}

func (this *Actor) Stop() {
	this.m_AcotrChan <- DESDORY_EVENT
}

func (this *Actor) Start(){
	if this.m_bStart == false{
		go ActorRoutine(this)
		this.m_bStart = true
	}
}

func (this *Actor) FindCall(funcName string) interface{} {
	funcName = strings.ToLower(funcName)
	fun, exist := this.m_CallMap[funcName]
	if exist == true {
		return fun
	}
	return nil
}

func (this *Actor) RegisterCall(funcName string, call interface{}) {
	switch call.(type) {
	case func(*IActor, []byte):
		log.Fatalln("actor error [%s] 消息定义函数不符合", funcName)
	}
	funcName = strings.ToLower(funcName)
	if this.FindCall(funcName) != nil {
		log.Fatalln("actor error [%s] 消息重复定义", funcName)
	}

	this.m_CallMap[funcName] = call
}

func (this *Actor)  SendMsg(sokcetId int, funcName string, params ...interface{}) {
	var io CallIO
	io.ActorId = this.m_Id
	io.SocketId = sokcetId
	io.Buff = base.GetPacket(funcName, params...)
	this.Send(io)
}

func (this *Actor) Send(io CallIO) {
	this.m_CallChan <- io
}

func (this *Actor) PacketFunc(id int, buff []byte) bool{
	var io CallIO
	io.Buff = buff
	io.SocketId = id

	bitstream := base.NewBitStream(io.Buff, len(io.Buff))
	funcName := bitstream.ReadString()
	funcName = strings.ToLower(funcName)
	pFunc := this.FindCall(funcName)
	if pFunc != nil{
		this.Send(io)
		return true
	}

	return false
}

func (this *Actor) call(io CallIO) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("actor call", err)
		}
	}()
	bitstream := base.NewBitStream(io.Buff, len(io.Buff))
	funcName := bitstream.ReadString()
	funcName = strings.ToLower(funcName)
	pFunc := this.FindCall(funcName)
	if pFunc != nil {
		f := reflect.ValueOf(pFunc)
		k := reflect.TypeOf(pFunc)
		strParams := reflect.TypeOf(pFunc).String()

		nCurLen := bitstream.ReadInt(8)
		params := make([]interface{}, nCurLen+1)
		var caller Caller
		caller.SocketId = io.SocketId
		caller.ActorId = io.ActorId
		params[0] = &caller
		for i := 1; i < nCurLen+1; i++  {
			switch bitstream.ReadInt(8) {
			case 1:
				params[i] = bitstream.ReadFlag()
			case 2:
				params[i] = bitstream.ReadFloat64()
			case 3:
				params[i] = bitstream.ReadFloat()
			case 4:
				params[i] = int8(bitstream.ReadInt(8))
			case 5:
				params[i] = uint8(bitstream.ReadInt(8))
			case 6:
				params[i] = int16(bitstream.ReadInt(16))
			case 7:
				params[i] = uint16(bitstream.ReadInt(16))
			case 8:
				params[i] = int32(bitstream.ReadInt(32))
			case 9:
				params[i] = uint32(bitstream.ReadInt(32))
			case 10:
				params[i] = int64(bitstream.ReadInt64(64))
			case 11:
				params[i] = uint64(bitstream.ReadInt64(64))
			case 12:
				params[i] = bitstream.ReadString()
			case 13:
				params[i] = bitstream.ReadInt(32)
			case 14:
				params[i] = uint(bitstream.ReadInt(32))
			case 15:
				nLen := bitstream.ReadInt(16)
				val := make([]bool, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = bitstream.ReadFlag()
				}
				params[i] = val
			case 16:
				nLen := bitstream.ReadInt(16)
				val := make([]float64, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = bitstream.ReadFloat64()
				}
				params[i] = val
			case 17:
				nLen := bitstream.ReadInt(16)
				val := make([]float32, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = bitstream.ReadFloat()
				}
				params[i] = val
			case 18:
				nLen := bitstream.ReadInt(16)
				val := make([]int8, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = int8(bitstream.ReadInt(8))
				}
				params[i] = val
			case 19:
				nLen := bitstream.ReadInt(16)
				val := make([]uint8, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = uint8(bitstream.ReadInt(8))
				}
				params[i] = val
			case 20:
				nLen := bitstream.ReadInt(16)
				val := make([]int16, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = int16(bitstream.ReadInt(16))
				}
				params[i] = val
			case 21:
				nLen := bitstream.ReadInt(16)
				val := make([]uint16, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = uint16(bitstream.ReadInt(16))
				}
				params[i] = val
			case 22:
				nLen := bitstream.ReadInt(16)
				val := make([]int32, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = int32(bitstream.ReadInt(32))
				}
				params[i] = val
			case 23:
				nLen := bitstream.ReadInt(16)
				val := make([]uint32, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = uint32(bitstream.ReadInt(32))
				}
				params[i] = val
			case 24:
				nLen := bitstream.ReadInt(16)
				val := make([]int64, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = int64(bitstream.ReadInt64(64))
				}
				params[i] = val
			case 25:
				nLen := bitstream.ReadInt(16)
				val := make([]uint64, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = uint64(bitstream.ReadInt64(64))
				}
				params[i] = val
			case 26:
				nLen := bitstream.ReadInt(16)
				val := make([]string, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = bitstream.ReadString()
				}
				params[i] = val
			case 27:
				nLen := bitstream.ReadInt(16)
				val := make([]int, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = bitstream.ReadInt(32)
				}
				params[i] = val
			case 28:
				nLen := bitstream.ReadInt(16)
				val := make([]uint, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = uint(bitstream.ReadInt(32))
				}
				params[i] = val
			case 29://*struct
				packet := base.GetMessage(bitstream.ReadString())
				packet.ReadData(bitstream)
				params[i] = packet
			case 30://[]*struct
				if k.In(i).Kind() != reflect.Slice{
					log.Printf("func [%s] params no fit, func params [%s], params [%v]", funcName, strParams, params)
					return
				}
				nLen := bitstream.ReadInt(16)
				val := reflect.MakeSlice(k.In(i), nLen, nLen)
				for i := 0; i < nLen; i++ {
					packet:= base.GetMessage(bitstream.ReadString())
					packet.ReadData(bitstream)
					val.Index(i).Set(reflect.ValueOf(packet))
				}
				params[i] = val.Interface()
			case 31://protobuf
				packet := message.GetPakcetByName(funcName)
				message.UnmarshalText(packet, bitstream.ReadString())
				params[i] = packet
			default:
				panic("func [%s] params type not supported")
			}
		}

		if k.NumIn()  != len(params) {
			log.Printf("func [%s] can not call, func params [%s], params [%v]", funcName, strParams, params)
			return
		}

		//params no fit
		for i := 0;  i< nCurLen+1; i++ {
			if k.In(i).Kind() != reflect.TypeOf(params[i]).Kind() {
				log.Println(k.In(i).Kind(), reflect.TypeOf(params[i]).Kind())
				log.Printf("func [%s] params no fit, func params [%s], params [%v]", funcName, strParams, params)
				return
			}
		}

		//fmt.Printf("func [%s]",funcName)

		if len(params) >= 1{
			in := make([]reflect.Value, len(params))
			for k, param := range params {
				in[k] = reflect.ValueOf(param)
			}
			f.Call(in)
		}else{
			f.Call(nil);
		}
	}
}

func ActorRoutine(pActor *Actor) bool {
	if pActor == nil {
		return false
	}

	bExit := false
	for {
		select {
		case io := <-pActor.m_CallChan:
			pActor.call(io)
		case msg := <-pActor.m_AcotrChan :
			if msg == DESDORY_EVENT{
				bExit = true
				break
			}
		case <- pActor.m_pTimer.C:
			if pActor.m_TimerCall != nil{
				pActor.m_TimerCall()
			}
		}
		if bExit{
			break
		}
	}

	pActor.Clear()
	return true
}
