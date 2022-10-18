package network

import (
	"fmt"
	"gonet/base"
	"gonet/common/timer"
	"gonet/rpc"
	"io"
	"sync/atomic"
	"time"
)

type IWebSocketClient interface {
	ISocket
}

type WebSocketClient struct {
	Socket
	server   *WebSocket
	sendChan chan []byte //对外缓冲队列
	timerId  *int64
}

func (w *WebSocketClient) Init(ip string, port int, params ...OpOption) bool {
	w.Socket.Init(ip, port, params...)
	w.timerId = new(int64)
	return true
}

func (w *WebSocketClient) Start() bool {
	if w.server == nil {
		return false
	}

	if w.connectType == CLIENT_CONNECT {
		w.sendChan = make(chan []byte, MAX_SEND_CHAN)
		timer.StoreTimerId(w.timerId, int64(w.clientId)+1<<32)
		timer.RegisterTimer(w.timerId, (HEART_TIME_OUT/3)*time.Second, func() {
			w.Update()
		})
	}

	if w.packetFuncList.Len() == 0 {
		w.packetFuncList = w.server.packetFuncList
	}
	if w.connectType == CLIENT_CONNECT {
		go w.SendLoop()
	}
	w.Run()
	return true
}

func (w *WebSocketClient) Stop() bool {
	timer.RegisterTimer(w.timerId, timer.TICK_INTERVAL, func() {
		timer.StopTimer(w.timerId)
		if atomic.CompareAndSwapInt32(&w.state, SSF_RUN, SSF_STOP) {
			if w.conn != nil {
				w.conn.Close()
			}
		}
	})
	return false
}

func (w *WebSocketClient) Send(head rpc.RpcHead, packet rpc.Packet) int {
	defer func() {
		if err := recover(); err != nil {
			base.TraceCode(err)
		}
	}()

	if w.connectType == CLIENT_CONNECT { //对外链接send不阻塞
		select {
		case w.sendChan <- packet.Buff:
		default: //网络太卡,tcp send缓存满了并且发送队列也满了
			w.OnNetFail(1)
		}
	} else {
		return w.DoSend(packet.Buff)
	}
	return 0
}

func (w *WebSocketClient) DoSend(buff []byte) int {
	if w.conn == nil {
		return 0
	}

	n, err := w.conn.Write(w.packetParser.Write(buff))
	handleError(err)
	if n > 0 {
		return n
	}

	return 0
}

func (w *WebSocketClient) OnNetFail(error int) {
	w.Stop()

	if w.connectType == CLIENT_CONNECT { //netgate对外格式统一
		stream := base.NewBitStream(make([]byte, 32), 32)
		stream.WriteInt(int(DISCONNECTINT), 32)
		stream.WriteInt(int(w.clientId), 32)
		w.HandlePacket(stream.GetBuffer())
	} else {
		w.CallMsg(rpc.RpcHead{}, "DISCONNECT", w.clientId)
	}
	if w.server != nil {
		w.server.DelClinet(w)
	}
}

func (w *WebSocketClient) Close() {
	if w.connectType == CLIENT_CONNECT {
		//close(w.sendChan)
	}
	w.Socket.Close()
	if w.server != nil {
		w.server.DelClinet(w)
	}
}

func (w *WebSocketClient) Run() bool {
	var buff = make([]byte, w.receiveBufferSize)
	w.SetState(SSF_RUN)
	loop := func() bool {
		defer func() {
			if err := recover(); err != nil {
				base.TraceCode(err)
			}
		}()

		if w.conn == nil {
			return false
		}

		n, err := w.conn.Read(buff)
		if err == io.EOF {
			fmt.Printf("远程链接：%s已经关闭！\n", w.conn.RemoteAddr().String())
			w.OnNetFail(0)
			return false
		}
		if err != nil {
			handleError(err)
			w.OnNetFail(0)
			return false
		}
		if n > 0 {
			w.packetParser.Read(buff[:n])
		}
		w.heartTime = int(time.Now().Unix()) + HEART_TIME_OUT
		return true
	}

	for {
		if !loop() {
			break
		}
	}

	w.Close()
	fmt.Printf("%s关闭连接", w.ip)
	return true
}

// heart
func (w *WebSocketClient) Update() bool {
	now := int(time.Now().Unix())
	if w.heartTime < now {
		w.OnNetFail(2)
		return false
	}
	return true
}

func (w *WebSocketClient) SendLoop() bool {
	for {
		defer func() {
			if err := recover(); err != nil {
				base.TraceCode(err)
			}
		}()

		select {
		case buff := <-w.sendChan:
			if buff == nil { //信道关闭
				return false
			} else {
				w.DoSend(buff)
			}
		}
	}
	return true
}
