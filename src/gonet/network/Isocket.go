package network

import (
	"gonet/base"
	"bytes"
	"fmt"
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
	MAX_WRITE_CHAN = 32
)

type (
	HandleFunc func(int,[]byte) bool//回调函数
	Socket struct {
		m_Conn                 net.Conn
		//m_Reader				*bufio.Reader
		//m_Writer				*bufio.Writer
		m_nPort                int
		m_sIP                  string
		m_nState			   int
		m_nConnectType		   int
		m_MaxReceiveBufferSize int
		m_MaxSendBufferSize    int

		m_ClientId int
		m_Seq      int64

		m_TotalNum     int
		m_AcceptedNum  int
		m_ConnectedNum int

		m_SendTimes     int
		m_ReceiveTimes  int
		m_bShuttingDown bool
		m_PacketFuncList	*base.Vector//call back

		m_bHalf		bool
		m_nHalfSize int
		m_pInBuffer []byte
	}

	ISocket interface {
		Init(string, int) bool
		Start() bool
		Stop() bool
		Restart() bool
		Connect() bool
		Disconnect(bool) bool
		OnNetFail(int)
		Clear()
		Close()
		Send([]byte) int
		SendById(int, []byte) int
		SendMsg(string, ...interface{})
		SendMsgById(int,string, ...interface{})
		CallMsg(string, ...interface{})//回调消息处理

		GetState() int
		SetMaxSendBufferSize(int)
		GetMaxSendBufferSize()int
		SetMaxReceiveBufferSize(int)
		GetMaxReceiveBufferSize()int
		BindPacketFunc(HandleFunc)
		SetConnectType(int)
		SetTcpConn(net.Conn)
		ReceivePacket(int,	[]byte)
		HandlePacket(int,	[]byte)
	}
)

// virtual
func (this *Socket) Init(string, int) bool {
	this.m_PacketFuncList = base.NewVector()
	this.m_nState = SSF_SHUT_DOWN
	this.m_MaxReceiveBufferSize = 1024
	this.m_MaxSendBufferSize =1024
	this.m_nConnectType = SERVER_CONNECT
	this.m_bHalf = false
	this.m_nHalfSize = 0
	return true
}

func (this *Socket) Start() bool {
	return true
}
func (this *Socket) Stop() bool {
	this.m_bShuttingDown = true
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

func (this *Socket) GetState() int{
	return  this.m_nState
}

func (this *Socket) Send([]byte) int{
	return  0
}

func (this *Socket) SendById(int, []byte) int{
	return 0
}

func (this *Socket) SendMsg(funcName string, params  ...interface{}){
}

func (this *Socket) SendMsgById(int,string, ...interface{}){
}

func (this *Socket) Clear() {
	this.m_nState = SSF_SHUT_DOWN
	//this.m_nConnectType = CLIENT_CONNECT
	this.m_Conn = nil
	//this.m_Reader = nil
	//this.m_Writer = nil
	this.m_MaxSendBufferSize = 1024
	this.m_MaxReceiveBufferSize = 1024
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

func (this *Socket) GetMaxReceiveBufferSize() int{
	return  this.m_MaxReceiveBufferSize
}

func (this *Socket) SetMaxReceiveBufferSize(maxReceiveSize int){
	this.m_MaxReceiveBufferSize = maxReceiveSize
}

func (this *Socket) GetMaxSendBufferSize() int{
	return  this.m_MaxSendBufferSize
}

func (this *Socket) SetMaxSendBufferSize(maxSendSize int){
	this.m_MaxSendBufferSize = maxSendSize
}

func (this *Socket) SetConnectType(nType int){
	this.m_nConnectType = nType
}

func (this *Socket) SetTcpConn(conn net.Conn){
	this.m_Conn = conn
	//this.m_Reader = bufio.NewReader(conn)
	//this.m_Writer = bufio.NewWriter(conn)
}

func (this *Socket) BindPacketFunc(callfunc HandleFunc){
	this.m_PacketFuncList.Push_back(callfunc)
}

func (this *Socket) CallMsg(funcName string, params ...interface{}){
	buff := base.GetPacket(funcName, params...)
	this.HandlePacket(this.m_ClientId, buff)
}

func (this *Socket) HandlePacket(Id int, buff []byte){
	for _,v := range this.m_PacketFuncList.Array() {
		if (v.(HandleFunc)(Id, buff)){
			break
		}
	}
}

func (this *Socket) ReceivePacket(Id int, dat []byte) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("ReceivePacket", err) // 接受包错误
		}
	}()
	//找包结束
	seekToTcpEnd := func(buff []byte) (bool, int) {
		nLen := bytes.Index(buff, []byte(base.TCP_END))
		if nLen != -1{
			return true, nLen+base.TCP_END_LENGTH
		}
		return false, 0
	}

	buff := append(this.m_pInBuffer, dat...)
	this.m_pInBuffer = []byte{}
	nCurSize := 0
	//fmt.Println(this.m_pInBuffer)
