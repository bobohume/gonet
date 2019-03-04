package network

import (
	"fmt"
	"io"
	"log"
	"net"
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
	m_WriteChan   chan []byte
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

	this.m_WriteChan = make(chan []byte, MAX_WRITE_CHAN)
	this.m_nState = SSF_CONNECT
	this.m_Conn.(*net.TCPConn).SetNoDelay(true)
	//this.m_Conn.SetKeepAlive(true)
	//this.m_Conn.SetKeepAlivePeriod(5*time.Second)
	go serverclientRoutine(this)
	//go serverclientWriteRoutine(this)
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
	if this.m_PacketFuncList.Len() > 0 {
		this.CallMsg("DISCONNECT", this.m_ClientId)
	}else if (this.m_pServer != nil && this.m_pServer.m_PacketFuncList.Len() > 0){
		this.m_pServer.CallMsg("DISCONNECT", this.m_ClientId)
	}
}

func (this *ServerSocketClient) Close() {
	close(this.m_WriteChan)
	this.Socket.Close()
	if this.m_pServer != nil {
		this.m_pServer.DelClinet(this)
	}
}

//防止消息过快
func (this *ServerSocketClient) SendNoBlock(buff []byte) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("WriteBuf", err)
		}
	}()

	if this.m_nConnectType == CLIENT_CONNECT{
		select {
		case this.m_WriteChan <- buff: //chan满后再写即阻塞，select进入default分支报错
		default:
			break
		}
	}else{
		this.m_WriteChan <- buff
	}
}

func serverclientRoutine(pClient *ServerSocketClient) bool {
	if pClient.m_Conn == nil {
		return false
	}

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("serverclientRoutine", err)
		}
	}()

	for {
		if pClient.m_bShuttingDown {
			break
		}

		var buff = make([]byte, pClient.m_MaxReceiveBufferSize)
		//n, err := io.ReadFull(pClient.m_Reader, buff)
		n, err := pClient.m_Conn.Read(buff)
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
		if n > 0 {
			pClient.ReceivePacket(pClient.m_ClientId, buff[:n])
		}
	}

	pClient.Close()
	fmt.Printf("%s关闭连接", pClient.m_sIP)
	return true
}

func serverclientWriteRoutine(pClient *ServerSocketClient) bool {
	for {
		select {
		case buff := <-pClient.m_WriteChan :
			pClient.Send(buff)
		}

		if pClient.m_bShuttingDown {
			break
		}
	}

	pClient.Close()
	return true
}