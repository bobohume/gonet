package network

import (
	"fmt"
	"gonet/base"
	"io"
	"log"
)

type IWebSocketClient interface {
	ISocket
}

type WebSocketClient struct {
	Socket
	m_pServer     *WebSocket
}

func (this *WebSocketClient) Start() bool {
	if this.m_nState != SSF_SHUT_DOWN{
		return false
	}

	if this.m_pServer == nil {
		return false
	}

	if this.m_PacketFuncList.Len() == 0 {
		this.m_PacketFuncList = this.m_pServer.m_PacketFuncList
	}
	this.m_nState = SSF_ACCEPT
	go wserverclientRoutine(this)
	return true
}

func (this *WebSocketClient) Send(buff []byte) int {
	defer func() {
		if err := recover(); err != nil{
			fmt.Println("WebSocketClient Send", err)
		}
	}()

	n, err := this.m_Conn.Write(buff)
	handleError(err)
	if n > 0 {
		return n
	}
	return 0
}

func (this *WebSocketClient) OnNetFail(error int) {
	this.Stop()

	if this.m_nConnectType == SERVER_CONNECT{
		this.CallMsg("DISCONNECT", this.m_ClientId)
	}else{//netgate对外格式统一
		stream := base.NewBitStream(make([]byte, 32), 32)
		stream.WriteInt(int(DISCONNECTINT), 32)
		stream.WriteInt(this.m_ClientId, 32)
		this.HandlePacket(this.m_ClientId, stream.GetBuffer())
	}
}

func (this *WebSocketClient) Close() {
	if this.m_Conn != nil{
		this.m_Conn.Close()
	}
	this.m_Conn = nil
	this.Socket.Close()
	if this.m_pServer != nil {
		this.m_pServer.DelClinet(this)
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