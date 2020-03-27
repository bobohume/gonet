package network

import (
	"gonet/base"
	"fmt"
	"gonet/rpc"
	"io"
	"log"
	"net"
)

type IClientSocket interface {
	ISocket
}

type ClientSocket struct {
	Socket
	m_nMaxClients int
	m_nMinClients int
}

func (this *ClientSocket) Init(ip string, port int) bool {
	if this.m_nPort == port || this.m_sIP == ip {
		return false
	}

	this.Socket.Init(ip, port)
	this.m_sIP = ip
	this.m_nPort = port
	fmt.Println(ip, port)
	return true
}
func (this *ClientSocket) Start() bool {
	this.m_bShuttingDown = false

	if this.m_sIP == "" {
		this.m_sIP = "127.0.0.1"
	}

	if this.Connect(){
		this.m_Conn.(*net.TCPConn).SetNoDelay(true)
		go this.Run()
	}
	//延迟，监听关闭
	//defer ln.Close()
	return true
}

func (this *ClientSocket) Stop() bool {
	if this.m_bShuttingDown {
		return true
	}

	this.m_bShuttingDown = true
	this.Close()
	return true
}

func (this *ClientSocket) SendMsg(funcName string, params  ...interface{}){
	buff := rpc.Marshal(funcName, params...)
	buff = base.SetTcpEnd(buff)
	this.Send(buff)
}

func (this *ClientSocket) Send(buff []byte) int {
	defer func() {
		if err := recover(); err != nil{
			base.TraceCode(err)
		}
	}()

	if this.m_Conn == nil{
		return 0
	}

	n, err := this.m_Conn.Write(buff)
	handleError(err)
	if n > 0 {
		return n
	}
	//this.m_Writer.Flush()
	return 0
}

func (this *ClientSocket) Restart() bool {
	return true
}

func (this *ClientSocket) Connect() bool {
	if this.m_nState == SSF_CONNECT{
		return false
	}

	var strRemote = fmt.Sprintf("%s:%d", this.m_sIP, this.m_nPort)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", strRemote)
	if err != nil {
		log.Printf("%v", err)
	}
	ln, err1 := net.DialTCP("tcp4", nil, tcpAddr)
	if err1 != nil {
		return false
	}

	this.m_nState = SSF_CONNECT
	this.SetTcpConn(ln)
	fmt.Printf("连接成功，请输入信息！\n")
	this.CallMsg("COMMON_RegisterRequest")
	return true
}

func (this *ClientSocket) OnDisconnect() {
}

func (this *ClientSocket) OnNetFail(int) {
    this.Stop()
	this.CallMsg("DISCONNECT", this.m_ClientId)
}

func (this *ClientSocket) Run() bool {
	if this.m_Conn == nil {
		return false
	}

	defer func() {
		if err := recover(); err != nil {
			base.TraceCode(err)
		}
	}()

	var buff= make([]byte, this.m_ReceiveBufferSize)
	for {
		if this.m_bShuttingDown {
			break
		}

		n, err := this.m_Conn.Read(buff)
		if err == io.EOF {
			fmt.Printf("远程链接：%s已经关闭！\n", this.m_Conn.RemoteAddr().String())
			this.OnNetFail(0)
			break
		}

		if err != nil {
			handleError(err)
			this.OnNetFail(0)
			break
		}
		if n > 0 {
			this.ReceivePacket(this.m_ClientId, buff[:n])
		}
	}

	this.Close()
	return true
}