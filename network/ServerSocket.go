package network

import (
	"fmt"
	"gonet/rpc"
	"log"
	"net"
	"sync"
	"sync/atomic"

	"github.com/xtaci/kcp-go"
)

type IServerSocket interface {
	ISocket

	AssignClientId() uint32
	GetClientById(uint32) *ServerSocketClient
	LoadClient() *ServerSocketClient
	AddClinet(*net.TCPConn, string, int) *ServerSocketClient
	DelClinet(*ServerSocketClient) bool
	StopClient(uint32)
}

type ServerSocket struct {
	Socket
	m_nClientCount int
	m_nMaxClients  int
	m_nMinClients  int
	m_nIdSeed      uint32
	m_ClientList   map[uint32]*ServerSocketClient
	m_ClientLocker *sync.RWMutex
	m_Listen       *net.TCPListener
	m_Lock         sync.Mutex
	m_KcpListern   net.Listener
}

type ClientChan struct {
	pClient *ServerSocketClient
	state   int
	id      int
}

type WriteChan struct {
	buff []byte
	id   int
}

func (this *ServerSocket) Init(ip string, port int, params ...OpOption) bool {
	this.Socket.Init(ip, port, params...)
	this.m_ClientList = make(map[uint32]*ServerSocketClient)
	this.m_ClientLocker = &sync.RWMutex{}
	this.m_sIP = ip
	this.m_nPort = port
	return true
}
func (this *ServerSocket) Start() bool {
	if this.m_sIP == "" {
		this.m_sIP = "127.0.0.1"
	}

	var strRemote = fmt.Sprintf("%s:%d", this.m_sIP, this.m_nPort)
	//初始tcp
	tcpAddr, err := net.ResolveTCPAddr("tcp4", strRemote)
	if err != nil {
		log.Fatalf("%v", err)
	}
	this.m_Listen, err = net.ListenTCP("tcp4", tcpAddr)
	if err != nil {
		log.Fatalf("%v", err)
		return false
	}

	//初始kcp
	this.m_KcpListern, err = kcp.Listen(strRemote)
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Printf("启动监听，等待链接！\n")
	//延迟，监听关闭
	//defer ln.Close()
	go this.Run()
	go this.RunKcp()
	return true
}

func (this *ServerSocket) AssignClientId() uint32 {
	return atomic.AddUint32(&this.m_nIdSeed, 1)
}

func (this *ServerSocket) GetClientById(id uint32) *ServerSocketClient {
	this.m_ClientLocker.RLock()
	client, exist := this.m_ClientList[id]
	this.m_ClientLocker.RUnlock()
	if exist == true {
		return client
	}

	return nil
}

func (this *ServerSocket) AddClinet(conn net.Conn, addr string, connectType int) *ServerSocketClient {
	pClient := this.LoadClient()
	if pClient != nil {
		pClient.Init("", 0)
		pClient.m_pServer = this
		pClient.m_ReceiveBufferSize = this.m_ReceiveBufferSize
		pClient.SetMaxPacketLen(this.GetMaxPacketLen())
		pClient.m_ClientId = this.AssignClientId()
		pClient.m_sIP = addr
		pClient.SetConnectType(connectType)
		pClient.SetConn(conn)
		this.m_ClientLocker.Lock()
		this.m_ClientList[pClient.m_ClientId] = pClient
		this.m_ClientLocker.Unlock()
		pClient.Start()
		this.m_nClientCount++
		return pClient
	} else {
		log.Printf("%s", "无法创建客户端连接对象")
	}
	return nil
}

func (this *ServerSocket) DelClinet(pClient *ServerSocketClient) bool {
	this.m_ClientLocker.Lock()
	delete(this.m_ClientList, pClient.m_ClientId)
	this.m_ClientLocker.Unlock()
	return true
}

func (this *ServerSocket) StopClient(id uint32) {
	pClinet := this.GetClientById(id)
	if pClinet != nil {
		pClinet.Stop()
	}
}

func (this *ServerSocket) LoadClient() *ServerSocketClient {
	s := &ServerSocketClient{}
	return s
}

func (this *ServerSocket) Send(head rpc.RpcHead, packet rpc.Packet) int {
	pClient := this.GetClientById(head.SocketId)
	if pClient != nil {
		pClient.Send(head, packet)
	}
	return 0
}

func (this *ServerSocket) SendMsg(head rpc.RpcHead, funcName string, params ...interface{}) {
	pClient := this.GetClientById(head.SocketId)
	if pClient != nil {
		pClient.Send(head, rpc.Marshal(&head, &funcName, params...))
	}
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
	defer this.m_KcpListern.Close()
	this.Clear()
}

func (this *ServerSocket) Run() bool {
	for {
		tcpConn, err := this.m_Listen.AcceptTCP()
		handleError(err)
		if err != nil {
			return false
		}

		fmt.Printf("客户端：%s已连接！\n", tcpConn.RemoteAddr().String())
		//延迟，关闭链接
		//defer tcpConn.Close()
		this.handleConn(tcpConn, tcpConn.RemoteAddr().String())
	}
}

func (this *ServerSocket) RunKcp() bool {
	for {
		kcpConn, err := this.m_KcpListern.Accept()
		handleError(err)
		if err != nil {
			return false
		}

		fmt.Printf("kcp客户端：%s已连接！\n", kcpConn.RemoteAddr().String())
		//延迟，关闭链接
		//defer kcpConn.Close()
		this.handleConn(kcpConn, kcpConn.RemoteAddr().String())
	}
}

func (this *ServerSocket) handleConn(tcpConn net.Conn, addr string) bool {
	if tcpConn == nil {
		return false
	}

	pClient := this.AddClinet(tcpConn, addr, this.m_nConnectType)
	if pClient == nil {
		return false
	}

	return true
}
