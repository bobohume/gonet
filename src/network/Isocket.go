package network

import (
	"net"
	"base"
	"fmt"
	"bufio"
)

const (
	SSF_ACCEPT    		 = iota
	SSF_CONNECT    	 = iota
	SSF_SHUT_DOWN      = iota //已经关闭
)

const (
	CLIENT_CONNECT = iota
	SERVER_CONNECT = iota
)

type (
	HandleFunc func(int,[]byte) bool//回调函数
	Socket struct {
		m_Conn                 *net.TCPConn
		m_Reader				*bufio.Reader
		m_Writer				*bufio.Writer
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
		SendByID(int, []byte) int
		SendMsg(string, ...interface{})
		SendMsgByID(int,string, ...interface{})

		GetState() int
		SetMaxSendBufferSize(int)
		GetMaxSendBufferSize()int
		SetMaxReceiveBufferSize(int)
		GetMaxReceiveBufferSize()int
		BindPacketFunc(HandleFunc)
		SetConnectType(int)
		SetTcpConn(*net.TCPConn)
		ReceivePacket(int,	[]byte)
		HandlePacket(int,	[]byte)
		CallPacket(string, ...interface{})//回调消息处理
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

func (this *Socket) SendByID(int, []byte) int{
	return 0
}

func (this *Socket) SendMsg(funcName string, params  ...interface{}){
}

func (this *Socket) SendMsgByID(int,string, ...interface{}){
}

func (this *Socket) Clear() {
	this.m_nState = SSF_SHUT_DOWN
	this.m_nConnectType = CLIENT_CONNECT
	this.m_Conn = nil
	this.m_Reader = nil
	this.m_Writer = nil
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

func (this *Socket) SetTcpConn(conn *net.TCPConn){
	this.m_Conn = conn
	this.m_Reader = bufio.NewReader(conn)
	this.m_Writer = bufio.NewWriter(conn)
}

func (this *Socket) BindPacketFunc(callfunc HandleFunc){
	this.m_PacketFuncList.Push_back(callfunc)
}

func (this *Socket) CallPacket(funcName string, params ...interface{}){
	buff := base.GetPacket(funcName, params...)
	buff = base.SetTcpEnd(buff)
	this.HandlePacket(this.m_ClientId, buff)
}

func (this *Socket) HandlePacket(Id int, buff []byte){
	for _,v := range this.m_PacketFuncList.Array() {
		if (v.(HandleFunc)(Id, buff)){
			break
		}
	}
}

func (this *Socket) ReceivePacket(Id int,	buff []byte){
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("ReceivePacket", err) // 接受包错误
		}
	}()

	nPacketSize := 0
	nBufferSize := len(buff)
	bFindFlag := false
	//找包结束
	seekToTcpEnd := func(buff []byte) (bool, int){
		nLen := len(buff)
		for	i := 0; i < nLen - 1; i++ {
			if (buff[i] == base.TCP_END[0] &&  buff[i+1] == base.TCP_END[1]){
				return true, i+2
			}
		}
		return false, 0
	}

	if this.m_bHalf{
		this.m_bHalf = false
		this.m_pInBuffer = append(this.m_pInBuffer, buff...)
		bFindFlag, nPacketSize = seekToTcpEnd(this.m_pInBuffer)
		if bFindFlag {
			//fmt.Println(this.m_nHalfSize)
			if nBufferSize == nPacketSize{		//完整包
				this.HandlePacket(Id, this.m_pInBuffer[:nPacketSize-2])
			}else if ( nBufferSize > nPacketSize){
				this.HandlePacket(Id, this.m_pInBuffer[:nPacketSize-2])
				this.ReceivePacket(Id,	this.m_pInBuffer[nPacketSize:])//继续解析
			}
		}else{//丢弃包
			fmt.Println("丢弃一个不完整的包")
		}
	}else{
		bFindFlag, nPacketSize = seekToTcpEnd(buff)
		//fmt.Println(bFindFlag, nPacketSize, nBufferSize)
		if bFindFlag {
			if nBufferSize == nPacketSize{		//完整包
				this.HandlePacket(Id, buff[:nPacketSize-2])
			}else if ( nBufferSize > nPacketSize){
				//fmt.Println(buff[:nPacketSize-2])
				//fmt.Println(buff[nPacketSize:])
				this.HandlePacket(Id, buff[:nPacketSize-2])
				this.ReceivePacket(Id,	buff[nPacketSize:])//继续解析
			}
		}else{
			this.m_bHalf = true
			this.m_nHalfSize = nBufferSize
			this.m_pInBuffer = buff
		}
	}
}
