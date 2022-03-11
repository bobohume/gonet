package base

import (
	"gonet/base/vector"
	"log"
	"sync"
	"time"
)
/*
* Snowflake
*
* 1                                               39               52             64
* +-----------------------------------------------+----------------+---------------+
* | timestamp(ms)                                 | workerid       | sequence      |
* +-----------------------------------------------+----------------+---------------+
* | 0000000000 0000000000 0000000000 0000000000 0 | 0000000000 000 | 0000000000 00 |
* +-----------------------------------------------+----------------+---------------+
*
* 1. 39位时间截(毫秒级)，注意这是时间截的差值（当前时间截 - 开始时间截)。可以使用约7年: (1L << 39) / (1000L * 60 * 60 * 24 * 365) = 8
* 2. 13位数据机器位，可以部署在8096个节点
* 3. 12位序列，毫秒内的计数，同一机器，同一时间截并发4096个序号
*/
const (
	twepoch        = int64(1483228800000)             //开始时间截 (2017-01-01)
	workeridBits   = uint(13)                         //机器id所占的位数
	sequenceBits   = uint(12)                         //序列所占的位数
	workeridMax    = int64(-1 ^ (-1 << workeridBits)) //支持的最大机器id数量
	sequenceMask   = int64(-1 ^ (-1 << sequenceBits)) //
	workeridShift  = sequenceBits                     //机器id左移位数
	timestampShift = sequenceBits + workeridBits      //时间戳左移位数
	WorkeridMax    = workeridMax					  //集群自增量
)

type(
	Snowflake struct {
		sequence int64
		workerid int64
		timestamp int64
		sync.Mutex
	}

	ISnowflake interface {
		Init(workerid int64)
		UUID() int64
	}

	WorkIdQue struct {//workid que
		m_WorkMap map[uint32] int
		m_IdelVec *vector.Vector
		m_Id 	 int
	}

	IWorkIdQue interface {
		Init(int)
		Add(string) int
		Del(string) int
	}
)

func (this *Snowflake) Init(workerid int64){
	if workerid < 0 || workerid > workeridMax {
		log.Fatalln("workerid must be between 0 and 1023")
		return
	}

	this.workerid = workerid
	log.Println("snowflake [  workid : ", workerid, "]")
}

// Generate creates and returns a unique snowflake ID
func (this *Snowflake) UUID() int64 {
	this.Lock()
	now := time.Now().UnixNano() / 1000000
	if this.timestamp == now {
		this.sequence = (this.sequence + 1) & sequenceMask

		if this.sequence == 0 {
			for now <= this.timestamp {
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		this.sequence = 0
	}

	this.timestamp = now
	r := int64((now-twepoch)<<timestampShift | (this.workerid << workeridShift) | (this.sequence))
	this.Unlock()
	return r
}

func ParseUUID(id int64) (ts int64, workerId int64, seq int64) {
	seq = id & sequenceMask
	workerId = (id >> workeridShift) & workeridMax
	ts = (id >> timestampShift) + twepoch
	//t = time.Unix(ts/1000, (ts%1000)*1000000)
	return ts, workerId, seq
}

//----------WorkIdQue----------//
func (this *WorkIdQue) Init(id int){
	this.m_WorkMap	= make(map[uint32] int)
	this.m_IdelVec	= vector.NewVector()
	this.m_Id		= id
}

func (this *WorkIdQue) Add(val string) int{
	nVal := ToHash(val)
	nId, bExist := this.m_WorkMap[nVal]
	if bExist{
		return nId
	}

	if !this.m_IdelVec.Empty(){
		back := this.m_IdelVec.Back()
		nId = back.(int)
		this.m_IdelVec.PopBack()
		this.m_WorkMap[nVal] = nId
		return back.(int)
	}

	nId = this.m_Id
	this.m_WorkMap[nVal] = nId
	this.m_Id++
	return nId
}

func (this *WorkIdQue) Del(val string) int{
	nVal := ToHash(val)
	nId, bExist := this.m_WorkMap[nVal]
	if !bExist{
		return -1
	}
	delete(this.m_WorkMap, nVal)
	this.m_IdelVec.PushFront(nId)
	return nId
}

var(
	UUID = ISnowflake(&Snowflake{})
)

/*
* +-------------------------------------------------+--------+------------------+
* | timestamp(ms)                                 	| 随机数  | sequence         |
* +-------------------------------------------------+--------+------------------+
* | 0000000000 0000000000 0000000000 0000000000 000 | 0000   | 0000000000 000000|
* +-------------------------------------------------+--------+------------------+
*/
/*var(
	g_SeedId int32
)

func Uuid() int64{
	var uid int64
	atomic.AddInt32(&g_SeedId,1)
	curTime := 122192928000000000 / 1000000 + uint64(time.Now().UnixNano()/1000000)
	uid |= int64((curTime) << 20) & (0x7FFFFFFFFFF00000) //时间
	uid |= int64(uint64(RAND.RandI(0, 0xF)) << 16) & (0x00000000000F0000) //随机数
	uid |= int64(uint64(g_SeedId % 0xFFFF)) & (0x000000000000FFFF)//自增ID
	return uid
}*/