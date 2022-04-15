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
		m_kcp bool
	}

	OpOption func(*Op)

	Socket struct {
		m_Conn              net.Conn
		m_nPort             int
		m_sIP               string
		m_State             int32
		m_nConnectType      int
		m_ReceiveBufferSize int //单次接收缓存

		m_ClientId uint32
		m_Seq      int64

		m_TotalNum     int
		m_AcceptedNum  int
		m_ConnectedNum int

		m_SendTimes      int
		m_ReceiveTimes   int
		m_PacketFuncList *vector.Vector //call back

		m_bHalf        bool
		m_nHalfSize    int
		m_PacketParser PacketParser
		m_HeartTime    int
		m_bKcp         bool
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
		op.m_kcp = true
	}
}

// virtual
func (this *Socket) Init(ip string, port int, params ...OpOption) bool {
	op := &Op{}
	op.applyOpts(params)
	this.m_PacketFuncList = vector.NewVector()
	this.SetState(SSF_NULL)
	this.m_ReceiveBufferSize = 1024
	this.m_nConnectType = SERVER_CONNECT
	this.m_bHalf = false
	this.m_nHalfSize = 0
	this.m_HeartTime = 0
	this.m_PacketParser = NewPacketParser(PacketConfig{Func: this.HandlePacket})
	if op.m_kcp {
		this.m_bKcp = true
	}
	return true
}

func (this *Socket) Start() bool {
	return true
}

func (this *Socket) Stop() bool {
	if this.m_Conn != nil && atomic.CompareAndSwapInt32(&this.m_State, SSF_RUN, SSF_STOP) {
		this.m_Conn.Close()
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
	return this.m_ClientId
}

func (this *Socket) GetState() int32 {
	return atomic.LoadInt32(&this.m_State)
}

func (this *Socket) SetState(state int32) {
	atomic.StoreInt32(&this.m_State, state)
}

func (this *Socket) SendMsg(head rpc.RpcHead, funcName string, params ...interface{}) {
}

func (this *Socket) Send(rpc.RpcHead, rpc.Packet) int {
	return 0
}

func (this *Socket) Clear() {
	this.SetState(SSF_NULL)
	//this.m_nConnectType = CLIENT_CONNECT
	this.m_Conn = nil
	this.m_ReceiveBufferSize = 1024
	this.m_bHalf = false
	this.m_nHalfSize = 0
	this.m_HeartTime = 0
}

func (this *Socket) Close() {
	this.Clear()
}

func (this *Socket) GetMaxPacketLen() int {
	return this.m_PacketParser.m_MaxPacketLen
}

func (this *Socket) SetMaxPacketLen(maxReceiveSize int) {
	this.m_PacketParser.m_MaxPacketLen = maxReceiveSize
}

func (this *Socket) GetReceiveBufferSize() int {
	return this.m_ReceiveBufferSize
}

func (this *Socket) SetReceiveBufferSize(maxSendSize int) {
	this.m_ReceiveBufferSize = maxSendSize
}

func (this *Socket) SetConnectType(nType int) {
	this.m_nConnectType = nType
}

func (this *Socket) SetConn(conn net.Conn) {
	this.m_Conn = conn
}

func (this *Socket) BindPacketFunc(callfunc PacketFunc) {
	this.m_PacketFuncList.PushBack(callfunc)
}

func (this *Socket) CallMsg(head rpc.RpcHead, funcName string, params ...interface{}) {
	this.HandlePacket(rpc.Marshal(&head, &funcName, params...).Buff)
}

func (this *Socket) HandlePacket(buff []byte) {
	packet := rpc.Packet{Id: this.m_ClientId, Buff: buff}
	for _, v := range this.m_PacketFuncList.Values() {
		if v.(PacketFunc)(packet) {
			break
		}
	}
}
