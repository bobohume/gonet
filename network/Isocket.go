package network

import (
	"gonet/base/vector"
	"gonet/rpc"
	"net"
)

const (
	SSF_ACCEPT    		 = iota
	SSF_CONNECT    	 = iota
	SSF_SHUT_DOWN      = iota //已经关闭
)

const (
	CLIENT_CONNECT = iota//对外
	SERVER_CONNECT = iota//对内
)

const (
	MAX_SEND_CHAN = 100
)

type (
	PacketFunc func(packet rpc.Packet) bool//回调函数
	HandlePacket func(buff []byte)
	Socket struct {
		m_Conn                 net.Conn
		m_nPort                int
		m_sIP                  string
		m_nState			   int
		m_nConnectType		   int
		m_ReceiveBufferSize    int//单次接收缓存

		m_ClientId uint32
		m_Seq      int64

		m_TotalNum     int
		m_AcceptedNum  int
		m_ConnectedNum int

		m_SendTimes     int
		m_ReceiveTimes  int
		m_bShuttingDown bool
		m_PacketFuncList	*vector.Vector//call back

		m_bHalf		bool
		m_nHalfSize int
		m_PacketParser PacketParser
	}

	ISocket interface {
		Init(string, int) bool
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
		Send(rpc.RpcHead, []byte) int
		CallMsg(string, ...interface{})//回调消息处理

		GetId() uint32
		GetState() int
		SetReceiveBufferSize(int)
		GetReceiveBufferSize()int
		SetMaxPacketLen(int)
		GetMaxPacketLen()int
		BindPacketFunc(PacketFunc)
		SetConnectType(int)
		SetTcpConn(net.Conn)
		HandlePacket([]byte)
	}
)

// virtual
func (this *Socket) Init(string, int) bool {
	this.m_PacketFuncList = vector.NewVector()
	this.m_nState = SSF_SHUT_DOWN
	this.m_ReceiveBufferSize =1024
	this.m_nConnectType = SERVER_CONNECT
	this.m_bHalf = false
	this.m_nHalfSize = 0
	this.m_PacketParser = NewPacketParser(PacketConfig{Func:this.HandlePacket})
	return true
}

func (this *Socket) Start() bool {
	return true
}

func (this *Socket) Stop() bool {
	this.m_bShuttingDown = true
	return true
}

func (this *Socket) Run()bool {
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

func (this *Socket) GetId() uint32{
	return this.m_ClientId
}

func (this *Socket) GetState() int{
	return  this.m_nState
}

func (this *Socket) SendMsg(head rpc.RpcHead, funcName string, params  ...interface{}){
}

func (this *Socket) Send(rpc.RpcHead, []byte) int{
	return  0
}

func (this *Socket) Clear() {
	this.m_nState = SSF_SHUT_DOWN
	//this.m_nConnectType = CLIENT_CONNECT
	this.m_Conn = nil
	this.m_ReceiveBufferSize = 1024
	this.m_bShuttingDown = false
	this.m_bHalf = false
	this.m_nHalfSize = 0
}

func (this *Socket) Close() {
	if this.m_Conn != nil{
		this.m_Conn.Close()
	}
	this.Clear()
}

func (this *Socket) GetMaxPacketLen() int{
	return  this.m_PacketParser.m_MaxPacketLen
}

func (this *Socket) SetMaxPacketLen(maxReceiveSize int){
	this.m_PacketParser.m_MaxPacketLen = maxReceiveSize
}

func (this *Socket) GetReceiveBufferSize() int{
	return  this.m_ReceiveBufferSize
}

func (this *Socket) SetReceiveBufferSize(maxSendSize int){
	this.m_ReceiveBufferSize = maxSendSize
}

func (this *Socket) SetConnectType(nType int){
	this.m_nConnectType = nType
}

func (this *Socket) SetTcpConn(conn net.Conn){
	this.m_Conn = conn
}

func (this *Socket) BindPacketFunc(callfunc PacketFunc){
	this.m_PacketFuncList.PushBack(callfunc)
}

func (this *Socket) CallMsg(funcName string, params ...interface{}){
	buff := rpc.Marshal(rpc.RpcHead{}, funcName, params...)
	this.HandlePacket(buff)
}

func (this *Socket) HandlePacket(buff []byte){
	packet := rpc.Packet{Id:this.m_ClientId, Buff:buff}
	for _,v := range this.m_PacketFuncList.Values() {
		if (v.(PacketFunc)(packet)){
			break
		}
	}
}