package network

import (
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"base"
)

const (
	ADD_CLIENT	= iota
	DEL_CLIENT= iota
	CLOSE_CLIENT=iota
)

type IServerSocket interface {
	ISocket

	AssignClientId() int
	GetClientById(int) *ServerSocketClient
	LoadClient() *ServerSocketClient
	AddClinet(*net.TCPConn, string, int) *ServerSocketClient
	DelClinet(*ServerSocketClient) bool
	StopClient(int)
}

type ServerSocket struct {
	Socket
	m_nClientCount  int
	m_nMaxClients   int
	m_nMinClients   int
	m_nIdSeed       int32
	m_bShuttingDown bool
	m_bCanAccept    bool
	m_bNagle        bool
	m_ClientList    map[int]*ServerSocketClient
	m_ClientChan	chan ClientChan
	m_SendChan 		chan SendChan
	m_Listen        *net.TCPListener
	m_Pool          sync.Pool
	m_Lock          sync.Mutex
}

type ClientChan struct{
	pClient *ServerSocketClient
	state int
	id int
}

type SendChan struct {
	buff	[]byte
	id		int
}

func (this *ServerSocket) Init(ip string, port int) bool {
	this.Socket.Init(ip, port)
	this.m_ClientList = make(map[int]*ServerSocketClient)
	this.m_ClientChan = make(chan ClientChan, 1000)
	this.m_SendChan = make(chan SendChan, 1000)
	this.m_sIP = ip
	this.m_nPort = port
	this.m_Pool = sync.Pool{
		New: func() interface{} {
			var s = &ServerSocketClient{}
			return s
		},
	}
	return true
}
func (this *ServerSocket) Start() bool {
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
	go serverRoutine(this)
	go timeRoutine(this)
	return true
}

func (this *ServerSocket) AssignClientId() int {
	atomic.AddInt32(&this.m_nIdSeed, 1)
	return int(this.m_nIdSeed)
}

func (this *ServerSocket) GetClientById(id int) *ServerSocketClient {
	client, exist := this.m_ClientList[id]
	if exist == true {
		return client
	}

	return nil
}

func (this *ServerSocket) AddClinet(tcpConn *net.TCPConn, addr string, connectType int) *ServerSocketClient {
	pClient := this.LoadClient()
	if pClient != nil {
		pClient.Socket.Init("", 0)
		pClient.m_pServer = this
		pClient.m_ClientId = this.AssignClientId()
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

func (this *ServerSocket) NotifyActor(pClient *ServerSocketClient, state int){
	if pClient != nil {
		var clientChan ClientChan
		clientChan.pClient = pClient
		clientChan.state = state
		clientChan.id = pClient.m_ClientId
		this.m_ClientChan <- clientChan
	}
}

func (this *ServerSocket) DelClinet(pClient *ServerSocketClient) bool {
	this.m_Pool.Put(pClient)
	this.NotifyActor(pClient, DEL_CLIENT)
	return true
}

func (this *ServerSocket) StopClient(id int){
	var clientChan ClientChan
	clientChan.pClient = nil
	clientChan.state = CLOSE_CLIENT
	clientChan.id = id
	this.m_ClientChan <- clientChan
}

func (this *ServerSocket) LoadClient() *ServerSocketClient {
	s := this.m_Pool.Get().(*ServerSocketClient)
	s.m_MaxReceiveBufferSize = this.m_MaxReceiveBufferSize
	s.m_MaxSendBufferSize = this.m_MaxSendBufferSize
	return s
}

func (this *ServerSocket) Stop() bool {
	if this.m_bShuttingDown {
		return true
	}

	this.m_bShuttingDown = true
	this.m_nState = SSF_SHUT_DOWN
	return true
}

func (this *ServerSocket) SendByID(id int, buff  []byte) int{
	var sendChan SendChan
	sendChan.buff = base.SetTcpEnd(buff)
	sendChan.id = id
	this.m_SendChan <- sendChan
	return  0
}

func (this *ServerSocket) SendMsgByID(id int, funcName string, params ...interface{}){
	var sendChan SendChan
	sendChan.buff = base.SetTcpEnd(base.GetPacket(funcName, params...))
	sendChan.id = id

	this.m_SendChan <- sendChan
}

func (this *ServerSocket) Restart() bool {
	return true
}
func (this *ServerSocket) Connect() bool {
	return true
}
func (this *ServerSocket) Disconnect(bool) bool {
	return true
}

func (this *ServerSocket) OnNetFail(int) {
}

func (this *ServerSocket) Close() {
	defer this.m_Listen.Close()
	this.Clear()
	//this.m_Pool.Put(this)
}

func SendClient(pClient *ServerSocketClient, buff []byte){
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("SendRpc", err) // 这里的err其实就是panic传入的内容，55
		}
	}()

	if pClient != nil{
		pClient.Send(buff)
	}
}

func serverRoutine(server *ServerSocket) {
	for {
		tcpConn, err := server.m_Listen.AcceptTCP()
		handleError(err)
		if err != nil {
			return
		}

		fmt.Printf("客户端：%s已连接！\n", tcpConn.RemoteAddr().String())
		//延迟，关闭链接
		//defer tcpConn.Close()
		handleConn(server, tcpConn, tcpConn.RemoteAddr().String())
	}
}

func timeRoutine(pServer *ServerSocket){
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
				go SendClient(pClient, sendChan.buff)
			}
		}
	}
}

func handleConn(server *ServerSocket, tcpConn *net.TCPConn, addr string) bool {
	if tcpConn == nil {
		return false
	}

	pClient := server.AddClinet(tcpConn, addr, server.m_nConnectType)
	if pClient == nil {
		return false
	}

	return true
}
