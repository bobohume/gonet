package actor

import (
	"base"
	"fmt"
	"log"
	"message"
	"reflect"
	"strings"
	"sync/atomic"
	"time"
	"unsafe"
)

var(
	g_IdSeed int64
)
type (
	Actor struct {
		m_CallChan  chan CallIO//rpc chan
		m_AcotrChan chan int//use for states
		m_Id       	 int64
		m_CallMap	 map[string]interface{}//rpc
		//m_pActorMgr  *ActorMgr
		m_pTimer 	 *time.Ticker//定时器
		m_TimerCall	func()//定时器触发函数
		m_bStart	bool
		m_SocketId	int
		m_CallId	int64
	}

	IActor interface {
		Init(chanNum int)
		Stop()
		Start()
		FindCall(funcName string) interface{}
		RegisterCall(funcName string, call interface{})//这里在回调的第一个参数为默认附加参数，为CALLER 信息， 同线程为ACOTRID,remote为SOCKETID
		SendMsg(funcName string, params ...interface{})
		Send(io CallIO)
		PacketFunc(id int, buff []byte) bool//回调函数
		RegisterTimer(duration time.Duration, fun interface{})//注册定时器,时间为纳秒 1000 * 1000 * 1000
		GetId() int64
		GetCallId() int64
		GetSocketId() int//rpc is safe
		SendNoBlock(io CallIO)
	}

	CallIO struct {
		SocketId int
		ActorId int64
		Buff []byte
	}
)

const (
	DESDORY_EVENT = iota
)

/*func SendMsg(pActor IActor, sokcetId int, funcName string, params ...interface{}){
	var io CallIO
	io.ActorId = pActor.GetId()
	io.SocketId = sokcetId
	io.Buff = base.GetPacket(funcName, params...)

	if pActor != nil{
		pActor.Send(io)
	}
}*/

func AssignActorId() int64 {
	atomic.AddInt64(&g_IdSeed, 1)
	return int64(g_IdSeed)
}

func (this *Actor) GetId() int64 {
	return this.m_Id
}

func (this *Actor) GetSocketId() int {
	return this.m_SocketId
}

func (this *Actor) GetCallId() int64 {
	return this.m_CallId
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
	this.m_pTimer.Stop()
	this.m_pTimer = time.NewTicker(duration)
	this.m_TimerCall = fun.(func())
}

func (this *Actor) clear() {
	this.m_Id = 0
	this.m_CallId = 0
	this.m_SocketId = 0
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
		go this.run()
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

func (this *Actor) SendMsg(funcName string, params ...interface{}) {
	var io CallIO
	io.ActorId = this.m_Id
	io.SocketId = 0
	io.Buff = base.GetPacket(funcName, params...)
	this.Send(io)
}


func (this *Actor) Send(io CallIO) {
	//go func() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Send", err)
		}
	}()

	this.m_CallChan <- io
	//}()
}

