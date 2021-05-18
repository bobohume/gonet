package network

import (
	"encoding/binary"
	"fmt"
	"gonet/base"
)

const (
	PACKET_LEN_BYTE  = 1
	PACKET_LEN_WORD  = 2
	PACKET_LEN_DWORD = 4
)

// --------------
// | len | data |
// --------------
type (
	PacketParser struct {
		m_PacketLen   	int
		m_MaxPacketLen	int
		m_LittleEndian  bool
		m_MaxPacketBuffer []byte//max receive buff
		m_PacketFunc 	HandlePacket
	}

	PacketConfig struct {
		MaxPacketLen	*int
		Func 	HandlePacket
	}
)

func NewPacketParser(conf PacketConfig) PacketParser {
	p := PacketParser{}
	p.m_PacketLen = PACKET_LEN_DWORD
	p.m_MaxPacketLen = base.MAX_PACKET
	p.m_LittleEndian = true
	if conf.Func != nil{
		p.m_PacketFunc = conf.Func
	}else{
		p.m_PacketFunc = func(buff []byte) {
		}
	}
	return p
}

func (this *PacketParser) readLen(buff []byte) (bool, int){
	nLen := len(buff)
	if nLen < this.m_PacketLen{
		return false, 0
	}

	bufMsgLen := buff[:this.m_PacketLen]
	// parse len
	var msgLen int
	switch this.m_PacketLen {
	case PACKET_LEN_BYTE:
		msgLen = int(bufMsgLen[0])
	case PACKET_LEN_WORD:
		if this.m_LittleEndian {
			msgLen = int(binary.LittleEndian.Uint16(bufMsgLen))
		} else {
			msgLen = int(binary.BigEndian.Uint16(bufMsgLen))
		}
	case PACKET_LEN_DWORD:
		if this.m_LittleEndian {
			msgLen = int(binary.LittleEndian.Uint32(bufMsgLen))
		} else {
			msgLen = int(binary.BigEndian.Uint32(bufMsgLen))
		}
	}

	if msgLen + this.m_PacketLen <= nLen{
		return true, msgLen + this.m_PacketLen
	}

	return false, 0
}

func (this *PacketParser) Read(dat []byte) bool {
	buff := append(this.m_MaxPacketBuffer, dat...)
	this.m_MaxPacketBuffer = []byte{}
	nCurSize := 0
	//fmt.Println(this.m_MaxPacketBuffer)
ParsePacekt:
	nPacketSize := 0
	nBufferSize := len(buff[nCurSize:])
	bFindFlag := false
	bFindFlag, nPacketSize = this.readLen(buff[nCurSize:])
	//fmt.Println(bFindFlag, nPacketSize, nBufferSize)
	if bFindFlag{
		if nBufferSize == nPacketSize{		//完整包
			this.m_PacketFunc(buff[nCurSize + this.m_PacketLen : nCurSize + nPacketSize])
			nCurSize += nPacketSize
		}else if ( nBufferSize > nPacketSize){
			this.m_PacketFunc(buff[nCurSize + this.m_PacketLen : nCurSize + nPacketSize])
			nCurSize += nPacketSize
			goto ParsePacekt
		}
	}else if nBufferSize < this.m_MaxPacketLen{
		this.m_MaxPacketBuffer = buff[nCurSize:]
	}else{
		fmt.Println("超出最大包限制，丢弃该包")
		return false
	}
	return true
}

func (this *PacketParser) Write(dat []byte) []byte {
	// get len
	msgLen := len(dat)
	// check len
	if msgLen  + this.m_PacketLen > base.MAX_PACKET{
		fmt.Println("write over base.MAX_PACKET")
	}

	msg := make([]byte, this.m_PacketLen + msgLen)
	// write len
	switch this.m_PacketLen {
	case PACKET_LEN_BYTE:
		msg[0] = byte(msgLen)
	case PACKET_LEN_WORD:
		if this.m_LittleEndian {
			binary.LittleEndian.PutUint16(msg, uint16(msgLen))
		} else {
			binary.BigEndian.PutUint16(msg, uint16(msgLen))
		}
	case PACKET_LEN_DWORD:
		if this.m_LittleEndian {
			binary.LittleEndian.PutUint32(msg, uint32(msgLen))
		} else {
			binary.BigEndian.PutUint32(msg, uint32(msgLen))
		}
	}

	copy(msg[this.m_PacketLen:], dat)
	return msg
}