ParsePacekt:
	nPacketSize := 0
	nBufferSize := len(buff[nCurSize:])
	bFindFlag := false
	bFindFlag, nPacketSize = seekToTcpEnd(buff[nCurSize:])
	//fmt.Println(bFindFlag, nPacketSize, nBufferSize)
	if bFindFlag {
		if nBufferSize == nPacketSize { //完整包
			this.HandlePacket(Id, buff[nCurSize:nCurSize+nPacketSize-base.TCP_END_LENGTH])
			nCurSize += nPacketSize
		} else if (nBufferSize > nPacketSize) {
			this.HandlePacket(Id, buff[nCurSize:nCurSize+nPacketSize-base.TCP_END_LENGTH])
			nCurSize += nPacketSize
			goto ParsePacekt
		}
	} else if nBufferSize < base.MAX_PACKET {
		this.m_pInBuffer = buff[nCurSize:]
	} else {
		fmt.Println("超出最大包限制，丢弃该包")
	}
}

//tcp粘包固定包头
/*func (this *Socket) ReceivePacket(Id int, dat []byte){
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("ReceivePacket", err) // 接受包错误
		}
	}()
	//找包结束
	seekToTcpEnd := func(buff []byte) (bool, int){
		nLen := len(buff)
		if nLen < base.PACKET_HEAD_SIZE{
			return false, 0
		}

		nSize := base.BytesToInt(buff[0:4])
		if nSize + base.PACKET_HEAD_SIZE <= nLen{
			return true, nSize+base.PACKET_HEAD_SIZE
		}
		return false, 0
	}

	buff := append(this.m_pInBuffer, dat...)
	this.m_pInBuffer = []byte{}
	nCurSize := 0
	//fmt.Println(this.m_pInBuffer)
ParsePacekt:
	nPacketSize := 0
	nBufferSize := len(buff[nCurSize:])
	bFindFlag := false
	bFindFlag, nPacketSize = seekToTcpEnd(buff[nCurSize:])
	//fmt.Println(bFindFlag, nPacketSize, nBufferSize)
	if bFindFlag{
		if nBufferSize == nPacketSize{		//完整包
			this.HandlePacket(Id, buff[nCurSize+base.PACKET_HEAD_SIZE:nCurSize+nPacketSize])
			nCurSize += nPacketSize
		}else if ( nBufferSize > nPacketSize){
			this.HandlePacket(Id, buff[nCurSize+base.PACKET_HEAD_SIZE:nCurSize+nPacketSize])
			nCurSize += nPacketSize
			goto ParsePacekt
		}
	}else if nBufferSize < base.MAX_PACKET{
		this.m_pInBuffer = buff[nCurSize:]
	}else{
		fmt.Println("超出最大包限制，丢弃该包")
	}
}*/