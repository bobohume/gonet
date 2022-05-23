package network

import (
	"fmt"
	"gonet/base"
	"gonet/rpc"
	"io"
	"log"
	"net"

	"github.com/xtaci/kcp-go"
)

type IClientSocket interface {
	ISocket
}

type ClientSocket struct {
	Socket
	maxClients int
	minClients int
}

func (c *ClientSocket) Init(ip string, port int, params ...OpOption) bool {
	if c.port == port || c.ip == ip {
		return false
	}

	c.Socket.Init(ip, port, params...)
	c.ip = ip
	c.port = port
	fmt.Println(ip, port)
	return true
}
func (c *ClientSocket) Start() bool {
	if c.ip == "" {
		c.ip = "127.0.0.1"
	}

	if c.Connect() {
		go c.Run()
	}
	//延迟，监听关闭
	//defer ln.Close()
	return true
}

func (c *ClientSocket) SendMsg(head rpc.RpcHead, funcName string, params ...interface{}) {
	c.Send(head, rpc.Marshal(&head, &funcName, params...))
}

func (c *ClientSocket) Send(head rpc.RpcHead, packet rpc.Packet) int {
	defer func() {
		if err := recover(); err != nil {
			base.TraceCode(err)
		}
	}()

	if c.conn == nil {
		return 0
	}

	n, err := c.conn.Write(c.packetParser.Write(packet.Buff))
	handleError(err)
	if n > 0 {
		return n
	}
	//c.m_Writer.Flush()
	return 0
}

func (c *ClientSocket) Restart() bool {
	return true
}

func (c *ClientSocket) Connect() bool {
	var strRemote = fmt.Sprintf("%s:%d", c.ip, c.port)
	connectStr := "Tcp"
	if c.isKcp {
		ln, err1 := kcp.Dial(strRemote)
		if err1 != nil {
			return false
		}
		c.SetConn(ln)
		connectStr = "Kcp"
	} else {
		tcpAddr, err := net.ResolveTCPAddr("tcp4", strRemote)
		if err != nil {
			log.Printf("%v", err)
		}
		ln, err1 := net.DialTCP("tcp4", nil, tcpAddr)
		if err1 != nil {
			return false
		}
		c.SetConn(ln)
	}

	fmt.Printf("%s 连接成功，请输入信息！\n", connectStr)
	c.CallMsg(rpc.RpcHead{}, "COMMON_RegisterRequest")
	return true
}

func (c *ClientSocket) OnDisconnect() {
}

func (c *ClientSocket) OnNetFail(int) {
	c.Stop()
	c.CallMsg(rpc.RpcHead{}, "DISCONNECT", c.clientId)
}

func (c *ClientSocket) Run() bool {
	c.SetState(SSF_RUN)
	var buff = make([]byte, c.receiveBufferSize)
	loop := func() bool {
		defer func() {
			if err := recover(); err != nil {
				base.TraceCode(err)
			}
		}()

		if c.conn == nil {
			return false
		}

		n, err := c.conn.Read(buff)
		if err == io.EOF {
			fmt.Printf("远程链接：%s已经关闭！\n", c.conn.RemoteAddr().String())
			c.OnNetFail(0)
			return false
		}
		if err != nil {
			handleError(err)
			c.OnNetFail(0)
			return false
		}
		if n > 0 {
			c.packetParser.Read(buff[:n])
		}
		return true
	}

	for {
		if !loop() {
			break
		}
	}

	c.Close()
	return true
}
