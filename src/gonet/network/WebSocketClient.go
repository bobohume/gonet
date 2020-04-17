package network

import (
	"fmt"
	"gonet/base"
	"io"
)

type IWebSocketClient interface {
	ISocket
}

type WebSocketClient struct {
	Socket
	m_pServer     *WebSocket
	m_SendChan	chan []byte//对外缓冲队列
}

func (this *WebSocketClient) Init(ip string, port int) bool {
	if this.m_nConnectType == CLIENT_CONNECT {
		this.m_SendChan = make(chan []byte, MAX_SEND_CHAN)
	}
	this.Socket.Init(ip, port)
	return true
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
	if this.m_nConnectType == CLIENT_CONNECT {
		go this.SendLoop()
	}
	this.Run()
	return true
}

func (this *WebSocketClient) Send(buff []byte) int {
	defer func() {
		if err := recover(); err != nil{
			base.TraceCode(err)
		}
	}()

	if this.m_nConnectType == CLIENT_CONNECT  {//对外链接send不阻塞
		select {
		case this.m_SendChan <- buff:
		default://网络太卡,tcp send缓存满了并且发送队列也满了
			this.OnNetFail(1)
			if this.m_pServer != nil {
				this.m_pServer.DelClinet(this)
			}
		}
	}else{
		return this.DoSend(buff)
	}
	return 0
}

func (this *WebSocketClient) DoSend(buff []byte) int {
	if this.m_Conn == nil{
		return 0
	}

	n, err := this.m_Conn.Write(buff)
	handleError(err)
	if n > 0 {
		return n
	}

	return 0
}

func (this *WebSocketClient) OnNetFail(error int) {
	this.Stop()
	
	if this.m_nConnectType == CLIENT_CONNECT{//netgate对外格式统一
		stream := base.NewBitStream(make([]byte, 32), 32)
		stream.WriteInt(int(DISCONNECTINT), 32)
		stream.WriteInt(this.m_ClientId, 32)
		this.HandlePacket(this.m_ClientId, stream.GetBuffer())
	}else{
		this.CallMsg("DISCONNECT", this.m_ClientId)
	}
}

func (this *WebSocketClient) Close() {
	if this.m_nConnectType == CLIENT_CONNECT {
		close(this.m_SendChan)
	}
	this.Socket.Close()
	if this.m_pServer != nil {
		this.m_pServer.DelClinet(this)
	}
}

func (this *WebSocketClient) Run() bool {
	defer func() {
		if err := recover(); err != nil {
			base.TraceCode(err)
		}
	}()

	var buff= make([]byte, this.m_ReceiveBufferSize)
	for {
		if this.m_bShuttingDown || this.m_Conn == nil{
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
	fmt.Printf("%s关闭连接", this.m_sIP)
	return true
}

func (this *WebSocketClient) SendLoop() bool {
	for {
		select {
		case buff := <-this.m_SendChan:
			if buff == nil{//信道关闭
				return false
			}else{
				this.DoSend(buff)
			}
		}
	}
	return true
}