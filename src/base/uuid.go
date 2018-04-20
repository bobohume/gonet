package base

import (
	"time"
	"sync/atomic"
)
var(
	g_SeedId int32
)

const epochStart = 122192928000000000 / 1000000

func Uuid1(id int) uint64{
	var uid uint64
	atomic.AddInt32(&g_SeedId,1)
	uid |= ((uint64(time.Now().Unix())) << 32) & (0xFFFFFFFF00000000) //时间
	uid |= (uint64(RANDOMMGR().RandI(0, 0xFF)) << 24) & (0x00000000FF000000) //随机数
	uid |= (uint64(id % 0xFF) << 16) & (0x0000000000FF0000)//固定ID
	uid |= (uint64(g_SeedId % 0xFFFF)) & (0x000000000000FFFF)//自增ID
	return uid
}

func Uuid() uint64{
	var uid uint64
	atomic.AddInt32(&g_SeedId,1)
	curTime := epochStart + uint64(time.Now().UnixNano()/1000000)
	uid |= ((curTime) << 20) & (0xFFFFFFFFFFF00000) //时间
	uid |= (uint64(RANDOMMGR().RandI(0, 0xF)) << 16) & (0x00000000000F0000) //随机数
	uid |= (uint64(g_SeedId % 0xFFFF)) & (0x000000000000FFFF)//自增ID
	return uid
}