package network

import (
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"base"
)

type IWebSocket interface {
	ISocket

	AssignClientId() int
	GetClientById(int) *WebSocketClient
	LoadClient() *WebSocketClient
	AddClinet(*net.TCPConn, string, int) *WebSocketClient
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
	m_ClientChan	chan WClientChan
	m_SendChan 		chan SendChan
	m_Listen        *net.TCPListener
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
	this.m_ClientChan = make(chan WClientChan, 1000)
	this.m_SendChan = make(chan SendChan, 1000)
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

	var strRemote = fmt.Sprintf("%s:%d", this.m_sIP, this.m_nPort)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", strRemote)
	if err != nil {
		log.Fatal("%v", err)
	}
	ln, err := net.ListenTCP("tcp4", tcpAddr)
	if err != nil {
		log.Fatal("%v", err)
		return false
	}

	fmt.Printf("启动监听，等待链接！\n")

	this.m_Listen = ln
	//延迟，监听关闭
	//defer ln.Close()
	this.m_nState = SSF_ACCEPT
	go wserverRoutine(this)
	go wtimeRoutine(this)
	return true
}

func (this *WebSocket) AssignClientId() int {
	atomic.AddInt32(&this.m_nIdSeed, 1)
	return int(this.m_nIdSeed)
}

func (this *WebSocket) GetClientById(id int) *WebSocketClient {
	client, exist := this.m_ClientList[id]
	if exist == true {
		return client
	}

	return nil
}

func (this *WebSocket) AddClinet(tcpConn *net.TCPConn, addr string, connectType int) *WebSocketClient {
	pClient := this.LoadClient()
	if pClient != nil {
		pClient.Socket.Init("", 0)
		pClient.m_pServer = this
		pClient.m_ClientId = this.AssignClientId()
		pClient.m_Conn = tcpConn
		pClient.m_sIP = addr
		pClient.SetConnectType(connectType)
		pClient.SetTcpConn(tcpConn)
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
	var sendChan SendChan
	sendChan.buff = buff
	sendChan.id = id
	this.m_SendChan <- sendChan
	return  0
}

func (this *WebSocket) SendMsgByID(id int, funcName string, params ...interface{}){
	var sendChan SendChan
	sendChan.buff = base.GetPacket(funcName, params...)
	sendChan.id = id
	this.m_SendChan <- sendChan
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
	defer this.m_Listen.Close()
	this.Clear()
	//this.m_Pool.Put(this)
}

func WSendClient(pClient *WebSocketClient, buff []byte){
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("SendRpc", err) // 这里的err其实就是panic传入的内容，55
		}
	}()

	if pClient != nil{
		pClient.Send(buff)
	}
}

func wserverRoutine(server *WebSocket) {
	for {
		tcpConn, err := server.m_Listen.AcceptTCP()
		handleError(err)
		if err != nil {
			return
		}

		fmt.Printf("客户端：%s已连接！\n", tcpConn.RemoteAddr().String())
		//延迟，关闭链接
		//defer tcpConn.Close()
		whandleConn(server, tcpConn, tcpConn.RemoteAddr().String())
	}
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
		case sendChan := <- pServer.m_SendChan:
			pClient := pServer.GetClientById(sendChan.id)
			if pClient != nil{
				go WSendClient(pClient, sendChan.buff)
			}
		}
	}
}

func whandleConn(server *WebSocket, tcpConn *net.TCPConn, addr string) bool {
	if tcpConn == nil {
		return false
	}

	pClient := server.AddClinet(tcpConn, addr, server.m_nConnectType)
	if pClient == nil {
		return false
	}

	return true
}