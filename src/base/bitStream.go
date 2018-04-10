package base

//----------------bitsream---------------
//for example
//buf := make([]byte, 256)
//var bitstream base.BitStream
//bitstream.BuildPacketStream(buf, 256)
//bitstream.WriteInt(1000, 16)
// or
//bitstream := NewBitStream(buf)
//----------------------------------------

import (
	"fmt"
)

const (
	Bit8   = 8
	Bit16  = 16
	Bit32  = 32
	Bit64  = 64
	Bit128 = 128
)

type (
	BitStream struct {
		dataPtr        []byte
		bitNum         int
		flagNum        int
		tailFlag       bool
		bufSize        int
		bitsLimite     int
		error          bool
		maxReadBitNum  int
		maxWriteBitNum int
		zipflag        byte   //zip flag
		tmpBuf         []byte //zip buf
		zipSize        int    //zip length
	}

	IBitStream interface {
		BuildPacketStream([]byte, int) bool
		setBuffer([]byte, int, int)
		GetBuffer() []byte
		GetBytePtr() []byte
		GetReadByteSize() int
		GetCurPos() int
		GetPosition() int
		GetStreamSize() int
		SetPosition(int) bool
		clear()

		WriteBits(int, []byte)
		ReadBits(int, []byte)
		WriteInt(int, int)
		ReadInt(int) int
		ReadFlag() bool
		WriteFlag()
		WriteString(string)
		ReadString() string

		WriteInt64(int64, int)
		ReadInt64(int) int64
		WriteFloat(float32)
		ReadFloat() float32
		WriteFLoat64(float64)
		ReadFloat64() float64
	}
)

func (this *BitStream) BuildPacketStream(buffer []byte, writeSize int) bool {
	if writeSize == 0 {
		return false
	}

	this.setBuffer(buffer, writeSize, -1)
	this.SetPosition(0)
	fmt.Print("")
	return true
}

func (this *BitStream) setBuffer(bufPtr []byte, size int, maxSize int) {
	this.dataPtr = bufPtr
	this.bitNum = 0
	this.flagNum = 0
	this.tailFlag = false
	this.bufSize = size
	this.maxReadBitNum = size << 3
	if maxSize < 0 {
		maxSize = size
	}
	this.maxWriteBitNum = maxSize << 3
	this.bitsLimite = size << 3
	this.error = false
}

func (this *BitStream) GetBuffer() []byte {
	return this.dataPtr[0:this.GetPosition()]
}

func (this *BitStream) GetBytePtr() []byte {
	return this.dataPtr[this.GetPosition():]
}

func (this *BitStream) GetReadByteSize() int {
	return (this.maxReadBitNum >> 3) - this.GetPosition()
}

func (this *BitStream) getCurPos() int {
	return this.bitNum
}

func (this *BitStream) clear() {
	var buff []byte
	buff = make([]byte, this.bufSize)
	this.dataPtr = buff
}

func (this *BitStream) GetPosition() int {
	return (this.bitNum + 7) >> 3
}

func (this *BitStream) GetStreamSize() int {
	return this.bufSize
}

func (this *BitStream) SetPosition(pos int) bool {
	Assert(pos == 0 || this.flagNum == 0, "不正确的setPosition调用")
	if pos != 0 && this.flagNum != 0 {
		return false
	}

	this.bitNum = pos << 3
	this.flagNum = 0
	return true
}

func (this *BitStream) WriteBits(bitCount int, bitPtr []byte) {
	if bitCount == 0 {
		return
	}

	if this.tailFlag {
		this.error = true
		Assert(false, "Out of range write")
		return
	}

	if (bitCount & 0x7) != 0 {
		bitCount = (bitCount & ^0x7) + 8
	}

	if bitCount+this.bitNum > this.maxWriteBitNum {
		this.error = true
		Assert(false, "Out of range write")
		return
	}

	byteCount := (bitCount + 7) >> 3
	for i, v := range bitPtr[:byteCount] {
		this.dataPtr[(this.bitNum>>3)+i] = v
	}
	this.bitNum += bitCount
}

