package network

import (
	"fmt"
	"io"
	"log"
	"base"
	"strings"
	"crypto/sha1"
	"encoding/base64"
)

const(
	NEXT_FRAME = 0x0
	END_FRAME = 0x80
	ERROR_FRAME = 0xFF00
	INCOMPLETE_FRAME = 0xFE00
	OPENING_FRAME = 0x3300
	CLOSING_FRAME = 0x3400
	// 未完成的帧
	INCOMPLETE_TEXT_FRAME = 0x01
	INCOMPLETE_BINARY_FRAME = 0x02
	// 文本帧与二进制帧
	TEXT_FRAME = 0x81
	BINARY_FRAME = 0x82
	PING_FRAME = 0x19
	PONG_FRAME = 0x1A
	// 关闭连接
	CLOSE_FRAME = 0x08
)

type IWebSocketClient interface {
	ISocket
}

type WebSocketClient struct {
	Socket
	m_pServer     *WebSocket
}

func isWebSocket(buff []byte) bool{
	data := string(buff)
	if strings.Index(data, "Sec-WebSocket-Key") == -1{
		return false
	}

	if strings.Index(data,"GET") == -1{
		return false
	}

	return true
}

//------------------------------------------------------------------//
//						HTTP 报头格式
//------------------------------------------------------------------//
/*
	0                   1                   2                   3
	0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	+-+-+-+-+-------+-+-------------+-------------------------------+
	|F|R|R|R| opcode|M| Payload len |    Extended payload length    |
	|I|S|S|S|  (4)  |A|     (7)     |             (16/64)           |
	|N|V|V|V|       |S|             |   (if payload len==126/127)   |
	| |1|2|3|       |K|             |                               |
	+-+-+-+-+-------+-+-------------+ - - - - - - - - - - - - - - - +
	|     Extended payload length continued, if payload len == 127  |
	+ - - - - - - - - - - - - - - - +-------------------------------+
	|                               |Masking-key, if MASK set to 1  |
	+-------------------------------+-------------------------------+
	| Masking-key (continued)       |          Payload Data         |
	+-------------------------------- - - - - - - - - - - - - - - - +
	:                     Payload Data continued ...                :
	+ - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - +
	|                     Payload Data continued ...                |
	+---------------------------------------------------------------+
*/
func readFrame(data []byte) []byte{
	//第一个字节：FIN + RSV1-3 + OPCODE
	var maskkey [4]byte
	pos :=		0
	payloadLen := 0
	fin := 		data[0] >> 7
	rsv1 :=		data[0] >> 6 & 1
	rsv2 := 	data[0] >> 5 & 1
	rsv3 := 	data[0] >> 4 & 1
	opcode := 	data[0] & 15
	pos++
	log.Println(fin,rsv1,rsv2,rsv3,opcode)

	mask := data[pos] >> 7
	payloadLen = int(data[pos] & 0x7f)
	pos++

	if payloadLen == 126 {
		payloadLen = int(base.BytesToInt16(data[pos:2]))
		pos+= 2
	}else if payloadLen == 127 {
		payloadLen = int(base.BytesToInt64(data[pos:8]))
		pos+=8
	}

	if mask == 1{
		for i:=0; i < 4; i++{
			maskkey[i] = data[pos+i]
		}
		pos += 4
	}

	if mask ==1 {
		for i := 0; i < payloadLen; i++ {
			j := i % 4
			data[pos+i] = data[pos+i] ^ maskkey[j]
		}
	}

	data = data[pos:]
	return data
}

func writeFrame(data []byte) []byte{
	nLen := len(data)
	payloadFieldExtraBytes := 0
	if nLen < 126{
		payloadFieldExtraBytes = 0
	}else if nLen < (1<<16){
		payloadFieldExtraBytes = 2
	}else{
		payloadFieldExtraBytes = 8
	}

	frameHeader := make([]byte, 2 + payloadFieldExtraBytes)
	frameHeader[0] = BINARY_FRAME

	if nLen < 126{
		frameHeader[1] = byte(nLen)
	}else if nLen < (1<<16){
		frameHeader[1] = 0x7e
		buff := base.Htons(uint16(nLen))
		for i:=0; i < 2; i++{
			frameHeader[2+i] = buff[i]
		}
	} else{
		frameHeader[1] = 0x7f
		buff := base.Htonl(uint64(nLen))
		for i:=0; i < 8; i++{
			frameHeader[2+i] = buff[i]
		}
	}

	frameHeader = append(frameHeader, data...)
	return frameHeader
}

func handshark(data []byte, pClient *WebSocketClient)  bool {
	parseHandshake := func(content string) map[string]string {
		headers := make(map[string]string, 10)
		lines := strings.Split(content, "\r\n")
		for _,line := range lines {
			if len(line) >= 0 {
				words := strings.Split(line, ":")
				if len(words) == 2 {
					headers[strings.Trim(words[0]," ")] = strings.Trim(words[1], " ")
				}
			}
		}
		return headers
	}

	if isWebSocket(data){
		headers := parseHandshake(string(data))
		secWebsocketKey := headers["Sec-WebSocket-Key"]
		guid := "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
		// 计算Sec-WebSocket-Accept
		h := sha1.New()
		io.WriteString(h, secWebsocketKey + guid)
		accept := make([]byte, 28)
		base64.StdEncoding.Encode(accept, h.Sum(nil))
		response := "HTTP/1.1 101 Switching Protocols\r\n"
		response = response + "Sec-WebSocket-Accept: " + string(accept) + "\r\n"
		response = response + "Connection: Upgrade\r\n"
		response = response + "Upgrade: websocket\r\n\r\n"
		pClient.m_Conn.Write([]byte(response))
		pClient.m_nState = SSF_CONNECT
		return true
	}
	return false
}

func (this *WebSocketClient) Start() bool {
	if this.m_nState != SSF_SHUT_DOWN{
		return false
	}

	if this.m_pServer == nil {
		return false
	}

	this.m_nState = SSF_ACCEPT
	go wserverclientRoutine(this)
	return true
}

func (this *WebSocketClient) Send(buff []byte) int {
	if len(buff) > this.m_MaxSendBufferSize{
		log.Print(" SendError size",len(buff))
		return  0
	}

	buff = writeFrame(buff)
	n, err := this.m_Conn.Write(buff)
	handleError(err)
	if n > 0 {
		return n
	}
	return 0
}

func (this *WebSocketClient) ReceivePacket(Id int, buff []byte){
	if this.m_PacketFuncList.Len() > 0 {
		this.Socket.ReceivePacket(this.m_ClientId, buff)
	}else if (this.m_pServer != nil && this.m_pServer.m_PacketFuncList.Len() > 0){
		this.m_pServer.Socket.ReceivePacket(this.m_ClientId, buff)
	}
}

func (this *WebSocketClient) OnNetFail(error int) {
	this.Stop()
	this.CallPacket("DISCONNECT", this.m_ClientId)
}

func (this *WebSocketClient) Close() {
	this.Socket.Close()
	if this.m_pServer != nil {
		this.m_pServer.DelClinet(this)
	}
}

func wserverclientRoutine(pClient *WebSocketClient) bool {
	if pClient.m_Conn == nil {
		return false
	}

	for {
		if pClient.m_bShuttingDown {
			break
		}

		var buff = make([]byte, pClient.m_MaxReceiveBufferSize)
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

		if pClient.m_nState == SSF_ACCEPT{
			handshark(buff, pClient)
		}else{
			buff = readFrame(buff)
		}

		if string(buff[:n]) == "exit" {
			fmt.Printf("远程链接：%s退出！\n", pClient.m_Conn.RemoteAddr().String())
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