//防止消息过快,主要在player里面
func (this *Actor) SendNoBlock(io CallIO) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("SendNoBlock", err)
		}
	}()

	select {
	case this.m_CallChan <- io: //chan满后再写即阻塞，select进入default分支报错
	default:
		break
	}
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
	funcName := ""
	bitstream := base.NewBitStream(io.Buff, len(io.Buff))
	funcName = bitstream.ReadString()
	funcName = strings.ToLower(funcName)
	pFunc := this.FindCall(funcName)
	if pFunc != nil {
		f := reflect.ValueOf(pFunc)
		k := reflect.TypeOf(pFunc)
		strParams := reflect.TypeOf(pFunc).String()

		nCurLen := bitstream.ReadInt(8)
		params := make([]interface{}, nCurLen)
		this.m_SocketId = io.SocketId
		this.m_CallId = io.ActorId
		for i := 0; i < nCurLen; i++  {
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


			case 21:
				nLen := bitstream.ReadInt(16)
				val := make([]bool, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = bitstream.ReadFlag()
				}
				params[i] = val
			case 22:
				nLen := bitstream.ReadInt(16)
				val := make([]float64, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = bitstream.ReadFloat64()
				}
				params[i] = val
			case 23:
				nLen := bitstream.ReadInt(16)
				val := make([]float32, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = bitstream.ReadFloat()
				}
				params[i] = val
			case 24:
				nLen := bitstream.ReadInt(16)
				val := make([]int8, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = int8(bitstream.ReadInt(8))
				}
				params[i] = val
			case 25:
				nLen := bitstream.ReadInt(16)
				val := make([]uint8, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = uint8(bitstream.ReadInt(8))
				}
				params[i] = val
			case 26:
				nLen := bitstream.ReadInt(16)
				val := make([]int16, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = int16(bitstream.ReadInt(16))
				}
				params[i] = val
			case 27:
				nLen := bitstream.ReadInt(16)
				val := make([]uint16, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = uint16(bitstream.ReadInt(16))
				}
				params[i] = val
			case 28:
				nLen := bitstream.ReadInt(16)
				val := make([]int32, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = int32(bitstream.ReadInt(32))
				}
				params[i] = val
			case 29:
				nLen := bitstream.ReadInt(16)
				val := make([]uint32, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = uint32(bitstream.ReadInt(32))
				}
				params[i] = val
			case 30:
				nLen := bitstream.ReadInt(16)
				val := make([]int64, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = int64(bitstream.ReadInt64(64))
				}
				params[i] = val
			case 31:
				nLen := bitstream.ReadInt(16)
				val := make([]uint64, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = uint64(bitstream.ReadInt64(64))
				}
				params[i] = val
			case 32:
				nLen := bitstream.ReadInt(16)
				val := make([]string, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = bitstream.ReadString()
				}
				params[i] = val
			case 33:
				nLen := bitstream.ReadInt(16)
				val := make([]int, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = bitstream.ReadInt(32)
				}
				params[i] = val
			case 34:
				nLen := bitstream.ReadInt(16)
				val := make([]uint, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = uint(bitstream.ReadInt(32))
				}
				params[i] = val
			case 35://[]struct
				if k.In(i).Kind() != reflect.Slice{
					log.Printf("func [%s] params no fit, func params [%s], params [%v]", funcName, strParams, params)
					return
				}
				nLen := bitstream.ReadInt(16)
				val := reflect.MakeSlice(k.In(i), nLen, nLen)
				for i := 0; i < nLen; i++ {
					packet:= base.GetMessage(bitstream.ReadString())
					base.ReadData(packet, bitstream)
					val.Index(i).Set(reflect.ValueOf(packet).Elem())
				}
				params[i] = val.Interface()


			case 41:
				nLen := bitstream.ReadInt(16)
				aa := bool(false)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
				val := reflect.New(tVal).Elem()
				for i := 0; i < nLen; i++ {
					val.Index(i).SetBool(bitstream.ReadFlag())
				}
				params[i] = val.Interface()
			case 42:
				nLen := bitstream.ReadInt(16)
				aa := float64(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
				val := reflect.New(tVal).Elem()
				for i := 0; i < nLen; i++ {
					val.Index(i).SetFloat(bitstream.ReadFloat64())
				}
				params[i] = val.Interface()
			case 43:
				nLen := bitstream.ReadInt(16)
				aa := float32(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
				val := reflect.New(tVal).Elem()
				for i := 0; i < nLen; i++ {
					val.Index(i).SetFloat(float64(bitstream.ReadFloat()))
				}
				params[i] = val.Interface()
			case 44:
				nLen := bitstream.ReadInt(16)
				aa := int8(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
				val := reflect.New(tVal).Elem()
				for i := 0; i < nLen; i++ {
					val.Index(i).SetInt(int64(bitstream.ReadInt(8)))
				}
				params[i] = val.Interface()
			case 45:
				nLen := bitstream.ReadInt(16)
				aa := uint8(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
				val := reflect.New(tVal).Elem()
				for i := 0; i < nLen; i++ {
					val.Index(i).SetInt(int64(bitstream.ReadInt(8)))
				}
				params[i] = val.Interface()
			case 46:
				nLen := bitstream.ReadInt(16)
				aa := int16(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
				val := reflect.New(tVal).Elem()
				for i := 0; i < nLen; i++ {
					val.Index(i).SetInt(int64(bitstream.ReadInt(16)))
				}
				params[i] = val.Interface()
			case 47:
				nLen := bitstream.ReadInt(16)
				aa := uint16(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
				val := reflect.New(tVal).Elem()
				for i := 0; i < nLen; i++ {
					val.Index(i).SetInt(int64(bitstream.ReadInt(16)))
				}
				params[i] = val.Interface()
			case 48:
				nLen := bitstream.ReadInt(16)
				aa := int32(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
				val := reflect.New(tVal).Elem()
				for i := 0; i < nLen; i++ {
					val.Index(i).SetInt(int64(bitstream.ReadInt(32)))
				}
				params[i] = val.Interface()
			case 49:
				nLen := bitstream.ReadInt(16)
				aa := uint32(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
				val := reflect.New(tVal).Elem()
				for i := 0; i < nLen; i++ {
					val.Index(i).SetInt(int64(bitstream.ReadInt(32)))
				}
				params[i] = val.Interface()
			case 50:
				nLen := bitstream.ReadInt(16)
				aa := int64(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
				val := reflect.New(tVal).Elem()
				for i := 0; i < nLen; i++ {
					val.Index(i).SetInt(int64(bitstream.ReadInt64(64)))
				}
				params[i] = val.Interface()
			case 51:
				nLen := bitstream.ReadInt(16)
				aa := uint64(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
				val := reflect.New(tVal).Elem()
				for i := 0; i < nLen; i++ {
					val.Index(i).SetInt(int64(bitstream.ReadInt64(64)))
				}
				params[i] = val.Interface()
			case 52:
				nLen := bitstream.ReadInt(16)
				aa := string("")
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
				val := reflect.New(tVal).Elem()
				for i := 0; i < nLen; i++ {
					val.Index(i).SetString(bitstream.ReadString())
				}
				params[i] = val.Interface()
			case 53:
				nLen := bitstream.ReadInt(16)
				aa := int(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
				val := reflect.New(tVal).Elem()
				for i := 0; i < nLen; i++ {
					val.Index(i).SetInt(int64(bitstream.ReadInt(32)))
				}
				params[i] = val.Interface()
			case 54:
				nLen := bitstream.ReadInt(16)
				aa := uint(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
				val := reflect.New(tVal).Elem()
				for i := 0; i < nLen; i++ {
					val.Index(i).SetInt(int64(bitstream.ReadInt(32)))
				}
				params[i] = val.Interface()
			/*case 55://[*]struct
				if k.In(i).Kind() != reflect.Array{
					log.Printf("func [%s] params no fit, func params [%s], params [%v]", funcName, strParams, params)
					return
				}
				nLen := bitstream.ReadInt(16)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(k.In(i)))
				val := reflect.New(tVal).Elem()
				arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
				for i := 0; i < nLen; i++ {
					value :=  unsafe.Pointer(unsafe.Pointer(arrayPtr))
					packet:= base.GetMessage(bitstream.ReadString())
					base.ReadData(packet, bitstream)
					arrayPtr = arrayPtr + unsafe.Sizeof(packet)
					*(*unsafe.Pointer)(value) = unsafe.Pointer(reflect.ValueOf(packet).Pointer())
				}

				params[i] = val.Interface()*/

			case 61:
				val := new(bool)
				*val = bitstream.ReadFlag()
				params[i] = val
			case 62:
				val := new(float64)
				*val = bitstream.ReadFloat64()
				params[i] = val
			case 63:
				val := new(float32)
				*val = bitstream.ReadFloat()
				params[i] = val
			case 64:
				val := new(int8)
				*val = int8(bitstream.ReadInt(8))
				params[i] = val
			case 65:
				val := new(uint8)
				*val = uint8(bitstream.ReadInt(8))
				params[i] = val
			case 66:
				val := new(int16)
				*val = int16(bitstream.ReadInt(16))
				params[i] = val
			case 67:
				val := new(uint16)
				*val = uint16(bitstream.ReadInt(16))
				params[i] = val
			case 68:
				val := new(int32)
				*val = int32(bitstream.ReadInt(32))
				params[i] = val
			case 69:
				val := new(uint32)
				*val = uint32(bitstream.ReadInt(32))
				params[i] = val
			case 70:
				val := new(int64)
				*val = int64(bitstream.ReadInt64(64))
				params[i] = val
			case 71:
				val := new(uint64)
				*val = uint64(bitstream.ReadInt64(64))
				params[i] = val
			case 72:
				val := new(string)
				*val = bitstream.ReadString()
				params[i] = val
			case 73:
				val := new(int)
				*val = bitstream.ReadInt(32)
				params[i] = val
			case 74:
				val := new(uint)
				*val = uint(bitstream.ReadInt(32))
				params[i] = val
			case 75://*struct
				packet := base.GetMessage(bitstream.ReadString())
				base.ReadData(packet, bitstream)
				params[i] = packet



			case 81:
				nLen := bitstream.ReadInt(16)
				val := make([]*bool, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = new(bool)
					*val[i] = bitstream.ReadFlag()
				}
				params[i] = val
			case 82:
				nLen := bitstream.ReadInt(16)
				val := make([]*float64, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = new(float64)
					*val[i] = bitstream.ReadFloat64()
				}
				params[i] = val
			case 83:
				nLen := bitstream.ReadInt(16)
				val := make([]*float32, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = new(float32)
					*val[i] = bitstream.ReadFloat()
				}
				params[i] = val
			case 84:
				nLen := bitstream.ReadInt(16)
				val := make([]*int8, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = new(int8)
					*val[i] = int8(bitstream.ReadInt(8))
				}
				params[i] = val
			case 85:
				nLen := bitstream.ReadInt(16)
				val := make([]*uint8, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = new(uint8)
					*val[i] = uint8(bitstream.ReadInt(8))
				}
				params[i] = val
			case 86:
				nLen := bitstream.ReadInt(16)
				val := make([]*int16, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = new(int16)
					*val[i] = int16(bitstream.ReadInt(16))
				}
				params[i] = val
			case 87:
				nLen := bitstream.ReadInt(16)
				val := make([]*uint16, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = new(uint16)
					*val[i] = uint16(bitstream.ReadInt(16))
				}
				params[i] = val
			case 88:
				nLen := bitstream.ReadInt(16)
				val := make([]*int32, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = new(int32)
					*val[i] = int32(bitstream.ReadInt(32))
				}
				params[i] = val
			case 89:
				nLen := bitstream.ReadInt(16)
				val := make([]*uint32, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = new(uint32)
					*val[i] = uint32(bitstream.ReadInt(32))
				}
				params[i] = val
			case 90:
				nLen := bitstream.ReadInt(16)
				val := make([]*int64, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = new(int64)
					*val[i] = int64(bitstream.ReadInt64(64))
				}
				params[i] = val
			case 91:
				nLen := bitstream.ReadInt(16)
				val := make([]*uint64, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = new(uint64)
					*val[i] = uint64(bitstream.ReadInt64(64))
				}
				params[i] = val
			case 92:
				nLen := bitstream.ReadInt(16)
				val := make([]*string, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = new(string)
					*val[i] = bitstream.ReadString()
				}
				params[i] = val
			case 93:
				nLen := bitstream.ReadInt(16)
				val := make([]*int, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = new(int)
					*val[i] = bitstream.ReadInt(32)
				}
				params[i] = val
			case 94:
				nLen := bitstream.ReadInt(16)
				val := make([]*uint, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = new(uint)
					*val[i] = uint(bitstream.ReadInt(32))
				}
				params[i] = val
			case 95://[]*struct
				if k.In(i).Kind() != reflect.Slice{
					log.Printf("func [%s] params no fit, func params [%s], params [%v]", funcName, strParams, params)
					return
				}
				nLen := bitstream.ReadInt(16)
				val := reflect.MakeSlice(k.In(i), nLen, nLen)
				for i := 0; i < nLen; i++ {
					packet:= base.GetMessage(bitstream.ReadString())
					base.ReadData(packet, bitstream)
					val.Index(i).Set(reflect.ValueOf(packet))
				}
				params[i] = val.Interface()


			case 101:
				nLen := bitstream.ReadInt(16)
				aa := bool(false)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
				val := reflect.New(tVal).Elem()
				arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
				for i := 0; i < nLen; i++ {
					value :=  (**bool)(unsafe.Pointer(arrayPtr))
					arrayPtr = arrayPtr + base.SIZE_BOOL
					val1 := bitstream.ReadFlag()
					*value = &val1
				}
				params[i] = val.Interface()
			case 102:
				nLen := bitstream.ReadInt(16)
				aa := float64(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
				val := reflect.New(tVal).Elem()
				arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
				for i := 0; i < nLen; i++ {
					value :=  (**float64)(unsafe.Pointer(arrayPtr))
					arrayPtr = arrayPtr + base.SIZE_FLOAT64
					val1 := bitstream.ReadFloat64()
					*value = &val1
				}
				params[i] = val.Interface()
			case 103:
				nLen := bitstream.ReadInt(16)
				aa := float32(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
				val := reflect.New(tVal).Elem()
				arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
				for i := 0; i < nLen; i++ {
					value :=  (**float32)(unsafe.Pointer(arrayPtr))
					arrayPtr = arrayPtr + base.SIZE_FLOAT32
					val1 := float32(bitstream.ReadFloat64())
					*value =  &val1
				}
				params[i] = val.Interface()
			case 104:
				nLen := bitstream.ReadInt(16)
				aa := int8(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
				val := reflect.New(tVal).Elem()
				arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
				for i := 0; i < nLen; i++ {
					value :=  (**int8)(unsafe.Pointer(arrayPtr))
					arrayPtr = arrayPtr + base.SIZE_INT8
					val1 := int8(bitstream.ReadInt(8))
					*value =  &val1
				}
				params[i] = val.Interface()
			case 105:
				nLen := bitstream.ReadInt(16)
				aa := uint8(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
				val := reflect.New(tVal).Elem()
				arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
				for i := 0; i < nLen; i++ {
					value :=  (**uint8)(unsafe.Pointer(arrayPtr))
					arrayPtr = arrayPtr + base.SIZE_UINT8
					val1 := uint8(bitstream.ReadInt(8))
					*value = &val1
				}
				params[i] = val.Interface()
			case 106:
				nLen := bitstream.ReadInt(16)
				aa := int16(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
				val := reflect.New(tVal).Elem()
				arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
				for i := 0; i < nLen; i++ {
					value :=  (**int16)(unsafe.Pointer(arrayPtr))
					arrayPtr = arrayPtr + base.SIZE_INT16
					val1 := int16(bitstream.ReadInt(16))
					*value =&val1
				}
				params[i] = val.Interface()
			case 107:
				nLen := bitstream.ReadInt(16)
				aa := uint16(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
				val := reflect.New(tVal).Elem()
				arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
				for i := 0; i < nLen; i++ {
					value :=  (**uint16)(unsafe.Pointer(arrayPtr))
					arrayPtr = arrayPtr + base.SIZE_UINT16
					val1 := uint16(bitstream.ReadInt(16))
					*value = &val1
				}
				params[i] = val.Interface()
			case 108:
				nLen := bitstream.ReadInt(16)
				aa := int32(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
				val := reflect.New(tVal).Elem()
				arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
				for i := 0; i < nLen; i++ {
					value :=  (**int32)(unsafe.Pointer(arrayPtr))
					arrayPtr = arrayPtr + base.SIZE_INT32
					val1 := int32(bitstream.ReadInt(32))
					*value = &val1
				}
				params[i] = val.Interface()
			case 109:
				nLen := bitstream.ReadInt(16)
				aa := uint32(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
				val := reflect.New(tVal).Elem()
				arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
				for i := 0; i < nLen; i++ {
					value :=  (**uint32)(unsafe.Pointer(arrayPtr))
					arrayPtr = arrayPtr + base.SIZE_UINT32
					val1 := uint32(bitstream.ReadInt(32))
					*value = &val1
				}
				params[i] = val.Interface()
			case 110:
				nLen := bitstream.ReadInt(16)
				aa := int64(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
				val := reflect.New(tVal).Elem()
				arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
				for i := 0; i < nLen; i++ {
					value :=  (**int64)(unsafe.Pointer(arrayPtr))
					arrayPtr = arrayPtr + base.SIZE_INT64
					val1 := int64(bitstream.ReadInt64(64))
					*value =  &val1
				}
				params[i] = val.Interface()
			case 111:
				nLen := bitstream.ReadInt(16)
				aa := uint64(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
				val := reflect.New(tVal).Elem()
				arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
				for i := 0; i < nLen; i++ {
					value :=  (**uint64)(unsafe.Pointer(arrayPtr))
					arrayPtr = arrayPtr + base.SIZE_UINT64
					val1 := uint64(bitstream.ReadInt64(64))
					*value = &val1
				}
				params[i] = val.Interface()
			case 112:
				nLen := bitstream.ReadInt(16)
				aa := string("")
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
				val := reflect.New(tVal).Elem()
				arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
				for i := 0; i < nLen; i++ {
					value :=  (**string)(unsafe.Pointer(arrayPtr))
					arrayPtr = arrayPtr + base.SIZE_STRING
					val1 := string(bitstream.ReadString())
					*value = &val1
				}
				params[i] = val.Interface()
			case 113:
				nLen := bitstream.ReadInt(16)
				aa := int(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
				val := reflect.New(tVal).Elem()
				arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
				for i := 0; i < nLen; i++ {
					value :=  (**int)(unsafe.Pointer(arrayPtr))
					arrayPtr = arrayPtr + base.SIZE_INT
					val1 := bitstream.ReadInt(32)
					*value = &val1
				}
				params[i] = val.Interface()
			case 114:
				nLen := bitstream.ReadInt(16)
				aa := uint(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
				val := reflect.New(tVal).Elem()
				arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
				for i := 0; i < nLen; i++ {
					value :=  (**uint)(unsafe.Pointer(arrayPtr))
					arrayPtr = arrayPtr + base.SIZE_UINT
					val1 := uint(bitstream.ReadInt(32))
					*value = &val1
				}
				params[i] = val.Interface()
			case 115:
				if k.In(i).Kind() != reflect.Array{
					log.Printf("func [%s] params no fit, func params [%s], params [%v]", funcName, strParams, params)
					return
				}
				nLen := bitstream.ReadInt(16)
				val := reflect.New(k.In(i)).Elem()
				arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
				for i := 0; i < nLen; i++ {
					value :=  unsafe.Pointer(unsafe.Pointer(arrayPtr))
					packet:= base.GetMessage(bitstream.ReadString())
					base.ReadData(packet, bitstream)
					arrayPtr = arrayPtr + base.SIZE_PTR
					*(*unsafe.Pointer)(value) = unsafe.Pointer(reflect.ValueOf(packet).Pointer())
				}
				params[i] = val.Interface()


			case 120://protobuf
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
		for i := 0;  i< nCurLen; i++ {
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

func (this *Actor) loop() bool{
	defer func() {
		if err := recover(); err != nil{
			base.TraceCode()
		}
	}()

	select {
	case io := <-this.m_CallChan:
		this.call(io)
	case msg := <-this.m_AcotrChan :
		if msg == DESDORY_EVENT{
			return true
		}
	case <- this.m_pTimer.C:
		if this.m_TimerCall != nil{
			this.m_TimerCall()
		}
	}
	return false
}

func (this *Actor) run(){
	for {
		if this.loop(){
			break
		}
	}

	this.clear()
}
/*func ActorRoutine(pActor *Actor) bool {
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

	pActor.clear()
	return true
}*/