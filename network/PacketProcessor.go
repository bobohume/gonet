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
		packetLen       int
		maxPacketLen    int
		littleEndian    bool
		maxPacketBuffer []byte //max receive buff
		packetFunc      HandlePacket
	}

	PacketConfig struct {
		MaxPacketLen *int
		Func         HandlePacket
	}
)

func NewPacketParser(conf PacketConfig) PacketParser {
	p := PacketParser{}
	p.packetLen = PACKET_LEN_DWORD
	p.maxPacketLen = base.MAX_PACKET
	p.littleEndian = true
	if conf.Func != nil {
		p.packetFunc = conf.Func
	} else {
		p.packetFunc = func(buff []byte) {
		}
	}
	return p
}

func (p *PacketParser) readLen(buff []byte) (bool, int) {
	nLen := len(buff)
	if nLen < p.packetLen {
		return false, 0
	}

	bufMsgLen := buff[:p.packetLen]
	// parse len
	var msgLen int
	switch p.packetLen {
	case PACKET_LEN_BYTE:
		msgLen = int(bufMsgLen[0])
	case PACKET_LEN_WORD:
		if p.littleEndian {
			msgLen = int(binary.LittleEndian.Uint16(bufMsgLen))
		} else {
			msgLen = int(binary.BigEndian.Uint16(bufMsgLen))
		}
	case PACKET_LEN_DWORD:
		if p.littleEndian {
			msgLen = int(binary.LittleEndian.Uint32(bufMsgLen))
		} else {
			msgLen = int(binary.BigEndian.Uint32(bufMsgLen))
		}
	}

	if msgLen+p.packetLen <= nLen {
		return true, msgLen + p.packetLen
	}

	return false, 0
}

func (p *PacketParser) Read(dat []byte) bool {
	buff := append(p.maxPacketBuffer, dat...)
	p.maxPacketBuffer = []byte{}
	nCurSize := 0
	//fmt.Println(p.maxPacketBuffer)
ParsePacekt:
	nPacketSize := 0
	nBufferSize := len(buff[nCurSize:])
	bFindFlag := false
	bFindFlag, nPacketSize = p.readLen(buff[nCurSize:])
	//fmt.Println(bFindFlag, nPacketSize, nBufferSize)
	if bFindFlag {
		if nBufferSize == nPacketSize { //完整包
			p.packetFunc(buff[nCurSize+p.packetLen : nCurSize+nPacketSize])
			nCurSize += nPacketSize
		} else if nBufferSize > nPacketSize {
			p.packetFunc(buff[nCurSize+p.packetLen : nCurSize+nPacketSize])
			nCurSize += nPacketSize
			goto ParsePacekt
		}
	} else if nBufferSize < p.maxPacketLen {
		p.maxPacketBuffer = buff[nCurSize:]
	} else {
		fmt.Println("超出最大包限制，丢弃该包")
		return false
	}
	return true
}

func (p *PacketParser) Write(dat []byte) []byte {
	// get len
	msgLen := len(dat)
	// check len
	if msgLen+p.packetLen > base.MAX_PACKET {
		fmt.Println("write over base.MAX_PACKET")
	}

	msg := make([]byte, p.packetLen+msgLen)
	// write len
	switch p.packetLen {
	case PACKET_LEN_BYTE:
		msg[0] = byte(msgLen)
	case PACKET_LEN_WORD:
		if p.littleEndian {
			binary.LittleEndian.PutUint16(msg, uint16(msgLen))
		} else {
			binary.BigEndian.PutUint16(msg, uint16(msgLen))
		}
	case PACKET_LEN_DWORD:
		if p.littleEndian {
			binary.LittleEndian.PutUint32(msg, uint32(msgLen))
		} else {
			binary.BigEndian.PutUint32(msg, uint32(msgLen))
		}
	}

	copy(msg[p.packetLen:], dat)
	return msg
}

/*func (p *Socket) ReceivePacket(Id uint32, dat []byte) bool{
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

	buff := append(p.m_MaxReceiveBuffer, dat...)
	p.m_MaxReceiveBuffer = []byte{}
	nCurSize := 0
	//fmt.Println(p.m_MaxReceiveBuffer)
ParsePacekt:
	nPacketSize := 0
	nBufferSize := len(buff[nCurSize:])
	bFindFlag := false
	bFindFlag, nPacketSize = seekToTcpEnd(buff[nCurSize:])
	//fmt.Println(bFindFlag, nPacketSize, nBufferSize)
	if bFindFlag{
		if nBufferSize == nPacketSize{		//完整包
			p.HandlePacket(Id, buff[nCurSize+base.TCP_HEAD_SIZE:nCurSize+nPacketSize])
			nCurSize += nPacketSize
		}else if ( nBufferSize > nPacketSize){
			p.HandlePacket(Id, buff[nCurSize+base.TCP_HEAD_SIZE:nCurSize+nPacketSize])
			nCurSize += nPacketSize
			goto ParsePacekt
		}
	}else if nBufferSize < p.m_MaxReceiveBufferSize{
		p.m_MaxReceiveBuffer = buff[nCurSize:]
	}else{
		fmt.Println("超出最大包限制，丢弃该包")
		return false
	}
	return true
}*/

//tcp粘包特殊结束标志
/*func (p *Socket) ReceivePacket(Id int, dat []byte) bool{
	//找包结束
	seekToTcpEnd := func(buff []byte) (bool, int) {
		nLen := bytes.Index(buff, []byte(base.TCP_END))
		if nLen != -1{
			return true, nLen+base.TCP_END_LENGTH
		}
		return false, 0
	}

	buff := append(p.m_MaxReceiveBuffer, dat...)
	p.m_MaxReceiveBuffer = []byte{}
	nCurSize := 0
	//fmt.Println(p.m_MaxReceiveBuffer)
ParsePacekt:
	nPacketSize := 0
	nBufferSize := len(buff[nCurSize:])
	bFindFlag := false
	bFindFlag, nPacketSize = seekToTcpEnd(buff[nCurSize:])
	//fmt.Println(bFindFlag, nPacketSize, nBufferSize)
	if bFindFlag {
		if nBufferSize == nPacketSize { //完整包
			p.HandlePacket(Id, buff[nCurSize:nCurSize+nPacketSize-base.TCP_END_LENGTH])
			nCurSize += nPacketSize
		} else if (nBufferSize > nPacketSize) {
			p.HandlePacket(Id, buff[nCurSize:nCurSize+nPacketSize-base.TCP_END_LENGTH])
			nCurSize += nPacketSize
			goto ParsePacekt
		}
	}else if nBufferSize < p.m_MaxReceiveBufferSize{
		p.m_MaxReceiveBuffer = buff[nCurSize:]
	}else{
		fmt.Println("超出最大包限制，丢弃该包")
		return false
	}
	return true
}*/
