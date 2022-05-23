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

const (
	Bit8              = 8
	Bit16             = 16
	Bit32             = 32
	Bit64             = 64
	Bit128            = 128
	MAX_PACKET        = 1 * 1024 * 1024 //1MB
	MAX_CLIENT_PACKET = 10 * 1024       //10KB
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
		resize() bool

		WriteBits([]byte, int)
		ReadBits(int) []byte
		WriteInt(int, int)
		ReadInt(int) int
		ReadFlag() bool
		WriteFlag(bool) bool
		WriteString(string)
		ReadString() string

		WriteInt64(int64, int)
		ReadInt64(int) int64
		WriteFloat(float32)
		ReadFloat() float32
		WriteFloat64(float64)
		ReadFloat64() float64
	}
)

func (b *BitStream) BuildPacketStream(buffer []byte, writeSize int) bool {
	if writeSize <= 0 {
		return false
	}

	b.setBuffer(buffer, writeSize, -1)
	b.SetPosition(0)
	return true
}

func (b *BitStream) setBuffer(bufPtr []byte, size int, maxSize int) {
	b.dataPtr = bufPtr
	b.bitNum = 0
	b.flagNum = 0
	b.tailFlag = false
	b.bufSize = size
	b.maxReadBitNum = size << 3
	if maxSize < 0 {
		maxSize = size
	}
	b.maxWriteBitNum = maxSize << 3
	b.bitsLimite = size
	b.error = false
}

func (b *BitStream) GetBuffer() []byte {
	return b.dataPtr[0:b.GetPosition()]
}

func (b *BitStream) GetBytePtr() []byte {
	return b.dataPtr[b.GetPosition():]
}

func (b *BitStream) GetReadByteSize() int {
	return (b.maxReadBitNum >> 3) - b.GetPosition()
}

func (b *BitStream) GetCurPos() int {
	return b.bitNum
}

func (b *BitStream) GetPosition() int {
	return (b.bitNum + 7) >> 3
}

func (b *BitStream) GetStreamSize() int {
	return b.bufSize
}

func (b *BitStream) SetPosition(pos int) bool {
	Assert(pos == 0 || b.flagNum == 0, "不正确的setPosition调用")
	if pos != 0 && b.flagNum != 0 {
		return false
	}

	b.bitNum = pos << 3
	b.flagNum = 0
	return true
}

func (b *BitStream) clear() {
	var buff []byte
	buff = make([]byte, b.bufSize)
	b.dataPtr = buff
}

func (b *BitStream) resize() bool {
	//fmt.Println("BitStream Resize")
	b.dataPtr = append(b.dataPtr, make([]byte, b.bitsLimite)...)
	size := b.bitsLimite * 2
	if size <= 0 || size >= MAX_PACKET*2 {
		return false
	}
	b.bufSize = size
	b.maxReadBitNum = size << 3
	b.maxWriteBitNum = size << 3
	b.bitsLimite = size
	return true
}

func (b *BitStream) WriteBits(bitPtr []byte, bitCount int) {
	if bitCount == 0 {
		return
	}

	if b.tailFlag {
		b.error = true
		Assert(false, "Out of range write")
		return
	}

	if (bitCount & 0x7) != 0 {
		bitCount = (bitCount & ^0x7) + 8
	}

	for bitCount+b.bitNum > b.maxWriteBitNum {
		if !b.resize() {
			b.error = true
			Assert(false, "Out of range write")
			return
		}
	}

	bitNum := b.bitNum >> 3
	byteCount := (bitCount + 7) >> 3
	copy(b.dataPtr[bitNum:], bitPtr[:byteCount])
	/*for i, v := range bitPtr[:byteCount] {
		b.dataPtr[bitNum+i] = v
	}*/
	b.bitNum += bitCount
}

func (b *BitStream) ReadBits(bitCount int) []byte {
	if bitCount == 0 {
		return []byte{}
	}

	if b.tailFlag {
		b.error = true
		Assert(false, "Out of range read")
		return []byte{}
	}

	if (bitCount & 0x7) != 0 {
		bitCount = (bitCount & ^0x7) + 8
	}

	for bitCount+b.bitNum > b.maxReadBitNum {
		if !b.resize() {
			b.error = true
			Assert(false, "Out of range read")
			return []byte{}
		}
	}

	byteCount := (bitCount + 7) >> 3
	bitNum := b.bitNum >> 3
	stPtr := b.dataPtr[bitNum : bitNum+byteCount]
	b.bitNum += bitCount
	return stPtr
}

