package network

import (
	"fmt"
	"gonet/base"
	"gonet/common/timer"
	"gonet/rpc"
	"hash/crc32"
	"io"
	"log"
	"time"
)

const (
	IDLE_TIMEOUT    = iota
	CONNECT_TIMEOUT = iota
	CONNECT_TYPE    = iota
)

var (
	DISCONNECTINT = crc32.ChecksumIEEE([]byte("DISCONNECT"))
	HEART_PACKET  = crc32.ChecksumIEEE([]byte("heardpacket"))
)

type IServerSocketClient interface {
	ISocket
}

type ServerSocketClient struct {
	Socket
	server   *ServerSocket
	sendChan chan []byte //对外缓冲队列
	timerId  *int64
}

func handleError(err error) {
	if err == nil {
		return
	}
	log.Printf("错误：%s\n", err.Error())
}

func (s *ServerSocketClient) Init(ip string, port int, params ...OpOption) bool {
	s.Socket.Init(ip, port, params...)
	return true
}

func (s *ServerSocketClient) Start() bool {
	if s.server == nil {
		return false
	}

	if s.connectType == CLIENT_CONNECT {
		s.sendChan = make(chan []byte, MAX_SEND_CHAN)
		s.timerId = new(int64)
		*s.timerId = int64(s.clientId)
		timer.RegisterTimer(s.timerId, (HEART_TIME_OUT/3)*time.Second, func() {
			s.Update()
		})
	}

	if s.packetFuncList.Len() == 0 {
		s.packetFuncList = s.server.packetFuncList
	}
	//s.m_Conn.SetKeepAlive(true)
	//s.m_Conn.SetKeepAlivePeriod(5*time.Second)
	go s.Run()
	if s.connectType == CLIENT_CONNECT {
		go s.SendLoop()
	}
	return true
}

func (s *ServerSocketClient) Send(head rpc.RpcHead, packet rpc.Packet) int {
	defer func() {
		if err := recover(); err != nil {
			base.TraceCode(err)
		}
	}()

	if s.connectType == CLIENT_CONNECT { //对外链接send不阻塞
		select {
		case s.sendChan <- packet.Buff:
		default: //网络太卡,tcp send缓存满了并且发送队列也满了
			s.OnNetFail(1)
		}
	} else {
		return s.DoSend(packet.Buff)
	}
	return 0
}

func (s *ServerSocketClient) DoSend(buff []byte) int {
	if s.conn == nil {
		return 0
	}

	n, err := s.conn.Write(s.packetParser.Write(buff))
	handleError(err)
	if n > 0 {
		return n
	}

	return 0
}

func (s *ServerSocketClient) OnNetFail(error int) {
	s.Stop()
	if s.connectType == CLIENT_CONNECT { //netgate对外格式统一
		stream := base.NewBitStream(make([]byte, 32), 32)
		stream.WriteInt(int(DISCONNECTINT), 32)
		stream.WriteInt(int(s.clientId), 32)
		s.HandlePacket(stream.GetBuffer())
	} else {
		s.CallMsg(rpc.RpcHead{}, "DISCONNECT", s.clientId)
	}
	if s.server != nil {
		s.server.DelClinet(s)
	}
}

func (s *ServerSocketClient) Close() {
	if s.connectType == CLIENT_CONNECT {
		s.sendChan <- nil
		//close(s.sendChan)
		timer.StopTimer(s.timerId)
	}
	s.Socket.Close()
	if s.server != nil {
		s.server.DelClinet(s)
	}
}

func (s *ServerSocketClient) Run() bool {
	var buff = make([]byte, s.receiveBufferSize)
	s.SetState(SSF_RUN)
	loop := func() bool {
		defer func() {
			if err := recover(); err != nil {
				base.TraceCode(err)
			}
		}()

		if s.conn == nil {
			return false
		}

		n, err := s.conn.Read(buff)
		if err == io.EOF {
			fmt.Printf("远程链接：%s已经关闭！\n", s.conn.RemoteAddr().String())
			s.OnNetFail(0)
			return false
		}
		if err != nil {
			handleError(err)
			s.OnNetFail(0)
			return false
		}
		if n > 0 {
			//熔断
			if !s.packetParser.Read(buff[:n]) && s.connectType == CLIENT_CONNECT {
				s.OnNetFail(1)
				return false
			}
		}
		s.heartTime = int(time.Now().Unix()) + HEART_TIME_OUT
		return true
	}

	for {
		if !loop() {
			break
		}
	}

	s.Close()
	fmt.Printf("%s关闭连接", s.ip)
	return true
}

// heart
func (s *ServerSocketClient) Update() {
	now := int(time.Now().Unix())
	// timeout
	if s.heartTime < now {
		s.OnNetFail(2)
		return
	}
}

func (s *ServerSocketClient) SendLoop() bool {
	for {
		defer func() {
			if err := recover(); err != nil {
				base.TraceCode(err)
			}
		}()

		select {
		case buff := <-s.sendChan:
			if buff == nil { //信道关闭
				return false
			} else {
				s.DoSend(buff)
			}
		}
	}

	return true
}
