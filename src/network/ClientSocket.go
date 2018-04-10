package network

import (
	"fmt"
	"io"
	"log"
	"net"
	"base"
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

	this.Connect()
	//延迟，监听关闭
	//defer ln.Close()

	this.m_Conn.SetNoDelay(true)
	go clientRoutine(this)
	return true
}

func (this *ClientSocket) Stop() bool {
	if this.m_bShuttingDown {
		return true
	}

	this.Send([]byte("exit"))//通知服务器
	this.m_bShuttingDown = true
	return true
}

func (this *ClientSocket) SendMsg(funcName string, params  ...interface{}){
	buff := base.GetPacket(funcName, params...)
	buff = base.SetTcpEnd(buff)
	this.Send(buff)
}

func (this *ClientSocket) Send(buff []byte) int {
	n, err := this.m_Conn.Write(buff)
	handleError(err)
	if n > 0 {
		return n
	}
	return 0
}

func (this *ClientSocket) Restart() bool {
	return true
}

func (this *ClientSocket) Connect() bool {
	if this.m_nState == SSF_CONNECT{
		return false
	}
	this.m_nState = SSF_CONNECT
	var strRemote = fmt.Sprintf("%s:%d", this.m_sIP, this.m_nPort)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", strRemote)
	if err != nil {
		log.Printf("%v", err)
	}
	ln, err := net.DialTCP("tcp4", nil, tcpAddr)
	if err != nil {
		return false
	}

	this.m_Conn = ln
	fmt.Printf("连接成功，请输入信息！\n")
	this.CallPacket("COMMON_RegisterRequest")
	return true
}

func (this *ClientSocket) OnDisconnect() {
}
func (this *ClientSocket) OnNetFail(int) {
    this.Stop()
	this.CallPacket("DISCONNECT", this.m_ClientId)
}

func clientRoutine(pClient *ClientSocket) bool {
	if pClient.m_Conn == nil {
		return false
	}

	for {
		if pClient.m_bShuttingDown {
			break
		}

		var buff = make([]byte, pClient.m_MaxReceiveBufferSize)
		n, err := pClient.m_Conn.Read(buff)
		if err == io.EOF {
			fmt.Printf("远程链接：%s已经关闭！\n", pClient.m_Conn.RemoteAddr().String())
			pClient.OnNetFail(0)
			break
		}

        if err != nil {
            handleError(err)
            pClient.OnNetFail(0)
            break;
        }
		if string(buff[:n]) == "exit" {
			fmt.Printf("远程链接：%s退出！\n", pClient.m_Conn.RemoteAddr().String())
			pClient.OnNetFail(0)
			break
		}
		if n > 0 {
			pClient.ReceivePacket(pClient.m_ClientId, buff[:n])
		}
	}

	pClient.Close()
	return true
}