/*func (this *Socket) ReceivePacket(Id uint32, dat []byte) bool{
	//找包结束
	seekToTcpEnd := func(buff []byte) (bool, int){
		nLen := len(buff)
		if nLen < base.TCP_HEAD_SIZE{
			return false, 0
		}

		nSize := base.BytesToInt(buff[0:4])
		if nSize + base.TCP_HEAD_SIZE <= nLen{
			return true, nSize+base.TCP_HEAD_SIZE
		}
		return false, 0
	}

	buff := append(this.m_MaxReceiveBuffer, dat...)
	this.m_MaxReceiveBuffer = []byte{}
	nCurSize := 0
	//fmt.Println(this.m_MaxReceiveBuffer)
ParsePacekt:
	nPacketSize := 0
	nBufferSize := len(buff[nCurSize:])
	bFindFlag := false
	bFindFlag, nPacketSize = seekToTcpEnd(buff[nCurSize:])
	//fmt.Println(bFindFlag, nPacketSize, nBufferSize)
	if bFindFlag{
		if nBufferSize == nPacketSize{		//完整包
			this.HandlePacket(Id, buff[nCurSize+base.TCP_HEAD_SIZE:nCurSize+nPacketSize])
			nCurSize += nPacketSize
		}else if ( nBufferSize > nPacketSize){
			this.HandlePacket(Id, buff[nCurSize+base.TCP_HEAD_SIZE:nCurSize+nPacketSize])
			nCurSize += nPacketSize
			goto ParsePacekt
		}
	}else if nBufferSize < this.m_MaxReceiveBufferSize{
		this.m_MaxReceiveBuffer = buff[nCurSize:]
	}else{
		fmt.Println("超出最大包限制，丢弃该包")
		return false
	}
	return true
}*/

//tcp粘包特殊结束标志
/*func (this *Socket) ReceivePacket(Id int, dat []byte) bool{
	//找包结束
	seekToTcpEnd := func(buff []byte) (bool, int) {
		nLen := bytes.Index(buff, []byte(base.TCP_END))
		if nLen != -1{
			return true, nLen+base.TCP_END_LENGTH
		}
		return false, 0
	}

	buff := append(this.m_MaxReceiveBuffer, dat...)
	this.m_MaxReceiveBuffer = []byte{}
	nCurSize := 0
	//fmt.Println(this.m_MaxReceiveBuffer)
ParsePacekt:
	nPacketSize := 0
	nBufferSize := len(buff[nCurSize:])
	bFindFlag := false
	bFindFlag, nPacketSize = seekToTcpEnd(buff[nCurSize:])
	//fmt.Println(bFindFlag, nPacketSize, nBufferSize)
	if bFindFlag {
		if nBufferSize == nPacketSize { //完整包
			this.HandlePacket(Id, buff[nCurSize:nCurSize+nPacketSize-base.TCP_END_LENGTH])
			nCurSize += nPacketSize
		} else if (nBufferSize > nPacketSize) {
			this.HandlePacket(Id, buff[nCurSize:nCurSize+nPacketSize-base.TCP_END_LENGTH])
			nCurSize += nPacketSize
			goto ParsePacekt
		}
	}else if nBufferSize < this.m_MaxReceiveBufferSize{
		this.m_MaxReceiveBuffer = buff[nCurSize:]
	}else{
		fmt.Println("超出最大包限制，丢弃该包")
		return false
	}
	return true
}*/
