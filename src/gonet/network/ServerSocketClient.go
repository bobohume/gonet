package network

import (
	"fmt"
	"gonet/base"
	"hash/crc32"
	"io"
	"log"
	"net"
)

const (
	IDLE_TIMEOUT    = iota
	CONNECT_TIMEOUT = iota
	CONNECT_TYPE    = iota
)

var(
	DISCONNECTINT = crc32.ChecksumIEEE([]byte("DISCONNECT"))
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

	if this.m_PacketFuncList.Len() == 0 {
		this.m_PacketFuncList = this.m_pServer.m_PacketFuncList
	}
	this.m_nState = SSF_CONNECT
	this.m_Conn.(*net.TCPConn).SetNoDelay(true)
	//this.m_Conn.SetKeepAlive(true)
	//this.m_Conn.SetKeepAlivePeriod(5*time.Second)
	go serverclientRoutine(this)
	return true
}

func (this *ServerSocketClient) Send(buff []byte) int {
	defer func() {
		if err := recover(); err != nil{
			fmt.Println("ServerSocketClient Send", err)
		}
	}()

	n, err := this.m_Conn.Write(buff)
	handleError(err)
	if n > 0 {
		return n
	}
	//this.m_Writer.Flush()
	return 0
}

func (this *ServerSocketClient) OnNetFail(error int) {
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