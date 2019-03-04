package network

import (
	"fmt"
	"io"
	"log"
)

type IWebSocketClient interface {
	ISocket
}

type WebSocketClient struct {
	Socket
	m_pServer     *WebSocket
	m_WriteChan   chan []byte
}

func (this *WebSocketClient) Start() bool {
	if this.m_nState != SSF_SHUT_DOWN{
		return false
	}

	if this.m_pServer == nil {
		return false
	}

	this.m_WriteChan = make(chan []byte, MAX_WRITE_CHAN)
	this.m_nState = SSF_ACCEPT
	go wserverclientRoutine(this)
	//go wserverclientWriteRoutine(this)
	return true
}

func (this *WebSocketClient) Send(buff []byte) int {
	defer func() {
		if err := recover(); err != nil{
			fmt.Println("WebSocketClient Send", err)
		}
	}()

	if this.m_Conn == nil{
		return 0
	}

	if len(buff) > this.m_MaxSendBufferSize{
		log.Print(" SendError size",len(buff))
		return  0
	}

	n, err := this.m_Conn.Write(buff)
	handleError(err)
	if n > 0 {
		return n
	}
	return 0
}

func (this *WebSocketClient) ReceivePacket(Id int, buff []byte){
	if this.m_PacketFuncList.Len() > 0 {
		this.Socket.ReceivePacket(this.m_ClientId, buff)
	}else if (this.m_pServer != nil && this.m_pServer.m_PacketFuncList.Len() > 0){
		this.m_pServer.Socket.ReceivePacket(this.m_ClientId, buff)
	}
}

func (this *WebSocketClient) OnNetFail(error int) {
	this.Stop()
	if this.m_PacketFuncList.Len() > 0 {
		this.CallMsg("DISCONNECT", this.m_ClientId)
	}else if (this.m_pServer != nil && this.m_pServer.m_PacketFuncList.Len() > 0){
		this.m_pServer.CallMsg("DISCONNECT", this.m_ClientId)
	}
}

func (this *WebSocketClient) Close() {
	if this.m_Conn != nil{
		this.m_Conn.Close()
	}
	this.m_Conn = nil
	close(this.m_WriteChan)
	this.Socket.Close()
	if this.m_pServer != nil {
		this.m_pServer.DelClinet(this)
	}
}

func (this *WebSocketClient) SendNoBlock(buff []byte) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("WriteBuf", err)
		}
	}()


	select {
	case this.m_WriteChan <- buff: //chan满后再写即阻塞，select进入default分支报错
	default:
		break
	}
}

func wserverclientRoutine(pClient *WebSocketClient) bool {
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

func wserverclientWriteRoutine(pClient *WebSocketClient) bool {
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