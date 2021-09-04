package network

import (
	"fmt"
	"gonet/rpc"
	"log"
	"net/http"
	"sync"
	"sync/atomic"

	"golang.org/x/net/websocket"
)

type IWebSocket interface {
	ISocket

	AssignClientId() uint32
	GetClientById(uint32) *WebSocketClient
	LoadClient() *WebSocketClient
	AddClinet(*websocket.Conn, string, int) *WebSocketClient
	DelClinet(*WebSocketClient) bool
	StopClient(uint32)
}

type WebSocket struct {
	Socket
	m_nClientCount int
	m_nMaxClients  int
	m_nMinClients  int
	m_nIdSeed      uint32
	m_ClientList   map[uint32]*WebSocketClient
	m_ClientLocker *sync.RWMutex
	m_Lock         sync.Mutex
}

func (this *WebSocket) Init(ip string, port int) bool {
	this.Socket.Init(ip, port)
	this.m_ClientList = make(map[uint32]*WebSocketClient)
	this.m_ClientLocker = &sync.RWMutex{}
	this.m_sIP = ip
	this.m_nPort = port
	return true
}

func (this *WebSocket) Start() bool {
	if this.m_sIP == "" {
		this.m_sIP = "127.0.0.1"
	}

	go func() {
		var strRemote = fmt.Sprintf("%s:%d", this.m_sIP, this.m_nPort)
		http.Handle("/ws", websocket.Handler(this.wserverRoutine))
		err := http.ListenAndServe(strRemote, nil)
		if err != nil {
			fmt.Errorf("WebSocket ListenAndServe:%v", err)
		}
	}()

	fmt.Printf("WebSocket 启动监听，等待链接！\n")
	return true
}

func (this *WebSocket) AssignClientId() uint32 {
	return atomic.AddUint32(&this.m_nIdSeed, 1)
}

func (this *WebSocket) GetClientById(id uint32) *WebSocketClient {
	this.m_ClientLocker.RLock()
	client, exist := this.m_ClientList[id]
	this.m_ClientLocker.RUnlock()
	if exist == true {
		return client
	}

	return nil
}

func (this *WebSocket) AddClinet(tcpConn *websocket.Conn, addr string, connectType int) *WebSocketClient {
	pClient := this.LoadClient()
	if pClient != nil {
		pClient.Init("", 0)
		pClient.m_pServer = this
		pClient.m_ReceiveBufferSize = this.m_ReceiveBufferSize
		pClient.SetMaxPacketLen(this.GetMaxPacketLen())
		pClient.m_ClientId = this.AssignClientId()
		pClient.m_sIP = addr
		pClient.SetTcpConn(tcpConn)
		pClient.SetConnectType(connectType)
		this.m_ClientLocker.Lock()
		this.m_ClientList[pClient.m_ClientId] = pClient
		this.m_ClientLocker.Unlock()
		this.m_nClientCount++
		return pClient
	} else {
		log.Printf("%s", "无法创建客户端连接对象")
	}
	return nil
}

func (this *WebSocket) DelClinet(pClient *WebSocketClient) bool {
	this.m_ClientLocker.Lock()
	delete(this.m_ClientList, pClient.m_ClientId)
	this.m_ClientLocker.Unlock()
	return true
}

func (this *WebSocket) StopClient(id uint32) {
	pClinet := this.GetClientById(id)
	if pClinet != nil {
		pClinet.Stop()
	}
}

func (this *WebSocket) LoadClient() *WebSocketClient {
	s := &WebSocketClient{}
	return s
}

func (this *WebSocket) Send(head rpc.RpcHead, buff []byte) int {
	pClient := this.GetClientById(head.SocketId)
	if pClient != nil {
		pClient.Send(head, buff)
	}
	return 0
}

func (this *WebSocket) SendMsg(head rpc.RpcHead, funcName string, params ...interface{}) {
	pClient := this.GetClientById(head.SocketId)
	if pClient != nil {
		pClient.Send(head, rpc.Marshal(head, funcName, params...))
	}
}

func (this *WebSocket) Restart() bool {
	return true
}
func (this *WebSocket) Connect() bool {
	return true
}
func (this *WebSocket) Disconnect(bool) bool {
	return true
}

func (this *WebSocket) OnNetFail(int) {
}

func (this *WebSocket) Close() {
	this.Clear()
}

func (this *WebSocket) wserverRoutine(conn *websocket.Conn) {
	fmt.Printf("客户端：%s已连接！\n", conn.RemoteAddr().String())
	this.handleConn(conn, conn.RemoteAddr().String())
}

func (this *WebSocket) handleConn(tcpConn *websocket.Conn, addr string) bool {
	if tcpConn == nil {
		return false
	}

	tcpConn.PayloadType = websocket.BinaryFrame
	pClient := this.AddClinet(tcpConn, addr, this.m_nConnectType)
	if pClient == nil {
		return false
	}

	pClient.Start()
	return true
}
