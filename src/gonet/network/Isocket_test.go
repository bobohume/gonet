package network_test

import (
	"gonet/base"
	"bytes"
	"fmt"
	"gonet/rpc"
	"testing"
)

var(
	m_pInBuffer []byte
	m_pInBuffer1 []byte
	nTimes = 1000
)

const (
	TCP_END	  = "💞♡"					//解决tpc粘包半包,结束标志
	ARRAY_LEN = 100000					//800kb 100 * 1000 * 8
)

func TestFixHead(t *testing.T)  {
	t.Log("固定长度")
	buff := []byte{}
	for j := 0; j < 1; j++{
		buff = append(buff, SetTcpEnd(rpc.Marshal("test1",[ARRAY_LEN]int64{1,2,3,4,5,6}))...)
	}
	for i :=0;i < nTimes;i++{
		ReceivePacket(0,buff)
	}
	
}

func TestEndFlag(t *testing.T)  {
	t.Log("c语言结束标志", []byte(TCP_END))
	buff := []byte{}
	for j := 0; j < 1; j++{
		buff = append(buff, SetTcpEnd1(rpc.Marshal("test1", [ARRAY_LEN]int64{1,2,3,4,5,6}))...)
	}
	for i :=0;i < nTimes;i++{
		ReceivePacket1(0,buff)
	}
}

func SetTcpEnd(buff []byte) []byte{
	buff = append(base.IntToBytes(len(buff)), buff...)
	return buff
}

func SetTcpEnd1(buff []byte) []byte{
	buff = append(buff, []byte(TCP_END)...)
	return buff
}

func ReceivePacket(Id int, dat []byte){
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("ReceivePacket", err) // 接受包错误
		}
	}()
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

	buff := append(m_pInBuffer, dat...)
	m_pInBuffer = []byte{}
	nCurSize := 0
	//fmt.Println(this.m_pInBuffer)
ParsePacekt:
	nPacketSize := 0
	nBufferSize := len(buff[nCurSize:])
	bFindFlag := false
	bFindFlag, nPacketSize = seekToTcpEnd(buff[nCurSize:])
	//fmt.Println(bFindFlag, nPacketSize, nBufferSize)
	if bFindFlag{
		if nBufferSize == nPacketSize{		//完整包
			//this.HandlePacket(Id, buff[nCurSize+base.TCP_HEAD_SIZE:nCurSize+nPacketSize])
			nCurSize += nPacketSize
		}else if ( nBufferSize > nPacketSize){
			//this.HandlePacket(Id, buff[nCurSize+base.TCP_HEAD_SIZE:nCurSize+nPacketSize])
			nCurSize += nPacketSize
			goto ParsePacekt
		}
	}else if nBufferSize < base.MAX_PACKET{
		m_pInBuffer = buff[nCurSize:]
	}else{
		fmt.Println("超出最大包限制，丢弃该包")
	}
}

func  ReceivePacket1(Id int, dat []byte){
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("ReceivePacket1", err) // 接受包错误
		}
	}()
	//找包结束
	seekToTcpEnd := func(buff []byte) (bool, int){
		i := bytes.Index(buff, []byte(TCP_END))
		if i != -1{
			return true, i+2
		}
		return false, 0
	}

	buff := append(m_pInBuffer1, dat...)
	m_pInBuffer1 = []byte{}
	nCurSize := 0
	//fmt.Println(this.m_pInBuffer)
ParsePacekt:
	nPacketSize := 0
	nBufferSize := len(buff[nCurSize:])
	bFindFlag := false
	bFindFlag, nPacketSize = seekToTcpEnd(buff[nCurSize:])
	//fmt.Println(bFindFlag, nPacketSize, nBufferSize)
	if bFindFlag{
		if nBufferSize == nPacketSize{		//完整包
			nCurSize += nPacketSize
		}else if ( nBufferSize > nPacketSize){
			nCurSize += nPacketSize
			goto ParsePacekt
		}
	}else if nBufferSize < base.MAX_PACKET{
		m_pInBuffer1 = buff[nCurSize:]
	}else{
		fmt.Println("超出最大包限制，丢弃该包")
	}
}