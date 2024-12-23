package network

import (
	"gonet/base/vector"
	"gonet/rpc"
	"net"
	"sync/atomic"
)

const (
	SSF_NULL = iota
	SSF_RUN  = iota
	SSF_STOP = iota //已经关闭
)

const (
	CLIENT_CONNECT = iota //对外
	SERVER_CONNECT = iota //对内
)

const (
	MAX_SEND_CHAN  = 100
	HEART_TIME_OUT = 30
)

type (
	PacketFunc   func(packet rpc.Packet) bool //回调函数
	HandlePacket func(buff []byte)

	Op struct {
		kcp bool
	}

	OpOption func(*Op)

	Socket struct {
		conn              net.Conn
		port              int
		ip                string
		state             int32
		connectType       int
		receiveBufferSize int //单次接收缓存

		clientId uint32
		seq      int64

		totalNum     int
		acceptedNum  int
		connectedNum int

		sendTimes      int
		receiveTimes   int
		packetFuncList *vector.Vector[PacketFunc] //call back

		isHalf       bool
		halfSize     int
		packetParser PacketParser
		isKcp        bool
	}

	ISocket interface {
		Init(string, int, ...OpOption) bool
		Start() bool
		Stop() bool
		Run() bool
		Restart() bool
		Connect() bool
		Disconnect(bool) bool
		OnNetFail(int)
		Clear()
		Close()
		SendMsg(rpc.RpcHead, string, ...interface{})
		Send(rpc.RpcHead, rpc.Packet) int
		CallMsg(rpc.RpcHead, string, ...interface{}) //回调消息处理

		GetId() uint32
		GetState() int32
		SetState(int32)
		SetReceiveBufferSize(int)
		GetReceiveBufferSize() int
		SetMaxPacketLen(int)
		GetMaxPacketLen() int
		BindPacketFunc(PacketFunc)
		SetConnectType(int)
		SetConn(net.Conn)
		HandlePacket([]byte)
	}
)

func (op *Op) applyOpts(opts []OpOption) {
	for _, opt := range opts {
		opt(op)
	}
}

func WithKcp() OpOption {
	return func(op *Op) {
		op.kcp = true
	}
}

// virtual
func (this *Socket) Init(ip string, port int, params ...OpOption) bool {
	op := &Op{}
	op.applyOpts(params)
	this.packetFuncList = &vector.Vector[PacketFunc]{}
	this.SetState(SSF_NULL)
	this.receiveBufferSize = 1024
	this.connectType = SERVER_CONNECT
	this.isHalf = false
	this.halfSize = 0
	this.packetParser = NewPacketParser(PacketConfig{Func: this.HandlePacket})
	if op.kcp {
		this.isKcp = true
	}
	return true
}

func (this *Socket) Start() bool {
	return true
}

func (this *Socket) Stop() bool {
	if atomic.CompareAndSwapInt32(&this.state, SSF_RUN, SSF_STOP) {
		if this.conn != nil {
			this.conn.Close()
		}
	}
	return false
}

func (this *Socket) Run() bool {
	return true
}

func (this *Socket) Restart() bool {
	return true
}

func (this *Socket) Connect() bool {
	return true
}

func (this *Socket) Disconnect(bool) bool {
	return true
}

func (this *Socket) OnNetFail(int) {
	this.Stop()
}

func (this *Socket) GetId() uint32 {
	return this.clientId
}

func (this *Socket) GetState() int32 {
	return atomic.LoadInt32(&this.state)
}

func (this *Socket) SetState(state int32) {
	atomic.StoreInt32(&this.state, state)
}

func (this *Socket) SendMsg(head rpc.RpcHead, funcName string, params ...interface{}) {
}

func (this *Socket) Send(rpc.RpcHead, rpc.Packet) int {
	return 0
}

func (this *Socket) Clear() {
	this.SetState(SSF_NULL)
	//this.connectType = CLIENT_CONNECT
	this.conn = nil
	this.receiveBufferSize = 1024
	this.isHalf = false
	this.halfSize = 0
}

func (this *Socket) Close() {
	this.Clear()
}

func (this *Socket) GetMaxPacketLen() int {
	return this.packetParser.maxPacketLen
}

func (this *Socket) SetMaxPacketLen(maxReceiveSize int) {
	this.packetParser.maxPacketLen = maxReceiveSize
}

func (this *Socket) GetReceiveBufferSize() int {
	return this.receiveBufferSize
}

func (this *Socket) SetReceiveBufferSize(maxSendSize int) {
	this.receiveBufferSize = maxSendSize
}

func (this *Socket) SetConnectType(nType int) {
	this.connectType = nType
}

func (this *Socket) SetConn(conn net.Conn) {
	this.conn = conn
}

func (this *Socket) BindPacketFunc(callfunc PacketFunc) {
	this.packetFuncList.PushBack(callfunc)
}

func (this *Socket) CallMsg(head rpc.RpcHead, funcName string, params ...interface{}) {
	this.HandlePacket(rpc.Marshal(&head, &funcName, params...).Buff)
}

func (this *Socket) HandlePacket(buff []byte) {
	packet := rpc.Packet{Id: this.clientId, Buff: buff}
	for i := 0; i < this.packetFuncList.Len(); i++ {
		if this.packetFuncList.Get(i)(packet) {
			break
		}
	}
}