func (b *BitStream) WriteInt(value int, bitCount int) {
	b.WriteBits(IntToBytes(value), bitCount)
}

func (b *BitStream) ReadInt(bitCount int) int {
	var ret int
	buf := b.ReadBits(bitCount)
	ret = BytesToInt(buf)
	if bitCount == Bit32 {
		return int(ret)
	} else {
		ret &= (1 << uint32(bitCount)) - 1
	}

	return int(ret)
}

func (b *BitStream) ReadFlag() bool {
	if ((b.flagNum - (b.flagNum>>3)<<3) == 0) && !b.tailFlag {
		b.flagNum = b.bitNum
		if b.bitNum+8 < b.maxReadBitNum {
			b.bitNum += 8
		} else {
			if !b.resize() {
				b.tailFlag = true
			} else {
				b.bitNum += 8
			}
		}
	}

	if b.flagNum+1 > b.maxReadBitNum {
		b.error = true
		Assert(false, "Out of range read")
		return false
	}

	mask := 1 << uint32(b.flagNum&0x7)
	ret := (int(b.dataPtr[(b.flagNum>>3)]) & mask) != 0
	b.flagNum++
	return ret
}

func (b *BitStream) WriteFlag(value bool) bool {
	if ((b.flagNum - (b.flagNum>>3)<<3) == 0) && !b.tailFlag {
		b.flagNum = b.bitNum

		if b.bitNum+8 < b.maxWriteBitNum {
			b.bitNum += 8 //跳开8个用于写flag
		} else {
			if !b.resize() {
				b.tailFlag = true
			} else {
				b.bitNum += 8 //跳开8个用于写flag
			}
		}
	}

	if b.flagNum+1 > b.maxWriteBitNum {
		b.error = true
		Assert(false, "Out of range write")
		return false
	}

	if value {
		b.dataPtr[(b.flagNum >> 3)] |= (1 << uint32(b.flagNum&0x7))
	} else {
		b.dataPtr[(b.flagNum >> 3)] &= ^(1 << uint32(b.flagNum&0x7))
	}
	b.flagNum++
	return (value)
}

func (b *BitStream) ReadString() string {
	if b.ReadFlag() {
		nLen := b.ReadInt(Bit16)
		buf := b.ReadBits(nLen << 3)
		return string(buf)
	}
	return string("")
}

func (b *BitStream) WriteString(value string) {
	buf := []byte(value)
	nLen := len(buf)

	if b.WriteFlag(nLen > 0) {
		b.WriteInt(nLen, Bit16)
		b.WriteBits(buf, nLen<<3)
	}
}

func (b *BitStream) WriteInt64(value int64, bitCount int) {
	b.WriteBits(Int64ToBytes(value), bitCount)
}

func (b *BitStream) ReadInt64(bitCount int) int64 {
	var ret int64
	buf := b.ReadBits(bitCount)
	ret = BytesToInt64(buf)
	if bitCount == Bit64 {
		return int64(ret)
	} else {
		ret &= (1 << uint64(bitCount)) - 1
	}

	return int64(ret)
}

func (b *BitStream) WriteFloat(value float32) {
	b.WriteBits(Float32ToByte(value), Bit32)
}

func (b *BitStream) ReadFloat() float32 {
	var ret float32
	buf := b.ReadBits(Bit32)
	ret = BytesToFloat32(buf)

	return float32(ret)
}

func (b *BitStream) WriteFloat64(value float64) {
	b.WriteBits(Float64ToByte(value), Bit64)
}

func (b *BitStream) ReadFloat64() float64 {
	var ret float64
	buf := b.ReadBits(Bit64)
	ret = BytesToFloat64(buf)

	return float64(ret)
}

func NewBitStream(buf []byte, nLen int) *BitStream {
	var bitstream BitStream
	bitstream.BuildPacketStream(buf, nLen)
	return &bitstream
}
