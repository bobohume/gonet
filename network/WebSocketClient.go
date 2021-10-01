package network

import (
	"fmt"
	"gonet/base"
	"gonet/common/timer"
	"gonet/rpc"
	"io"
	"time"
)

type IWebSocketClient interface {
	ISocket
}

type WebSocketClient struct {
	Socket
	m_pServer  *WebSocket
	m_SendChan chan []byte //对外缓冲队列
	m_TimerId  *int64
}

func (this *WebSocketClient) Init(ip string, port int, params ...OpOption) bool {
	this.Socket.Init(ip, port, params...)
	return true
}

func (this *WebSocketClient) Start() bool {
	if this.m_pServer == nil {
		return false
	}

	if this.m_nConnectType == CLIENT_CONNECT {
		this.m_SendChan = make(chan []byte, MAX_SEND_CHAN)
		this.m_TimerId = new(int64)
		*this.m_TimerId = int64(this.m_ClientId)
		timer.RegisterTimer(this.m_TimerId, (HEART_TIME_OUT/3)*time.Second, func() {
			this.Update()
		})
	}

	if this.m_PacketFuncList.Len() == 0 {
		this.m_PacketFuncList = this.m_pServer.m_PacketFuncList
	}
	if this.m_nConnectType == CLIENT_CONNECT {
		go this.SendLoop()
	}
	this.Run()
	return true
}

func (this *WebSocketClient) Send(head rpc.RpcHead, buff []byte) int {
	defer func() {
		if err := recover(); err != nil {
			base.TraceCode(err)
		}
	}()

	if this.m_nConnectType == CLIENT_CONNECT { //对外链接send不阻塞
		select {
		case this.m_SendChan <- buff:
		default: //网络太卡,tcp send缓存满了并且发送队列也满了
			this.OnNetFail(1)
		}
	} else {
		return this.DoSend(buff)
	}
	return 0
}

func (this *WebSocketClient) DoSend(buff []byte) int {
	if this.m_Conn == nil {
		return 0
	}

	n, err := this.m_Conn.Write(this.m_PacketParser.Write(buff))
	handleError(err)
	if n > 0 {
		return n
	}

	return 0
}

func (this *WebSocketClient) OnNetFail(error int) {
	this.Stop()

	if this.m_nConnectType == CLIENT_CONNECT { //netgate对外格式统一
		stream := base.NewBitStream(make([]byte, 32), 32)
		stream.WriteInt(int(DISCONNECTINT), 32)
		stream.WriteInt(int(this.m_ClientId), 32)
		this.HandlePacket(stream.GetBuffer())
	} else {
		this.CallMsg("DISCONNECT", this.m_ClientId)
	}
	if this.m_pServer != nil {
		this.m_pServer.DelClinet(this)
	}
}

func (this *WebSocketClient) Close() {
	if this.m_nConnectType == CLIENT_CONNECT {
		//close(this.m_SendChan)
		timer.StopTimer(this.m_TimerId)
	}
	this.Socket.Close()
	if this.m_pServer != nil {
		this.m_pServer.DelClinet(this)
	}
}

func (this *WebSocketClient) Run() bool {
	var buff = make([]byte, this.m_ReceiveBufferSize)
	this.SetState(SSF_RUN)
	loop := func() bool {
		defer func() {
			if err := recover(); err != nil {
				base.TraceCode(err)
			}
		}()

		if this.m_Conn == nil {
			return false
		}

		n, err := this.m_Conn.Read(buff)
		if err == io.EOF {
			fmt.Printf("远程链接：%s已经关闭！\n", this.m_Conn.RemoteAddr().String())
			this.OnNetFail(0)
			return false
		}
		if err != nil {
			handleError(err)
			this.OnNetFail(0)
			return false
		}
		if n > 0 {
			this.m_PacketParser.Read(buff[:n])
		}
		return true
	}

	for {
		if !loop() {
			break
		}
	}

	this.Close()
	fmt.Printf("%s关闭连接", this.m_sIP)
	return true
}

// heart
func (this *WebSocketClient) Update() bool {
	now := int(time.Now().Unix())
	if this.m_HeartTime < now {
		this.OnNetFail(2)
		return false
	}
	return true
}

func (this *WebSocketClient) SendLoop() bool {
	for {
		defer func() {
			if err := recover(); err != nil {
				base.TraceCode(err)
			}
		}()

		select {
		case buff := <-this.m_SendChan:
			if buff == nil { //信道关闭
				return false
			} else {
				this.DoSend(buff)
			}
		}
	}
	return true
}
