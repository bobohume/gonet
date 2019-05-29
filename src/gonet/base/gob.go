package base

import (
	"bytes"
	"encoding/gob"
	"reflect"
)

// gobencode
type GobNode struct {
	enc *gob.Encoder//gob解析器new耗费很大
	dec *gob.Decoder
	buf [2]*bytes.Buffer
}

type GobList struct {
	m_Chan chan int
	m_Gobs []*GobNode
}

func (this *GobList) Init(num int) {
	this.m_Gobs = make([]*GobNode, num)
	this.m_Chan = make(chan int, num)
	for i := 0; i < num; i++ {
		node := &GobNode{}
		node.buf[0] = &bytes.Buffer{}
		node.enc = gob.NewEncoder(node.buf[0])
		node.buf[1] = &bytes.Buffer{}
		node.dec = gob.NewDecoder(node.buf[1])
		this.m_Gobs[i] = node
		this.m_Chan <- i
	}
}

func (this *GobList) Marsh(data interface{}) []byte{
	n := <- this.m_Chan
	node := this.m_Gobs[n]
	node.buf[0].Reset()
	node.enc.Encode(data)
	this.m_Chan <- n
	return node.buf[0].Bytes()
}

func (this *GobList) UnMarsh(buf []byte, data interface{}) {
	n := <- this.m_Chan
	node := this.m_Gobs[n]
	node.buf[1].Reset()
	node.buf[1].Write(buf)
	node.dec.Decode(data)
	this.m_Chan <- n
}

func (this *GobList) UnMarshValue(buf []byte, data reflect.Value){
	n := <- this.m_Chan
	node := this.m_Gobs[n]
	node.buf[1].Reset()
	node.buf[1].Write(buf)
	node.dec.DecodeValue(data)
	this.m_Chan <- n
}

var(
	GOB *GobList
)

func init(){
	GOB = new(GobList)
	GOB.Init(10)
}