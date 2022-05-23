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
	clientCount  int
	maxClients   int
	minClients   int
	idSeed       uint32
	clientMap    map[uint32]*WebSocketClient
	clientLocker *sync.RWMutex
	lock         sync.Mutex
}

func (w *WebSocket) Init(ip string, port int, params ...OpOption) bool {
	w.Socket.Init(ip, port, params...)
	w.clientMap = make(map[uint32]*WebSocketClient)
	w.clientLocker = &sync.RWMutex{}
	w.ip = ip
	w.port = port
	return true
}

func (w *WebSocket) Start() bool {
	if w.ip == "" {
		w.ip = "127.0.0.1"
	}

	go func() {
		var strRemote = fmt.Sprintf("%s:%d", w.ip, w.port)
		http.Handle("/ws", websocket.Handler(w.wserverRoutine))
		err := http.ListenAndServe(strRemote, nil)
		if err != nil {
			fmt.Errorf("WebSocket ListenAndServe:%v", err)
		}
	}()

	fmt.Printf("WebSocket 启动监听，等待链接！\n")
	return true
}

func (w *WebSocket) AssignClientId() uint32 {
	return atomic.AddUint32(&w.idSeed, 1)
}

func (w *WebSocket) GetClientById(id uint32) *WebSocketClient {
	w.clientLocker.RLock()
	client, exist := w.clientMap[id]
	w.clientLocker.RUnlock()
	if exist == true {
		return client
	}

	return nil
}

func (w *WebSocket) AddClinet(tcpConn *websocket.Conn, addr string, connectType int) *WebSocketClient {
	client := w.LoadClient()
	if client != nil {
		client.Init("", 0)
		client.server = w
		client.receiveBufferSize = w.receiveBufferSize
		client.SetMaxPacketLen(w.GetMaxPacketLen())
		client.clientId = w.AssignClientId()
		client.ip = addr
		client.SetConn(tcpConn)
		client.SetConnectType(connectType)
		w.clientLocker.Lock()
		w.clientMap[client.clientId] = client
		w.clientLocker.Unlock()
		w.clientCount++
		return client
	} else {
		log.Printf("%s", "无法创建客户端连接对象")
	}
	return nil
}

func (w *WebSocket) DelClinet(client *WebSocketClient) bool {
	w.clientLocker.Lock()
	delete(w.clientMap, client.clientId)
	w.clientLocker.Unlock()
	return true
}

func (w *WebSocket) StopClient(id uint32) {
	client := w.GetClientById(id)
	if client != nil {
		client.Stop()
	}
}

func (w *WebSocket) LoadClient() *WebSocketClient {
	s := &WebSocketClient{}
	return s
}

func (w *WebSocket) Send(head rpc.RpcHead, packet rpc.Packet) int {
	client := w.GetClientById(head.SocketId)
	if client != nil {
		client.Send(head, packet)
	}
	return 0
}

func (w *WebSocket) SendMsg(head rpc.RpcHead, funcName string, params ...interface{}) {
	client := w.GetClientById(head.SocketId)
	if client != nil {
		client.Send(head, rpc.Marshal(&head, &funcName, params...))
	}
}

func (w *WebSocket) Restart() bool {
	return true
}
func (w *WebSocket) Connect() bool {
	return true
}
func (w *WebSocket) Disconnect(bool) bool {
	return true
}

func (w *WebSocket) OnNetFail(int) {
}

func (w *WebSocket) Close() {
	w.Clear()
}

func (w *WebSocket) wserverRoutine(conn *websocket.Conn) {
	fmt.Printf("客户端：%s已连接！\n", conn.RemoteAddr().String())
	w.handleConn(conn, conn.RemoteAddr().String())
}

func (w *WebSocket) handleConn(tcpConn *websocket.Conn, addr string) bool {
	if tcpConn == nil {
		return false
	}

	tcpConn.PayloadType = websocket.BinaryFrame
	client := w.AddClinet(tcpConn, addr, w.connectType)
	if client == nil {
		return false
	}

	client.Start()
	return true
}
