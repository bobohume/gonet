package network

import (
	"fmt"
	"io"
	"log"
)

const (
	IDLE_TIMEOUT    = iota
	CONNECT_TIMEOUT = iota
	CONNECT_TYPE    = iota
)

type IServerSocketClient interface {
	ISocket
}

type ServerSocketClient struct {
	Socket
	m_pServer     *ServerSocket
}

func handleError(err error) {
	if err == nil {
		return
	}
	log.Printf("错误：%s\n", err.Error())
}

func (this *ServerSocketClient) Start() bool {
	if this.m_nState != SSF_SHUT_DOWN{
		return false
	}

	if this.m_pServer == nil {
		return false
	}

	this.m_nState = SSF_CONNECT
	this.m_Conn.SetNoDelay(true)
	go serverclientRoutine(this)
	return true
}

func (this *ServerSocketClient) Send(buff []byte) int {
	defer func() {
		if err := recover(); err != nil{
			fmt.Println("ServerSocketClient Send", err)
		}
	}()

	if len(buff) > this.m_MaxSendBufferSize{
		log.Print(" SendError size",len(buff))
		return  0
	}

	n, err := this.m_Conn.Write(buff)
	handleError(err)
	if n > 0 {
		return n
	}
	//this.m_Writer.Flush()
	return 0
}

func (this *ServerSocketClient) ReceivePacket(Id int, buff []byte){
	if this.m_PacketFuncList.Len() > 0 {
		this.Socket.ReceivePacket(this.m_ClientId, buff)
	}else if (this.m_pServer != nil && this.m_pServer.m_PacketFuncList.Len() > 0){
		this.m_pServer.Socket.ReceivePacket(this.m_ClientId, buff)
	}
}

func (this *ServerSocketClient) OnNetFail(error int) {
	this.Stop()
	this.CallPacket("DISCONNECT", this.m_ClientId)
}

func (this *ServerSocketClient) Close() {
	this.Socket.Close()
	if this.m_pServer != nil {
		this.m_pServer.DelClinet(this)
	}
}

func serverclientRoutine(pClient *ServerSocketClient) bool {
	if pClient.m_Conn == nil {
		return false
	}

	for {
		if pClient.m_bShuttingDown {
			break
		}

		var buff = make([]byte, pClient.m_MaxReceiveBufferSize)
		n, err := io.ReadFull(pClient.m_Reader, buff)
		//n, err := pClient.m_Conn.Read(buff)
		if err == io.EOF {
			fmt.Printf("远程链接：%s已经关闭！\n", pClient.m_Conn.RemoteAddr().String())
			pClient.OnNetFail(0)
			break
		}
		if err != nil {
			handleError(err)
			pClient.OnNetFail(0)
			break
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
	fmt.Printf("%s关闭连接", pClient.m_sIP)
	return true
}