func (this *BitStream) ReadBits(bitCount int, bitPtr []byte) {
	if bitCount == 0 {
		return
	}

	if this.tailFlag {
		this.error = true
		Assert(false, "Out of range read")
		return
	}

	if (bitCount & 0x7) != 0 {
		bitCount = (bitCount & ^0x7) + 8
	}

	if bitCount+this.bitNum > this.maxReadBitNum {
		this.error = true
		Assert(false, "Out of range read")
		return
	}

	byteCount := (bitCount + 7) >> 3
	stPtr := this.dataPtr[(this.bitNum >> 3) : (this.bitNum>>3)+byteCount]
	for i, v := range stPtr[:] {
		bitPtr[i] = v
	}
	this.bitNum += bitCount
}

func (this *BitStream) WriteInt(value int, bitCount int) {
	this.WriteBits(bitCount, IntToBytes(value))
}

func (this *BitStream) ReadInt(bitCount int) int {
	var ret int
	buf := make([]byte, 4)
	this.ReadBits(bitCount, buf)
	ret = BytesToInt(buf)
	if bitCount == Bit32 {
		return int(ret)
	} else {
		ret &= (1 << uint32(bitCount)) - 1
	}

	return int(ret)
}

func (this *BitStream) ReadFlag() bool {
	if ((this.flagNum - (this.flagNum>>3)<<3) == 0) && !this.tailFlag {
		this.flagNum = this.bitNum
		if this.bitNum+8 < this.maxReadBitNum {
			this.bitNum += 8
		} else {
			this.tailFlag = true
		}
	}

	if this.flagNum+1 > this.maxReadBitNum {
		this.error = true
		Assert(false, "Out of range read")
		return false
	}

	mask := 1 << uint32(this.flagNum&0x7)
	ret := (int(this.dataPtr[(this.flagNum>>3)]) & mask) != 0
	this.flagNum++
	return ret
}

func (this *BitStream) WriteFlag(val bool) bool {
	if ((this.flagNum - (this.flagNum>>3)<<3) == 0) && !this.tailFlag {
		this.flagNum = this.bitNum

		if this.bitNum+8 < this.maxWriteBitNum {
			this.bitNum += 8 //Ray; 跳开8个用于写flag
		} else {
			this.tailFlag = true
		}
	}

	if this.flagNum+1 > this.maxWriteBitNum {
		this.error = true
		Assert(false, "Out of range write")
		return false
	}

	if val {
		this.dataPtr[(this.flagNum >> 3)] |= (1 << uint32(this.flagNum&0x7))
	} else {
		this.dataPtr[(this.flagNum >> 3)] &= ^(1 << uint32(this.flagNum&0x7))
	}
	this.flagNum++
	return (val)
}

func (this *BitStream) ReadString() string {
	if this.ReadFlag() {
		nLen := this.ReadInt(Bit16)
		buf := make([]byte, nLen)
		this.ReadBits(nLen<<3, buf)
		return string(buf)
	}
	return string("")
}

func (this *BitStream) WriteString(str string) {
	buf := []byte(str)
	nLen := len(buf)

	if this.WriteFlag(nLen > 0) {
		this.WriteInt(nLen, Bit16)
		this.WriteBits(nLen<<3, buf)
	}
}

func (this *BitStream) WriteInt64(value int64, bitCount int) {
	this.WriteBits(bitCount, Int64ToBytes(value))
}

func (this *BitStream) ReadInt64(bitCount int) int64 {
	var ret int64
	buf := make([]byte, 8)
	this.ReadBits(bitCount, buf)
	ret = BytesToInt64(buf)
	if bitCount == Bit64 {
		return int64(ret)
	} else {
		ret &= (1 << uint64(bitCount)) - 1
	}

	return int64(ret)
}

func (this *BitStream) WriteFloat(value float32) {
	this.WriteBits(Bit32, Float32ToByte(value))
}

func (this *BitStream) ReadFloat() float32 {
	var ret float32
	buf := make([]byte, 4)
	this.ReadBits(Bit32, buf)
	ret = ByteToFloat32(buf)

	return float32(ret)
}

func (this *BitStream) WriteFloat64(value float64) {
	this.WriteBits(Bit64, Float64ToByte(value))
}

func (this *BitStream) ReadFloat64() float64 {
	var ret float64
	buf := make([]byte, 8)
	this.ReadBits(Bit64, buf)
	ret = ByteToFloat64(buf)

	return float64(ret)
}

func NewBitStream(buf []byte, nLen int) *BitStream {
	var bitstream BitStream
	bitstream.BuildPacketStream(buf, nLen)
	return &bitstream
}
