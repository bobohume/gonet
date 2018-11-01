package network

import (
	"base"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"golang.org/x/net/websocket"
)

type IWebSocket interface {
	ISocket

	AssignClientId() int
	GetClientById(int) *WebSocketClient
	LoadClient() *WebSocketClient
	AddClinet(*websocket.Conn, string, int) *WebSocketClient
	DelClinet(*WebSocketClient) bool
	StopClient(int)
}

type WebSocket struct {
	Socket
	m_nClientCount  int
	m_nMaxClients   int
	m_nMinClients   int
	m_nIdSeed       int32
	m_bShuttingDown bool
	m_bCanAccept    bool
	m_bNagle        bool
	m_ClientList    map[int]*WebSocketClient
	m_ClientLocker	*sync.RWMutex
	m_ClientChan 	chan WClientChan
	m_Pool          sync.Pool
	m_Lock          sync.Mutex
}

type WClientChan struct{
	pClient *WebSocketClient
	state int
	id int
}

func (this *WebSocket) Init(ip string, port int) bool {
	this.Socket.Init(ip, port)
	this.m_ClientList = make(map[int]*WebSocketClient)
	this.m_ClientLocker = &sync.RWMutex{}
	this.m_ClientChan = make(chan WClientChan, 1000)
	this.m_sIP = ip
	this.m_nPort = port
	this.m_Pool = sync.Pool{
		New: func() interface{} {
			var s = &WebSocketClient{}
			return s
		},
	}
	return true
}

func (this *WebSocket) Start() bool {
	this.m_bShuttingDown = false

	if this.m_sIP == "" {
		this.m_sIP = "127.0.0.1"
	}

	go func() {
		var strRemote = fmt.Sprintf("%s:%d", this.m_sIP, this.m_nPort)
		http.Handle("/ws", websocket.Handler(this.wserverRoutine))
		err := http.ListenAndServe(strRemote, nil)
		if err != nil {
			fmt.Errorf("WebSocket ListenAndServe:", err)
		}
	}()

	fmt.Printf("WebSocket 启动监听，等待链接！\n")
	//延迟，监听关闭
	//defer ln.Close()
	this.m_nState = SSF_ACCEPT
	go wtimeRoutine(this)
	return true
}

func (this *WebSocket) AssignClientId() int {
	atomic.AddInt32(&this.m_nIdSeed, 1)
	return int(this.m_nIdSeed)
}

func (this *WebSocket) GetClientById(id int) *WebSocketClient {
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
		pClient.Socket.Init("", 0)
		pClient.m_pServer = this
		pClient.m_ClientId = this.AssignClientId()
		pClient.m_WebConn = tcpConn
		pClient.m_sIP = addr
		pClient.SetConnectType(connectType)
		this.NotifyActor(pClient, ADD_CLIENT)
		pClient.Start()
		this.m_nClientCount++
		return pClient
	} else {
		log.Print("%s", "无法创建客户端连接对象")
	}
	return nil
}

func (this *WebSocket) NotifyActor(pClient *WebSocketClient, state int){
	if pClient != nil {
		var clientChan WClientChan
		clientChan.pClient = pClient
		clientChan.state = state
		clientChan.id = pClient.m_ClientId
		this.m_ClientChan <- clientChan
	}
}

func (this *WebSocket) DelClinet(pClient *WebSocketClient) bool {
	this.m_Pool.Put(pClient)
	this.NotifyActor(pClient, DEL_CLIENT)
	return true
}

func (this *WebSocket) StopClient(id int){
	var clientChan WClientChan
	clientChan.pClient = nil
	clientChan.state = CLOSE_CLIENT
	clientChan.id = id
	this.m_ClientChan <- clientChan
}

func (this *WebSocket) LoadClient() *WebSocketClient {
	s := this.m_Pool.Get().(*WebSocketClient)
	s.m_MaxReceiveBufferSize = this.m_MaxReceiveBufferSize
	s.m_MaxSendBufferSize = this.m_MaxSendBufferSize
	return s
}

func (this *WebSocket) Stop() bool {
	if this.m_bShuttingDown {
		return true
	}

	this.m_bShuttingDown = true
	this.m_nState = SSF_SHUT_DOWN
	return true
}

func (this *WebSocket) SendByID(id int, buff  []byte) int{
	pClient := this.GetClientById(id)
	if pClient != nil{
		pClient.Send(base.SetTcpEnd(buff))
	}
	return  0
}

func (this *WebSocket) SendMsgByID(id int, funcName string, params ...interface{}){
	pClient := this.GetClientById(id)
	if pClient != nil{
		pClient.Send(base.SetTcpEnd(base.GetPacket(funcName, params...)))
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
	//this.m_Pool.Put(this)
}

func (this *WebSocket)wserverRoutine(conn *websocket.Conn){
	fmt.Printf("客户端：%s已连接！\n", conn.RemoteAddr().String())
	whandleConn(this, conn, conn.RemoteAddr().String())
}

func wtimeRoutine(pServer *WebSocket){
	for{
		select {
		case clientChan := <- pServer.m_ClientChan:
			if clientChan.state == ADD_CLIENT{
				if clientChan.pClient != nil {
					pServer.m_ClientList[clientChan.pClient.m_ClientId] = clientChan.pClient
				}
			}else if clientChan.state == DEL_CLIENT{
				delete(pServer.m_ClientList, clientChan.id)
			}else if clientChan.state == CLOSE_CLIENT{
				pClinet := pServer.GetClientById(clientChan.id)
				if pClinet != nil{
					pClinet.Stop()
				}
			}
		}
	}
}

func whandleConn(server *WebSocket, tcpConn *websocket.Conn, addr string) bool {
	if tcpConn == nil {
		return false
	}

	pClient := server.AddClinet(tcpConn, addr, server.m_nConnectType)
	if pClient == nil {
		return false
	}

	return true
}